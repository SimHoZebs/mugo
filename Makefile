include .env

.PHONY: server mobile emulator tidy build db sqlc adk orval

mobile:
	cd ./mobile/ && infisical run -- nr start

emulator:
	$(ANDROID_SDK_ROOT)/emulator/emulator -avd Medium_Phone_API_36.1 &

orval:
	cd ./mobile/ && infisical run -- nr orval

server:
	cd ./server/ && infisical run -- go run ./cmd/api/main.go

adk:
	cd ./server/ && infisical run -- go run ./cmd/adk/main.go web api webui

adk-api:
	cd ./server/ && infisical run -- go run ./cmd/adk/main.go web api

adk-help:
	cd ./server/ && infisical run -- go run ./cmd/adk/main.go --help

sqlc:
	cd ./server/ && go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

tidy:
	cd ./server/ && go mod tidy

build:
	cd ./server/ && go build -o ../server-api ./cmd/api/main.go

docker:
	infisical run -- docker compose up -d
