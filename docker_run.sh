#!/usr/bin/env bash
# Run
docker run -it -v `pwd`:/usr/local/app teecp sh -c 'go get ./... && GOOS=windows GOARCH=amd64 CGO_ENABLED=1  CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_CFLAGS="-I/usr/local/app/include" CGO_LDFLAGS="-L/usr/local/app/lib" go build -o teecp.exe .'
