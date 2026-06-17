include .env

.PHONY: server mobile emulator tidy build db sqlc orval migrate-up migrate-down migrate-force

migrate-up:
	cd ./server/ && infisical run -- go run ./cmd/migrate/main.go up

migrate-down:
	cd ./server/ && infisical run -- go run ./cmd/migrate/main.go down $(steps)

migrate-force:
	cd ./server/ && infisical run -- go run ./cmd/migrate/main.go force $(version)

mobile:
	cd ./mobile/ && infisical run -- nr start

emulator:
	$(ANDROID_SDK_ROOT)/emulator/emulator -avd Medium_Phone_API_36.1 &

orval:
	cd ./mobile/ && infisical run -- nr orval

server:
	cd ./server/ && infisical run -- go run ./cmd/api/main.go

sqlc:
	cd ./server/ && go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

tidy:
	cd ./server/ && go mod tidy

build:
	cd ./server/ && go build -o ../server-api ./cmd/api/main.go

docker:
	infisical run -- docker compose up -d
