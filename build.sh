#!/bin/sh

go get github.com/aws/aws-lambda-go/events

go build -ldflags="-s -w" -o bin/hello hello/main.go

exit 0