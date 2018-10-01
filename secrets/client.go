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

	"cloud.google.com/go/storage"
	"github.com/caigwatkin/go/errors"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type Client interface {
	Decrypt(secret Secret) ([]byte, error)
	DownloadAndDecryptAndCache(ctx context.Context, bucket, dir string, required Required) error
	Encrypt(secret []byte) (*Secret, error)
	FileName(domain, kind string) string
	Secret(domain, kind string) ([]byte, error)
}

type client struct {
	cloudkmsClient *cloudkms.Service
	env            string
	gcpProjectID   string
	keyRing        string
	key            string
	secrets        map[string][]byte
	storageClient  *storage.Client
}

func NewClient(ctx context.Context, env, gcpProjectID, keyRing, key string) (Client, error) {
	googleClient, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, errors.Wrap(err, "Failed initializing google client")
	}
	cloudkmsClient, err := cloudkms.New(googleClient)
	if err != nil {
		return nil, errors.Wrap(err, "Failed initializing cloudkms client")
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed initializing storage client")
	}
	return client{
		cloudkmsClient: cloudkmsClient,
		env:            env,
		gcpProjectID:   gcpProjectID,
		keyRing:        keyRing,
		key:            key,
		secrets:        make(map[string][]byte),
		storageClient:  storageClient,
	}, nil
}
