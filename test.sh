#!/bin/bash
cd src
go test -v --race -cover .
