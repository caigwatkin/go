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

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	go_context "github.com/caigwatkin/go/context"
	go_errors "github.com/caigwatkin/go/errors"
	go_log "github.com/caigwatkin/go/log"
	go_secrets "github.com/caigwatkin/go/secrets"
)

var (
	plaintext      string
	pathToFile     string
	gcpProjectID   string
	keyRing        string
	key            string
	env            string
	saveAsFileName string
	debug          bool
)

func init() {
	flag.BoolVar(&debug, "debug", true, "Debug mode on/off")
	flag.StringVar(&env, "env", "dev", "Friendly environment name, used for file naming")
	flag.StringVar(&pathToFile, "pathToFile", "", "Path to file to be encrypted. Required if no plaintext given")
	flag.StringVar(&gcpProjectID, "gcpProjectID", "", "GCP project ID which has cloudkms used for encryption")
	flag.StringVar(&key, "key", "", "Cloudkms key to use")
	flag.StringVar(&keyRing, "keyRing", "", "Cloudkms key ring to use")
	flag.StringVar(&plaintext, "plaintext", "", "Plaintext to be encrypted. Required if no pathToFile given")
	flag.StringVar(&saveAsFileName, "saveAsFileName", "", "Optional file name to save as")
	flag.Parse()
}

func main() {
	ctx := go_context.StartUp()

	log.Println("Initialising logger", os.Environ())
	logClient := go_log.NewClient(debug)
	logClient.Info(ctx, "Logger initialised",
		go_log.FmtBool(debug, "debug"),
		go_log.FmtString(env, "env"),
		go_log.FmtString(pathToFile, "pathToFile"),
		go_log.FmtString(gcpProjectID, "gcpProjectID"),
		go_log.FmtString(key, "key"),
		go_log.FmtString(keyRing, "keyRing"),
		go_log.FmtString(plaintext, "plaintext"),
		go_log.FmtString(saveAsFileName, "saveAsFileName"),
		go_log.FmtStrings(os.Environ(), "os.Environ()"),
	)

	logClient.Info(ctx, "Checking required flags")
	if err := checkRequiredFlags(); err != nil {
		logClient.Fatal(ctx, "Failed flag check", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Passed flag check")

	logClient.Info(ctx, "Creating secrets client")
	secretsClient, err := go_secrets.NewClient(ctx, env, gcpProjectID, keyRing, key)
	if err != nil {
		logClient.Fatal(ctx, "Failed creating secrets client", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Created secrets client")

	encrypt(ctx, logClient, secretsClient)
}

func checkRequiredFlags() error {
	if plaintext == "" && pathToFile == "" {
		return go_errors.New("Plaintext or file path must be provided")
	} else if plaintext != "" && pathToFile != "" {
		return go_errors.New("Only plaintext or file path can be encrypted at one time, not both")
	} else if env == "" {
		return go_errors.New("Missing env flag")
	} else if gcpProjectID == "" {
		return go_errors.New("Missing gcpProjectID flag")
	} else if key == "" {
		return go_errors.New("Missing key flag")
	} else if keyRing == "" {
		return go_errors.New("Missing keyRing flag")
	}
	return nil
}

func encrypt(ctx context.Context, logClient go_log.Client, secretsClient go_secrets.Client) {
	if pathToFile != "" {
		buf, err := ioutil.ReadFile(pathToFile)
		if err != nil {
			logClient.Fatal(ctx, "Failed reading file", go_log.FmtError(err))
		}
		plaintext = string(buf)
		logClient.Info(ctx, "Loaded from file", go_log.FmtString(plaintext, "plaintext"))
	}

	secret, err := secretsClient.Encrypt(plaintext)
	if err != nil {
		logClient.Fatal(ctx, "Failed encrypting plaintext", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Encrypted", go_log.FmtAny(secret, "secret"))

	if saveAsFileName != "" {
		saveAs(ctx, logClient, *secret)
	}
}

func saveAs(ctx context.Context, logClient go_log.Client, secret go_secrets.Secret) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logClient.Fatal(ctx, "Failed to get directory of process", go_log.FmtError(err))
	}
	if !strings.HasSuffix(saveAsFileName, ".json") {
		saveAsFileName = fmt.Sprintf("%s.json", saveAsFileName)
	}
	path := fmt.Sprintf("%s/%s", dir, saveAsFileName)
	b, err := json.MarshalIndent(secret, "", "\t")
	if err != nil {
		logClient.Fatal(ctx, "Failed to marshalling secret", go_log.FmtError(err))
	}
	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		logClient.Fatal(ctx, "Failed to save file", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Saved", go_log.FmtString(path, "path"))
}
