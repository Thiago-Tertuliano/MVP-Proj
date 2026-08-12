package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/usecase"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/config"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/external"
	pgrepo "github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/handler"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/middleware"
)

// New monta o router e o grafo de dependências (composição root).
func New(cfg *config.Config, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// ---- repositórios (infra) ----
	usuarioRepo := pgrepo.NewUsuarioRepoPG(pool)
	refreshRepo := pgrepo.NewRefreshTokenRepoPG(pool)
	artigoRepo := pgrepo.NewArtigoRepoPG(pool)

	// ---- ports concretos (infra) ----
	hasher := external.NewBcryptHasher(10)
	tokens := external.NewJWTService(cfg.JWTSecret)

	// ---- use cases (aplicação) ----
	tokenCfg := usecase.TokenConfig{
		AccessTTLMin:  cfg.JWTAccessTTLMin,
		RefreshTTLHor: cfg.JWTRefreshTTLHours,
	}
	registrarUC := usecase.NewRegistrarUsuario(usuarioRepo, hasher, tokens, usecase.RegistrarConfig{
		AccessTTLMin: cfg.JWTAccessTTLMin, RefreshTTLHor: cfg.JWTRefreshTTLHours,
	})
	loginUC := usecase.NewLoginUsuario(usuarioRepo, refreshRepo, hasher, tokens, tokenCfg)
	refreshUC := usecase.NewRefreshTokenUC(refreshRepo, usuarioRepo, tokens, tokenCfg)
	perfilUC := usecase.NewObterPerfil(usuarioRepo)

	criarArtigoUC := usecase.NewCriarArtigo(artigoRepo)
	obterArtigoUC := usecase.NewObterArtigo(artigoRepo)
	listarArtigosUC := usecase.NewListarArtigos(artigoRepo)
	atualizarArtigoUC := usecase.NewAtualizarArtigo(artigoRepo)
	publicarArtigoUC := usecase.NewPublicarArtigo(artigoRepo)

	// ---- handlers (apresentação) ----
	auth := handler.NewAuthHandler(registrarUC, loginUC, refreshUC, perfilUC)
	artigos := handler.NewArtigoHandler(criarArtigoUC, obterArtigoUC, listarArtigosUC, atualizarArtigoUC, publicarArtigoUC)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/registrar", auth.Registrar)
		api.Post("/auth/login", auth.Login)
		api.Post("/auth/refresh", auth.Refresh)

		api.Get("/artigos", artigos.Listar)
		api.Get("/artigos/{slug}", artigos.Obter)

		api.Group(func(pr chi.Router) {
			pr.Use(middleware.NewAutenticador(tokens).Proteger)
			pr.Get("/auth/me", auth.Me)
			pr.Post("/artigos", artigos.Criar)
			pr.Put("/artigos/{id}", artigos.Atualizar)
			pr.Post("/artigos/{id}/publicar", artigos.Publicar)
		})
	})

	return r
}
