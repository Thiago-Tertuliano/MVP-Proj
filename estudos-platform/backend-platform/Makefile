.PHONY: dev build run migrate-up migrate-down test lint content-job content-job-dry

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

content-job-dry: ## valida Courses.md + aula Go sem gravar
	go run ./cmd/content-job --dry-run

content-job:    ## upsert rascunhos no Postgres (não publica)
	go run ./cmd/content-job
