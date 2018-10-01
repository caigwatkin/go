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
	"encoding/base64"

	go_errors "github.com/caigwatkin/go/errors"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// Encrypt plaintext bytes into a secret
func (c client) Encrypt(plaintext []byte) (*Secret, error) {
	r, err := c.cloudkmsClient.Projects.Locations.KeyRings.CryptoKeys.Encrypt(c.cryptoKey, &cloudkms.EncryptRequest{
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
	resp, err := c.cloudkmsClient.Projects.Locations.KeyRings.CryptoKeys.Decrypt(c.cryptoKey, &cloudkms.DecryptRequest{
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
