package environment

import (
	"log"
	"os"
	"strconv"

	go_errors "github.com/caigwatkin/go/errors"
)

type Environment struct {
	App         string
	DatabaseUrl string
	Debug       bool
	Remote      bool
	Port        int64
}

// TODO: test
func New(app string) (env Environment, err error) {
	log.Println("Generating environment", app)

	databaseUrl := os.Getenv("DATABASE_URL")

	remote := false
	if osDyno := os.Getenv("DYNO"); osDyno != "" {
		remote = true
	}

	debug := !remote
	if osDebug := os.Getenv("DEBUG"); osDebug != "" {
		debug, err = strconv.ParseBool(osDebug)
		if err != nil {
			err = go_errors.Wrap(err, "Failed parsing environment variable DEBUG")
			return
		}
	}

	port := int64(8080)
	if osPort := os.Getenv("PORT"); osPort != "" {
		port, err = strconv.ParseInt(osPort, 10, 0)
		if err != nil {
			err = go_errors.Wrap(err, "Failed parsing environment variable PORT")
			return
		}
	}

	env = Environment{
		App:         app,
		DatabaseUrl: databaseUrl,
		Debug:       debug,
		Remote:      remote,
		Port:        port,
	}

	log.Println("Generated environment", env)
	return
}
