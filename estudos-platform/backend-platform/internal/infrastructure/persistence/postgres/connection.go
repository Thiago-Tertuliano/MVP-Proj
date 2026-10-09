package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/config"
)

func NewConnection(cfg *config.Config) *pgxpool.Pool {
	dsn := cfg.PostgresURL()

	// Postgres gerenciado gratuito pode estar "dormindo": o 1º acesso demora alguns segundos.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	slog.Info("conexão com o banco estabelecida")
	return pool
}
