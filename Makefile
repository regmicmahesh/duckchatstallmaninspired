.PHONY: server client
export ENV ?= dev

server:
	@go run ./cmd/server

client:
	@go run ./cmd/client localhost:8080

build: build-server build-client

build-server:
	@go build ./cmd/server

build-client:
	@go build ./cmd/client
