package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/config"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres/migration"
)

// Hash bcrypt de "senha1234" gravado pela migration 0004 (autor de demo).
const hashSenhaSeed = "$2a$10$GnWGzZlXrOg9M.dqhPxVR.qVwdW1.MKsO7lVtsv.iupzECQ089Em2"

func main() {
	down := flag.Bool("down", false, "reverte a última migration")
	flag.Parse()

	dsn := config.LoadDB().PostgresURL()

	src, err := iofs.New(migration.FS, ".")
	if err != nil {
		log.Fatalf("source das migrations: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}
	defer func() { _, _ = m.Close() }()

	if *down {
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("última migration revertida")
		return
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations aplicadas")

	if os.Getenv("APP_ENV") == "production" {
		bloquearSenhaDoSeed(dsn)
	}
}

// bloquearSenhaDoSeed impede login com a senha pública do autor de demo em produção.
// Só troca o hash se ainda for o do seed: uma senha definida depois pelo time é preservada.
func bloquearSenhaDoSeed(dsn string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("conexão para bloquear o seed: %v", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	tag, err := conn.Exec(ctx, `UPDATE usuarios SET senha_hash = '!' WHERE email = 'autor.seed@estudos.local' AND senha_hash = $1`, hashSenhaSeed)
	if err != nil {
		log.Fatalf("bloquear senha do seed: %v", err)
	}
	if tag.RowsAffected() > 0 {
		log.Println("senha pública do autor de demo bloqueada (produção)")
	}
}
