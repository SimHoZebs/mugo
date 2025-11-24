.PHONY: server mobile

mobile:
	cd ./mobile/ && infisical run -- nr start

orval:
	cd ./mobile/ && infisical run -- nr orval

server:
	cd ./server/ && infisical run -- go run ./cmd/api/main.go
