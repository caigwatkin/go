#!/bin/bash -eu

# Run from repo root:
# $ ./scripts/go/start.sh [cmd]

SCRIPT_DIR="$( cd "$(dirname "$0")" ; pwd -P )"
SCRIPT_DIR_NAME=${SCRIPT_DIR##*/}
SCRIPT_NAME=`basename $0`

API="github.com/caigwatkin/go"

export GO111MODULE=on

echo "${SCRIPT_NAME} -> START at `date '+%Y-%m-%d %H:%M:%S'`..."

echo "Formatting..."
go fmt ${API}/...
echo

echo "Vetting..."
go vet -mod vendor ${API}/...
echo

echo "Linting..."
revive -formatter stylish -config ${SCRIPT_DIR}/../../configs/revive.toml -exclude ${SCRIPT_DIR}/../../vendor/... ${SCRIPT_DIR}/...
echo

echo "Testing..."
go test -mod vendor ${API}/...
echo
