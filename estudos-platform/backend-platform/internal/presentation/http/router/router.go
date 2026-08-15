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
	r.Use(middleware.NewCORS(cfg.CORSAllowedOrigins).Handler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	usuarioRepo := pgrepo.NewUsuarioRepoPG(pool)
	refreshRepo := pgrepo.NewRefreshTokenRepoPG(pool)
	artigoRepo := pgrepo.NewArtigoRepoPG(pool)
	trilhaRepo := pgrepo.NewTrilhaRepoPG(pool)

	hasher := external.NewBcryptHasher(10)
	tokens := external.NewJWTService(cfg.JWTSecret)

	tokenCfg := usecase.TokenConfig{
		AccessTTLMin:  cfg.JWTAccessTTLMin,
		RefreshTTLHor: cfg.JWTRefreshTTLHours,
	}
	registrarUC := usecase.NewRegistrarUsuario(usuarioRepo, refreshRepo, hasher, tokens, usecase.RegistrarConfig{
		AccessTTLMin: cfg.JWTAccessTTLMin, RefreshTTLHor: cfg.JWTRefreshTTLHours,
	})
	loginUC := usecase.NewLoginUsuario(usuarioRepo, refreshRepo, hasher, tokens, tokenCfg)
	refreshUC := usecase.NewRefreshTokenUC(refreshRepo, usuarioRepo, tokens, tokenCfg)
	perfilUC := usecase.NewObterPerfil(usuarioRepo)
	logoutUC := usecase.NewLogoutUsuario(refreshRepo)

	criarArtigoUC := usecase.NewCriarArtigo(artigoRepo)
	obterArtigoUC := usecase.NewObterArtigo(artigoRepo)
	listarArtigosUC := usecase.NewListarArtigos(artigoRepo)
	atualizarArtigoUC := usecase.NewAtualizarArtigo(artigoRepo)
	publicarArtigoUC := usecase.NewPublicarArtigo(artigoRepo)

	criarTrilhaUC := usecase.NewCriarTrilha(trilhaRepo)
	obterTrilhaUC := usecase.NewObterTrilha(trilhaRepo)
	listarTrilhasUC := usecase.NewListarTrilhas(trilhaRepo)
	adicionarModuloUC := usecase.NewAdicionarModulo(trilhaRepo)
	publicarTrilhaUC := usecase.NewPublicarTrilha(trilhaRepo)

	auth := handler.NewAuthHandler(registrarUC, loginUC, refreshUC, perfilUC, logoutUC)
	artigos := handler.NewArtigoHandler(criarArtigoUC, obterArtigoUC, listarArtigosUC, atualizarArtigoUC, publicarArtigoUC)
	trilhas := handler.NewTrilhaHandler(criarTrilhaUC, obterTrilhaUC, listarTrilhasUC, adicionarModuloUC, publicarTrilhaUC)

	r.Route("/api/v1", func(api chi.Router) {
		api.Post("/auth/registrar", auth.Registrar)
		api.Post("/auth/login", auth.Login)
		api.Post("/auth/refresh", auth.Refresh)

		api.Get("/artigos", artigos.Listar)
		api.Get("/artigos/{slug}", artigos.Obter)
		api.Get("/trilhas", trilhas.Listar)
		api.Get("/trilhas/{slug}", trilhas.Obter)

		api.Group(func(pr chi.Router) {
			pr.Use(middleware.NewAutenticador(tokens).Proteger)
			pr.Get("/auth/me", auth.Me)
			pr.Post("/auth/logout", auth.Logout)
			pr.Post("/artigos", artigos.Criar)
			pr.Put("/artigos/{id}", artigos.Atualizar)
			pr.Post("/artigos/{id}/publicar", artigos.Publicar)
			pr.Post("/trilhas", trilhas.Criar)
			pr.Post("/trilhas/{id}/modulos", trilhas.AdicionarModulo)
			pr.Post("/trilhas/{id}/publicar", trilhas.Publicar)
		})
	})

	return r
}
