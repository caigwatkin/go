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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	go_context "github.com/caigwatkin/go/context"
	go_errors "github.com/caigwatkin/go/errors"
	go_log "github.com/caigwatkin/go/log"
	go_secrets "github.com/caigwatkin/go/secrets"
)

var (
	ciphertext         string
	debug              bool
	env                string
	gcpProjectID       string
	key                string
	keyRing            string
	pathToFile         string
	saveAsSecretType   string
	saveAsSecretDomain string
)

func init() {
	flag.StringVar(&ciphertext, "ciphertext", "", "Ciphertext to be decrypted. Required if no pathToFile given")
	flag.BoolVar(&debug, "debug", true, "Debug mode on/off")
	flag.StringVar(&env, "env", "dev", "Friendly environment name, used for file naming")
	flag.StringVar(&pathToFile, "pathToFile", "", "Path to file to be decrypted. Required if no ciphertext given")
	flag.StringVar(&gcpProjectID, "gcpProjectID", "", "GCP project ID which has cloudkms used for decryption")
	flag.StringVar(&key, "key", "", "Cloudkms key to use")
	flag.StringVar(&keyRing, "keyRing", "", "Cloudkms key ring to use")
	flag.StringVar(&saveAsSecretDomain, "saveAsSecretDomain", "", "Optional secret domain to use as file name for saving, must be provided if saveAsSecretType provided")
	flag.StringVar(&saveAsSecretType, "saveAsSecretType", "", "Optional secret type to use as file name for saving, must be provided if saveAsSecretDomain provided")
	flag.Parse()
}

func main() {
	ctx := go_context.StartUp()

	log.Println("Initialising logger", os.Environ())
	logClient := go_log.NewClient(debug)
	logClient.Info(ctx, "Logger initialised",
		go_log.FmtString(ciphertext, "ciphertext"),
		go_log.FmtBool(debug, "debug"),
		go_log.FmtString(env, "env"),
		go_log.FmtString(pathToFile, "pathToFile"),
		go_log.FmtString(gcpProjectID, "gcpProjectID"),
		go_log.FmtString(key, "key"),
		go_log.FmtString(keyRing, "keyRing"),
		go_log.FmtString(saveAsSecretType, "saveAsSecretType"),
		go_log.FmtString(saveAsSecretDomain, "saveAsSecretDomain"),
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

	decrypt(ctx, logClient, secretsClient)
}

func checkRequiredFlags() error {
	if (ciphertext != "") == (pathToFile != "") {
		return go_errors.New("Either `ciphertext` or `pathToFile` flag values must be provided, not both")
	} else if (saveAsSecretDomain != "") != (saveAsSecretType != "") {
		return go_errors.New("Both or neither `saveAsSecretDomain` and `saveAsSecretType` flag values must be provided")
	} else if env == "" {
		return go_errors.New("Missing `env` flag value")
	} else if gcpProjectID == "" {
		return go_errors.New("Missing `gcpProjectID` flag value")
	} else if key == "" {
		return go_errors.New("Missing `key` flag value")
	} else if keyRing == "" {
		return go_errors.New("Missing `keyRing` flag value")
	}
	return nil
}

func decrypt(ctx context.Context, logClient go_log.Client, secretsClient go_secrets.Client) {
	if pathToFile != "" {
		secret, err := go_secrets.SecretFromFile(pathToFile)
		if err != nil {
			logClient.Fatal(ctx, "Failed reading secret from file", go_log.FmtError(err))
		}
		logClient.Info(ctx, "Loaded from file", go_log.FmtAny(secret, "secret"))
		ciphertext = secret.Ciphertext
		logClient.Info(ctx, "Loaded", go_log.FmtString(ciphertext, "ciphertext"))
	}

	plaintext, err := secretsClient.Decrypt(ciphertext)
	if err != nil {
		logClient.Fatal(ctx, "Failed decrypting ciphertext", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Decrypted", go_log.FmtAny(plaintext, "plaintext"))

	if saveAsSecretType != "" {
		saveAs(ctx, logClient, []byte(plaintext))
	}
}

func saveAs(ctx context.Context, logClient go_log.Client, plaintext []byte) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logClient.Fatal(ctx, "Failed to get directory of process", go_log.FmtError(err))
	}
	path := fmt.Sprintf("%s/%s_%s_plaintext.json", dir, saveAsSecretDomain, saveAsSecretType)
	if err := ioutil.WriteFile(path, plaintext, 0644); err != nil {
		logClient.Fatal(ctx, "Failed to save file", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Saved", go_log.FmtString(path, "path"))
}
