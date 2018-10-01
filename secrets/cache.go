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
	"encoding/json"
	"fmt"

	"cloud.google.com/go/storage"
	go_errors "github.com/caigwatkin/go/errors"
)

// Secret returns an unencrypted secret from the cache if one exists, else errors
func (c client) Secret(domain, kind string) ([]byte, error) {
	if v, ok := c.secrets[cacheKey(domain, kind)]; ok {
		return v, nil
	}
	return nil, go_errors.Errorf("No secret for domain %q and type %q", domain, kind)
}

// Required secrets map, where the key is the domain of the secret and the values are the types of secrets
//
// This should map to the naming scheme of the encrypted secret file, e.g.:
//   - Secret file naming should be "secret_domain-secret_type-cloudkms_env.json"
//   - If a required secret is from some api, it is a key, the domain is "some_api" and the type "key"
//   - If the file was encrypted using cloudkms in a "dev" env, the file name is "some_api-key-cloudkms_dev.json"
type Required map[string][]string

// Combine two or more maps of required secrets into one
func Combine(to Required, from ...Required) Required {
	for _, v := range from {
		for k, vv := range v {
			if s, ok := to[k]; ok {
				to[k] = append(s, vv...)
			} else {
				to[k] = vv
			}
		}
	}
	return to
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
