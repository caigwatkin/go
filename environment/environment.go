package environment

import (
	"log"
	"os"
	"strconv"

	go_errors "github.com/caigwatkin/go/errors"
)

type Environment struct {
	App              string
	DatabaseUrl      string
	Debug            bool
	Remote           bool
	Port             int64
	WorkingDirectory string
}

// TODO: test
func New(app string) (env Environment, err error) {
	log.Println("Generating environment", app)

	databaseUrl := os.Getenv("DATABASE_URL")

	remote := os.Getenv("REMOTE") != "" ||
		os.Getenv("DYNO") != ""

	debug := !remote
	if osDebug := os.Getenv("DEBUG"); osDebug != "" {
		debug, err = strconv.ParseBool(osDebug)
		if err != nil {
			err = go_errors.Wrap(err, "Failed to parse environment variable DEBUG")
			return
		}
	}

	port := int64(8080)
	if osPort := os.Getenv("PORT"); osPort != "" {
		port, err = strconv.ParseInt(osPort, 10, 0)
		if err != nil {
			err = go_errors.Wrap(err, "Failed to parse environment variable PORT")
			return
		}
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		err = go_errors.Wrap(err, "Failed to get working directory")
		return
	}

	env = Environment{
		App:              app,
		DatabaseUrl:      databaseUrl,
		Debug:            debug,
		Remote:           remote,
		Port:             port,
		WorkingDirectory: workingDirectory,
	}

	log.Println("Generated environment", env)
	return
}
