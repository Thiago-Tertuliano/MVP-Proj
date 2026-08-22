.PHONY: dev build run migrate-up migrate-down test lint

dev:            ## roda API com hot reload (instale o air: go install github.com/air-verse/air@latest)
	air

build:          ## compila binário
	go build -o bin/api ./cmd/api

run:            ## roda direto
	go run ./cmd/api

migrate-up:     ## aplica migrations via golang-migrate (cmd/migrate)
	go run ./cmd/migrate

migrate-down:   ## reverte a última migration
	go run ./cmd/migrate -down

test:           ## testes
	go test -v -race ./internal/...

lint:           ## linter
	golangci-lint run ./...
