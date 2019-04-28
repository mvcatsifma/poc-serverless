.PHONY: build clean deploy

build:
	mkdir -p bin
	go get ./...
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/db db/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
