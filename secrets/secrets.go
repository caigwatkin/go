/*
Copyright 2018 Cai Gwatkin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package secrets

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	go_errors "github.com/caigwatkin/go/errors"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// Client interface for secrets
type Client interface {
	Decrypt(secret Secret) ([]byte, error)
	Encrypt(secret []byte) (*Secret, error)
	Secret(domain, kind string) ([]byte, error)
	SecretFromFile(pathToFile string) (*Secret, error)
	DownloadAndDecryptAndCache(ctx context.Context, bucket, dir string, required Required) error
}

type Config struct {
	Env             string
	GcpProjectId    string
	CloudkmsKeyRing string
	CloudkmsKey     string
}

// NewClient returns an implementation of the client interface that allows secret management
func NewClient(ctx context.Context, config Config) (Client, error) {
	cloudkmsService, err := cloudkms.NewService(ctx)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed initializing cloudkms service")
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed initializing storage client")
	}
	return client{
		cloudkmsService: cloudkmsService,
		env:             config.Env,
		cryptoKey:       fmt.Sprintf("projects/%s/locations/global/keyRings/%s/cryptoKeys/%s", config.GcpProjectId, config.CloudkmsKeyRing, config.CloudkmsKey),
		secrets:         make(map[string][]byte),
		storageClient:   storageClient,
	}, nil
}

// Required secrets map, where the key is the domain of the secret and the values are the types of secrets
//
// This should map to the naming scheme of the encrypted secret file, e.g.:
//   - Secret file naming should be "secret_domain-secret_type-cloudkms_env.json"
//   - If a required secret is from some api, it is a key, the domain is "some_api" and the type "key"
//   - If the file was encrypted using cloudkms in a "dev" env, the file name is "some_api-key-cloudkms_dev.json"
type Required map[string][]string

// ReduceRequired secrets into one set
func ReduceRequired(required ...Required) Required {
	reduced := make(Required)
	for i := 0; i < len(required); i++ {
		for newK, newV := range required[i] {
			newV = reduceRequiredTypes(nil, newV)
			if reducedV, ok := reduced[newK]; ok {
				reduced[newK] = reduceRequiredTypes(reducedV, newV)
				continue
			}
			reduced[newK] = newV
		}
	}
	return reduced
}

func reduceRequiredTypes(reduced, new []string) []string {
	for _, newv := range new {
		var exists bool
		for _, reducedV := range reduced {
			if newv == reducedV {
				exists = true
				break
			}
		}
		if !exists {
			reduced = append(reduced, newv)
		}
	}
	return reduced
}

type client struct {
	cloudkmsService *cloudkms.Service
	cryptoKey       string
	env             string
	secrets         map[string][]byte
	storageClient   *storage.Client
}

// Secret data model, a subset of properties of a cloudkms secret
type Secret struct {
	Name       string `json:"name,omitempty"`
	Ciphertext string `json:"ciphertext,omitempty"`
}

// Encrypt plaintext bytes into a secret
func (c client) Encrypt(plaintext []byte) (*Secret, error) {
	r, err := c.cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(c.cryptoKey, &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(plaintext),
	}).Do()
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed to encrypt plaintext")
	}
	return &Secret{
		Ciphertext: r.Ciphertext,
		Name:       r.Name,
	}, nil
}

// Decrypt a secret into plaintext bytes
func (c client) Decrypt(secret Secret) ([]byte, error) {
	resp, err := c.cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(c.cryptoKey, &cloudkms.DecryptRequest{
		Ciphertext: secret.Ciphertext,
	}).Do()
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed to decrypt ciphertext")
	}
	buf, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed decoding plaintext as base64")
	}
	return buf, nil
}

// Secret returns an unencrypted secret from the cache if one exists, else errors
func (c client) Secret(domain, kind string) ([]byte, error) {
	if v, ok := c.secrets[cacheKey(domain, kind)]; ok {
		return v, nil
	}
	return nil, go_errors.Errorf("No secret for domain %q and kind %q", domain, kind)
}

// SecretFromFile returns a secret from a file
func (c client) SecretFromFile(pathToFile string) (*Secret, error) {
	buf, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed reading file")
	}
	var s Secret
	if err := json.Unmarshal(buf, &s); err != nil {
		return nil, go_errors.Wrap(err, "Failed unmarshalling file buf into Secret")
	}
	return &s, nil
}

// DownloadAndDecryptAndCache required secrets from a GCP cloud bucket
func (c client) DownloadAndDecryptAndCache(ctx context.Context, bucket, dir string, required Required) error {
	b := c.storageClient.Bucket(bucket)
	for domain, kinds := range required {
		for _, kind := range kinds {
			s, err := c.download(ctx, b, dir, domain, kind)
			if err != nil {
				return go_errors.Wrap(err, "Failed downloading secret from bucket")
			}
			plaintext, err := c.Decrypt(*s)
			if err != nil {
				return go_errors.Wrap(err, "Failed decrypting secret")
			}
			c.secrets[cacheKey(domain, kind)] = plaintext
		}
	}
	return nil
}

func (c client) download(ctx context.Context, bucket *storage.BucketHandle, dir, domain, kind string) (*Secret, error) {
	var file string
	if dir != "" {
		file = fmt.Sprintf("%s/%s-%s-cloudkms_%s.json", dir, domain, kind, c.env)
	} else {
		file = fmt.Sprintf("%s-%s-cloudkms_%s.json", domain, kind, c.env)
	}
	fileObject := bucket.Object(file)
	reader, err := fileObject.NewReader(ctx)
	if err != nil {
		return nil, go_errors.Wrapf(err, "Failed opening file %q", file)
	}
	defer reader.Close()
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(reader); err != nil {
		return nil, go_errors.Wrap(err, "Failed reading from reader")
	}
	var s Secret
	if err := json.Unmarshal(buffer.Bytes(), &s); err != nil {
		return nil, go_errors.Wrap(err, "Failed unmarshalling buffer into Secret")
	}
	return &s, nil
}

func cacheKey(domain, kind string) string {
	return fmt.Sprintf("%s%s", domain, kind)
}
