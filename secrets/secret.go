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
	"encoding/json"
	"io/ioutil"

	go_errors "github.com/caigwatkin/go/errors"
)

// Secret data model, a subset of properties of a cloudkms secret
type Secret struct {
	Name       string `json:"name,omitempty"`
	Ciphertext string `json:"ciphertext,omitempty"`
}

// SecretFromFile returns a secret from a file
func SecretFromFile(pathToFile string) (*Secret, error) {
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
