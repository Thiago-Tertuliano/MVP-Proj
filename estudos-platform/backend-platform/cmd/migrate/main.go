package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres/migration"
)

func main() {
	down := flag.Bool("down", false, "reverte a última migration")
	flag.Parse()

	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		env("DB_USER", "estudos"),
		env("DB_PASSWORD", "estudos_dev"),
		env("DB_HOST", "localhost"),
		env("DB_PORT", "5433"),
		env("DB_NAME", "estudos_platform"),
		env("DB_SSL_MODE", "disable"),
	)

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
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
