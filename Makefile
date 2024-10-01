NAME := $(shell basename $(CURDIR))

fmt:
	go fmt ./...
.PHONY: fmt

lint: fmt
	golint ./...
.PHONY: lint

vet: fmt
	go vet ./...
.PHONY: vet

build: 
	go build .
.PHONY: build

run:
	go run .
.PHONY: run

clean:
	rm -f $(NAME) *.out *.html

test:
	go test ./...
.PHONY: test

cover:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out
.PHONY: cover