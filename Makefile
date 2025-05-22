SHELL:=/bin/bash

install-dev:
	mkdir -p ignored
	echo '*' > ignored/.gitignore
	touch ignored/.env

run-server:
	source ignored/.env && go run cmd/server/*.go

run-agent:
	source ignored/.env && go run cmd/agent/*.go

build-server-amd64-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/whispersync-server.exe ./cmd/server
build-server-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/whispersync-server-linux-amd64 ./cmd/server
build-server-amd64-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/whispersync-server-macos-amd64 ./cmd/server
build-server-arm64-macos:
	GOOS=darwin GOARCH=arm64 go build -o bin/whispersync-server-macos-arm64 ./cmd/server

build-agent-amd64-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/whispersync-agent.exe ./cmd/agent
build-agent-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/whispersync-agent-linux-amd64 ./cmd/agent
build-agent-amd64-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/whispersync-agent-macos-amd64 ./cmd/agent
build-agent-arm64-macos:
	GOOS=darwin GOARCH=arm64 go build -o bin/whispersync-agent-macos-arm64 ./cmd/agent

