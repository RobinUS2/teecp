#!/bin/bash
set -e
cd src
go vet -v .
go test -v -race -cover .
