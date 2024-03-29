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
	go_environment "github.com/caigwatkin/go/environment"
	go_errors "github.com/caigwatkin/go/errors"
	go_log "github.com/caigwatkin/go/log"
	go_secrets "github.com/caigwatkin/go/secrets"
)

var (
	ciphertext         string
	cloudkmsKey        string
	cloudkmsKeyRing    string
	env                string
	gcpProjectId       string
	pathToFile         string
	saveAsFileType     string
	saveAsSecretType   string
	saveAsSecretDomain string
)

func init() {
	flag.StringVar(&ciphertext, "ciphertext", "", "Ciphertext to be decrypted. Required if no `pathToFile` given")
	flag.StringVar(&cloudkmsKey, "cloudkmsKey", "", "Cloud KMS key to use")
	flag.StringVar(&cloudkmsKeyRing, "cloudkmsKeyRing", "", "Cloud KMS key ring to use")
	flag.StringVar(&env, "env", "dev", "Friendly environment name, used for file naming")
	flag.StringVar(&pathToFile, "pathToFile", "", "Path to file to be decrypted. Required if no `ciphertext` given")
	flag.StringVar(&gcpProjectId, "gcpProjectId", "", "GCP project ID which has Cloud KMS used for decryption")
	flag.StringVar(&saveAsFileType, "saveAsFileType", "json", "Optional file type to use as file name for saving")
	flag.StringVar(&saveAsSecretDomain, "saveAsSecretDomain", "", "Optional secret domain to use as file name for saving, must be provided if `saveAsSecretType` provided")
	flag.StringVar(&saveAsSecretType, "saveAsSecretType", "", "Optional secret type to use as file name for saving, must be provided if `saveAsSecretDomain` provided")
	flag.Parse()
}

func main() {
	osEnviron := os.Environ()
	log.Println("Starting", osEnviron)

	environment, err := go_environment.New("Decrypt")
	if err != nil {
		log.Fatal("Failed generating new environment", err)
	}

	ctx := go_context.StartUp()

	logClient := go_log.NewClient(ctx, go_log.Config{
		Env: environment,
	})

	logClient.Info(ctx, "Starting",
		go_log.FmtString(ciphertext, "ciphertext"),
		go_log.FmtString(env, "env"),
		go_log.FmtString(pathToFile, "pathToFile"),
		go_log.FmtString(gcpProjectId, "gcpProjectId"),
		go_log.FmtString(cloudkmsKey, "cloudkmsKey"),
		go_log.FmtString(cloudkmsKeyRing, "cloudkmsKeyRing"),
		go_log.FmtString(saveAsSecretType, "saveAsSecretType"),
		go_log.FmtString(saveAsSecretDomain, "saveAsSecretDomain"),
		go_log.FmtStrings(osEnviron, "osEnviron"),
	)

	logClient.Info(ctx, "Checking required flags")
	if err := checkRequiredFlags(); err != nil {
		logClient.Fatal(ctx, "Failed flag check", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Passed flag check")

	secretsClient, err := go_secrets.NewClient(ctx, go_secrets.Config{
		CloudkmsKey:     cloudkmsKey,
		CloudkmsKeyRing: cloudkmsKeyRing,
		Env:             env,
		GcpProjectId:    gcpProjectId,
	}, logClient)
	if err != nil {
		logClient.Fatal(ctx, "Failed creating secrets client", go_log.FmtError(err))
	}

	decrypt(ctx, logClient, secretsClient)
}

func checkRequiredFlags() error {
	if (ciphertext != "") == (pathToFile != "") {
		return go_errors.New("Either `ciphertext` or `pathToFile` flag values must be provided, not both")
	} else if (saveAsSecretDomain != "") != (saveAsSecretType != "") {
		return go_errors.New("Both or neither `saveAsSecretDomain` and `saveAsSecretType` flag values must be provided")
	} else if env == "" {
		return go_errors.New("Missing `env` flag value")
	} else if gcpProjectId == "" {
		return go_errors.New("Missing `gcpProjectId` flag value")
	} else if cloudkmsKey == "" {
		return go_errors.New("Missing `cloudkmsKey` flag value")
	} else if cloudkmsKeyRing == "" {
		return go_errors.New("Missing `cloudkmsKeyRing` flag value")
	}
	return nil
}

func decrypt(ctx context.Context, logClient go_log.Client, secretsClient go_secrets.Client) {
	secret := go_secrets.Secret{
		Ciphertext: ciphertext,
	}
	if pathToFile != "" {
		s, err := secretsClient.SecretFromFile(pathToFile)
		if err != nil {
			logClient.Fatal(ctx, "Failed reading secret from file", go_log.FmtError(err))
		}
		secret = *s
		logClient.Info(ctx, "Loaded from file", go_log.FmtAny(secret, "secret"))
	}

	logClient.Info(ctx, "Decrypting", go_log.FmtAny(secret, "secret"))
	plaintext, err := secretsClient.Decrypt(secret)
	if err != nil {
		logClient.Fatal(ctx, "Failed decrypting ciphertext", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Decrypted", go_log.FmtBytes(plaintext, "plaintext"))

	if saveAsSecretType != "" {
		saveAs(ctx, logClient, plaintext)
	}
}

func saveAs(ctx context.Context, logClient go_log.Client, plaintext []byte) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logClient.Fatal(ctx, "Failed to get directory of process", go_log.FmtError(err))
	}
	path := fmt.Sprintf("%s/%s_%s_plaintext.%s", dir, saveAsSecretDomain, saveAsSecretType, saveAsFileType)
	if err := ioutil.WriteFile(path, plaintext, 0644); err != nil {
		logClient.Fatal(ctx, "Failed to save file", go_log.FmtError(err))
	}
	logClient.Info(ctx, "Saved", go_log.FmtString(path, "path"))
}
