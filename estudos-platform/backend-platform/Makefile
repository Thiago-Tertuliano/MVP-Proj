.PHONY: dev build run migrate-up migrate-down test

dev:            ## roda API com hot reload (instale o air: go install github.com/air-verse/air@latest)
	air

build:          ## compila binário
	go build -o bin/api ./cmd/api

run:            ## roda direto
	go run ./cmd/api

migrate-up:     ## aplica migrations
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0001_create_usuario.up.sql

migrate-down:   ## reverte migrations
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0001_create_usuario.down.sql

test:           ## testes
	go test -v -race ./internal/...

lint:           ## linter
	golangci-lint run ./...