.PHONY: dev build run migrate-up migrate-down test lint

dev:            ## roda API com hot reload (instale o air: go install github.com/air-verse/air@latest)
	air

build:          ## compila binário
	go build -o bin/api ./cmd/api

run:            ## roda direto
	go run ./cmd/api

migrate-up:     ## aplica migrations 0001 + 0002 + 0003
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0001_create_usuario.up.sql
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0002_create_artigos.up.sql
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0003_create_trilhas.up.sql

migrate-down:   ## reverte migrations (ordem inversa)
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0003_create_trilhas.down.sql
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0002_create_artigos.down.sql
	docker exec -i estudos-postgres psql -U estudos -d estudos_platform \
	  < internal/infrastructure/persistence/postgres/migration/0001_create_usuario.down.sql

test:           ## testes
	go test -v -race ./internal/...

lint:           ## linter
	golangci-lint run ./...
