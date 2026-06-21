include .env

.PHONY: server mobile emulator tidy build docker db sqlc orval migrate-up migrate-down migrate-force test-server lint-mobile test-mobile verify dev dev-server dev-mobile doctor doctor-android stop-dev

migrate-up:
	cd ./server/ && infisical run -- go run ./cmd/migrate/main.go up

migrate-down:
	cd ./server/ && infisical run -- go run ./cmd/migrate/main.go down $(steps)

migrate-force:
	cd ./server/ && infisical run -- go run ./cmd/migrate/main.go force $(version)

mobile:
	cd ./mobile/ && infisical run -- nr start

emulator:
	./scripts/dev.sh android

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

db: docker

dev:
	./scripts/dev.sh

dev-server:
	./scripts/dev.sh server

dev-mobile:
	./scripts/dev.sh mobile

doctor:
	./scripts/dev-doctor.sh

doctor-android:
	./scripts/dev-doctor.sh --android

stop-dev:
	./scripts/dev.sh stop

test-server:
	cd ./server/ && go test ./...

lint-mobile:
	cd ./mobile/ && pnpm lint

test-mobile:
	cd ./mobile/ && pnpm test

verify: test-server build lint-mobile test-mobile
