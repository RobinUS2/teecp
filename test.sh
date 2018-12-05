#!/bin/bash
set -e
# go to this file
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

# vet
go vet -v .
cd forwarder && go vet -v . && cd ..
echo "OK vetting"

# test
go test -v -race -cover .
echo "OK tests"