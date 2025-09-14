.PHONY: test coverage lint

all: build lint test coverage

build:
	go build cmd/product-showcase/product-showcase.go
lint:
	go fmt ./...
test:
	go test ./...
coverage:
	go test -v -cover -coverprofile=coverage.out ./... &&\
	go tool cover -html=coverage.out -o coverage.html