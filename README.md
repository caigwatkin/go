# go

[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/caigwatkin/go/master/LICENSE)
[![Build Status](https://travis-ci.org/caigwatkin/go.svg?branch=master)](https://travis-ci.org/caigwatkin/go)
[![codecov](https://codecov.io/gh/caigwatkin/go/branch/master/graph/badge.svg)](https://codecov.io/gh/caigwatkin/go)
[![GolangCI](https://golangci.com/badges/github.com/caigwatkin/go.svg)](https://golangci.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/caigwatkin/go)](https://goreportcard.com/report/github.com/caigwatkin/go)

Golang package library for APIs

## Usage

```bash
go get -u github.com/caigwatkin/go
```

See [github.com/caigwatkin/slate](https://github.com/caigwatkin/slate) for usage in a Go API server.

## CI/CD

Using [Travis CI](https://travis-ci.org) for builds of commits and pull requests.

All changes are made to branches of `master`. The branch must be up to date with `master` and all commits must be signed with a [GPG key](https://gnupg.org).

The following status checks must pass before merging into master:

- [Travis CI](https://travis-ci.org) build passes
- [Codecov](https://codecov.io) meets minimum coverage requirements
- [GolangCI](https://golangci.com) finds no issues

## Dependency management

Using [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules)
