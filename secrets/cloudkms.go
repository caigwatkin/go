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
	"context"
	"encoding/base64"
	"strings"

	go_errors "github.com/caigwatkin/go/errors"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func newCloudkmsService(ctx context.Context) (*cloudkms.Service, error) {
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed initialising default google client cloudkms.CloudPlatformScope")
	}
	cloudKMSService, err := cloudkms.New(client)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed creating the kms client")
	}
	return cloudKMSService, nil
}

type cloudKMSSecret struct {
	Name       string `json:"name"`
	Ciphertext string `json:"ciphertext"`
}

func cloudKMSLoadAndDecrypt(c client, filename string) (string, error) {
	var s cloudKMSSecret
	err := load(filename, &s)
	if err != nil {
		return "", go_errors.Wrap(err, "Failed loading secret from file")
	}
	n := strings.Split(s.Name, "/")
	parentName := strings.Join(n[:8], "/")
	req := &cloudkms.DecryptRequest{
		Ciphertext: s.Ciphertext,
	}
	resp, err := c.cloudKMSService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parentName, req).Do()
	if err != nil {
		return "", err
	}
	p, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return "", go_errors.Wrap(err, "Failed decoding base64 plaintext to byte array")
	}
	return strings.Replace(string(p), "\n", "", -1), nil
}
