package schema

import (
	"context"
	"fmt"
	"net/http"

	go_environment "github.com/caigwatkin/go/environment"
	go_errors "github.com/caigwatkin/go/errors"
	go_log "github.com/caigwatkin/go/log"
	"github.com/xeipuuv/gojsonschema"
)

type Client interface {
	Validate(ctx context.Context, schemaFileName string, bytes []byte) error
}

type Config struct {
	Env go_environment.Environment
}

func NewClient(ctx context.Context, config Config, logClient go_log.Client, schemaFileNames []string) (Client, error) {
	logClient.Info(ctx, "Initializing", go_log.FmtAny(config, "config"), go_log.FmtStrings(schemaFileNames, "schemaFileNames"))

	type schemaAndFileNameAndError struct {
		Err      error
		FileName string
		Schema   *gojsonschema.Schema
	}
	ch := make(chan schemaAndFileNameAndError, len(schemaFileNames))
	for _, schemaFileName := range schemaFileNames {
		go func(schemaFileName string) {
			logClient.Info(ctx, "Loading", go_log.FmtString("schemaFileName", schemaFileName))
			schema, err := loadSchemaFromFile(schemaFileName, config.Env.WorkingDirectory)
			if err != nil {
				ch <- schemaAndFileNameAndError{
					Err: go_errors.Wrapf(err, "Failed to load schema %q from file", schemaFileName),
				}
				return
			}
			ch <- schemaAndFileNameAndError{
				FileName: schemaFileName,
				Schema:   schema,
			}
			logClient.Info(ctx, "Loaded", go_log.FmtString("schemaFileName", schemaFileName))
		}(schemaFileName)
	}

	schemaByFileName := make(map[string]*gojsonschema.Schema)

	var err error
	for range schemaFileNames {
		schemaAndFileNameAndError := <-ch
		if schemaAndFileNameAndError.Err != nil {
			logClient.Error(ctx, "Failed to load schema from file in goroutine, will check others and return an error", go_log.FmtError(schemaAndFileNameAndError.Err))
			err = go_errors.Wrap(schemaAndFileNameAndError.Err, "Failed to load schema from file in goroutine")
		}
		schemaByFileName[schemaAndFileNameAndError.FileName] = schemaAndFileNameAndError.Schema
	}
	if err != nil {
		return nil, err
	}

	logClient.Info(ctx, "Initialized")
	return &client{
		config:           config,
		logClient:        logClient,
		schemaByFileName: schemaByFileName,
	}, nil
}

type client struct {
	config           Config
	schemaByFileName map[string]*gojsonschema.Schema
	logClient        go_log.Client
}

func loadSchemaFromFile(fileName, workingDirectory string) (*gojsonschema.Schema, error) {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/%s", workingDirectory, fileName)))
	if err != nil {
		return nil, go_errors.Wrapf(err, "Failed to initialize schema for file %q", fileName)
	}

	return schema, nil
}

func (c client) Validate(ctx context.Context, schemaFileName string, bytes []byte) error {
	c.logClient.Info(ctx, "Validating", go_log.FmtString(schemaFileName, "schemaFileName"), go_log.FmtInt(len(bytes), "len(bytes)"))

	if len(bytes) == 0 {
		return go_errors.NewStatus(http.StatusBadRequest, "BODY_MUST_EXIST")
	}

	schema, ok := c.schemaByFileName[schemaFileName]
	if !ok {
		return go_errors.Errorf("Schema %q not loaded", schemaFileName)
	}

	result, err := schema.Validate(gojsonschema.NewBytesLoader(bytes))
	if err != nil {
		return go_errors.NewStatus(http.StatusBadRequest, err.Error())
	}

	if !result.Valid() {
		var errorItems []go_errors.Item

		for _, resultError := range result.Errors() {
			errorItems = append(errorItems, go_errors.Item{
				Message: resultError.Description(),
				Field:   resultError.Field(),
			})
		}

		return go_errors.NewStatusWithItems(http.StatusBadRequest, "Failed schema validation", errorItems)
	}

	c.logClient.Info(ctx, "Validated")
	return nil
}
