package environment

import (
	"log"
	"os"
	"strconv"

	go_errors "github.com/caigwatkin/go/errors"
)

type Environment struct {
	App    string
	Debug  bool
	Remote bool
	Port   int64
}

func New(app string) (env Environment, err error) {
	log.Println("Generating environment", app)

	debug := false
	if osDebug := os.Getenv("DEBUG"); osDebug != "" {
		debug, err = strconv.ParseBool(osDebug)
		if err != nil {
			err = go_errors.Wrap(err, "Failed parsing environment variable DEBUG")
			return
		}
	}

	remote := false
	if osDyno := os.Getenv("DYNO"); osDyno != "" {
		remote = true
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
		App:    app,
		Debug:  debug,
		Remote: remote,
		Port:   port,
	}

	log.Println("Generated environment", env)
	return
}
