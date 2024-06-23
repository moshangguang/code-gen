#!/bin/bash

rm -rf ./bin/darwin/code-gen
GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin/code-gen main.go
