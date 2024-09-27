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

test:
	go test ./...
.PHONY: test