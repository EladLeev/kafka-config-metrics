SHELL := /bin/bash
export GOBIN := $(CWD)/.bin
NAME=kafka-config-metrics

build:
	GOARCH=amd64 GOOS=darwin go build -o ${NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${NAME}-linux main.go

clean:
	go clean
	rm ${NAME}-darwin
	rm ${NAME}-linux

test:
	go test -v ./...

test_coverage:
	go test -v ./... -coverprofile=coverage.out

dep:
	go mod download

tidy:
	go mod tidy

vet:
	go vet
