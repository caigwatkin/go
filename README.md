# go

[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/caigwatkin/go/master/LICENSE)
[![CircleCI](https://circleci.com/gh/caigwatkin/go.svg?style=svg)](https://circleci.com/gh/caigwatkin/go)
[![codecov](https://codecov.io/gh/caigwatkin/go/branch/master/graph/badge.svg)](https://codecov.io/gh/caigwatkin/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/caigwatkin/go)](https://goreportcard.com/report/github.com/caigwatkin/go)

Golang package library for APIs

## Usage

```bash
go get github.com/caigwatkin/go
```

See [github.com/caigwatkin/slate](https://github.com/caigwatkin/slate) for usage in a Go API server.

## CI/CD

Using [CircelCI](https://circleci.com) for builds of commits and pull requests.

All changes are made to branches of `master`. The branch must be up to date with `master` and all commits must be signed with a [GPG key](https://gnupg.org).

The following status checks must pass before merging into master:

- [CircelCI](https://circleci.com) build passes
- [Codecov](https://codecov.io) meets minimum coverage requirements
