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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	go_errors "github.com/caigwatkin/go/errors"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// Client interface for secrets
type Client interface {
	Value(source, kind string) (string, error)
}

// NewClient returns a client that satisfies the Client interface
func NewClient(ctx context.Context) (Client, error) {
	cloudKMSService, err := newCloudkmsService(ctx)
	if err != nil {
		return nil, go_errors.Wrap(err, "Failed creating cloudkms service")
	}
	s := client{
		cloudKMSService: cloudKMSService,
		secrets:         make(map[string]string),
	}
	return &s, nil
}

type client struct {
	cloudKMSService *cloudkms.Service
	secrets         map[string]string
}

// Value returns the decrypted value of a secret with the given source and kind
func (c *client) Value(source, kind string) (string, error) {
	filename := fmt.Sprintf("%s-%s-cloudkms_%s.json", source, kind, os.Getenv("ENV"))
	p, err := c.value(filename, cloudKMSLoadAndDecrypt)
	if err != nil {
		return "", go_errors.Wrapf(err, "Failed reading value of cloudkms secret %q", filename)
	}
	return p, nil
}

func (c *client) value(filename string, loadAndDecrypt func(s client, filename string) (string, error)) (string, error) {
	if v, ok := c.secrets[filename]; ok {
		return v, nil
	}
	secret, err := loadAndDecrypt(*c, filename)
	if err != nil {
		return "", go_errors.Wrap(err, "Failed decrypting secret")
	}
	c.secrets[filename] = string(secret)
	return c.secrets[filename], nil
}

func load(filename string, dst interface{}) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return go_errors.Wrap(err, "Failed to retrieve current directory")
	}
	f := fmt.Sprintf("%s/%s", dir, filename)
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return go_errors.Wrapf(err, "Failed read file %q", f)
	}
	if err := json.Unmarshal(buf, &dst); err != nil {
		return go_errors.Wrapf(err, "Failed unmarshal file %q into dst", f)
	}
	return nil
}
