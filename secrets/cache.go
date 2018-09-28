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
	"github.com/pkg/errors"
)

func (c client) PlaintextFromCache(domain, kind string) (string, error) {
	if v, ok := c.plaintextSecretsCache[cacheKey(domain, kind)]; ok {
		return v, nil
	}
	return "", errors.Errorf("No secret for domain %q and type %q", domain, kind)
}

type Required map[string][]string

func Combine(from, to Required) Required {
	for k, v := range from {
		if s, ok := to[k]; ok {
			to[k] = append(s, v...)
		} else {
			to[k] = v
		}
	}
	return to
}

func (c client) DownloadAndDecryptAndCache(ctx context.Context, bucket, dir string, required Required) error {
	b := c.storageClient.Bucket(bucket)
	for domain, kinds := range required {
		for _, kind := range kinds {
			s, err := download(ctx, b, c.FileName(domain, kind), dir)
			if err != nil {
				return errors.Wrap(err, "Failed downloading secret from bucket")
			}
			plaintext, err := c.Decrypt(s)
			if err != nil {
				return errors.Wrap(err, "Failed decrypting secret")
			}
			c.plaintextSecretsCache[cacheKey(domain, kind)] = plaintext
		}
	}
	return nil
}

func (c client) FileName(domain, kind string) string {
	return fmt.Sprintf("%s_%s_cloudkms-%s.json", domain, kind, c.env)
}

func download(ctx context.Context, bucket *storage.BucketHandle, file, dir string) (string, error) {
	fileObject := bucket.Object(dir + file)
	reader, err := fileObject.NewReader(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "Failed opening file %q", file)
	}
	defer reader.Close()
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(reader); err != nil {
		return "", errors.Wrap(err, "Failed reading from reader")
	}
	var s Secret
	if err := json.Unmarshal(buffer.Bytes(), &s); err != nil {
		return "", errors.Wrap(err, "Failed unmarshalling buffer into Secret")
	}
	return s.Ciphertext, nil
}

func cacheKey(domain, kind string) string {
	return fmt.Sprintf("%s%s", domain, kind)
}
