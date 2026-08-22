package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/usecase"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/external"
	pgrepo "github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres/repository"
)

func TestIntegration_RegistrarECriarArtigo(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") != "1" {
		t.Skip("defina RUN_INTEGRATION=1 com Postgres no ar (docker compose up -d)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := "postgres://estudos:estudos_dev@localhost:5433/estudos_platform?sslmode=disable"
	if v := os.Getenv("DATABASE_URL"); v != "" {
		dsn = v
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("postgres indisponível: %v", err)
	}

	usuarios := pgrepo.NewUsuarioRepoPG(pool)
	refresh := pgrepo.NewRefreshTokenRepoPG(pool)
	artigos := pgrepo.NewArtigoRepoPG(pool)
	hasher := external.NewBcryptHasher(10)
	tokens := external.NewJWTService("integration-secret-min-32-bytes!!")

	email := "it-" + uuid.New().String()[:8] + "@estudos.local"
	reg := usecase.NewRegistrarUsuario(usuarios, refresh, hasher, tokens, usecase.RegistrarConfig{
		AccessTTLMin: 15, RefreshTTLHor: 24,
	})
	auth, err := reg.Execute(ctx, dto.RegistrarRequest{Nome: "IT User", Email: email, Senha: "senha1234"})
	if err != nil {
		t.Fatalf("registrar: %v", err)
	}

	criar := usecase.NewCriarArtigo(artigos)
	artigo, err := criar.Execute(ctx, dto.CriarArtigoRequest{
		Titulo:   "Artigo de integração",
		Conteudo: []byte(`{"blocks":[{"type":"p","text":"ok"}]}`),
	}, auth.Usuario.ID)
	if err != nil {
		t.Fatalf("criar artigo: %v", err)
	}
	if artigo.Slug == "" || artigo.Status != "rascunho" {
		t.Fatalf("artigo inesperado: %+v", artigo)
	}
}
