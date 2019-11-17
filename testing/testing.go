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

package testing

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Unexpected string
	Desc       string
	At         int
	Input      interface{}
	Expected   interface{}
	Result     interface{}
}

func Errorf(e Error) string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf(`{
	"Unexpected": %q,
	"Desc": %q,
	"At": %d,
	"Input": "potentially unmarshallable",
	"Expected": "potentially unmarshallable",
	"Result": "potentially unmarshallable"
}`, e.Unexpected, e.Desc, e.At)
	}
	return string(b)
}
