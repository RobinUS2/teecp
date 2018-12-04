#!/bin/bash
set -e
# go to this file
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

# to src
cd src

# vet
go vet -v .
echo "OK vetting"

# test
go test -v -race -cover .
echo "OK tests"