package repository

import (
	"context"
	"time"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
)

// RoadmapResumo é a linha do catálogo (sem o grafo).
type RoadmapResumo struct {
	ID           string
	Slug         string
	Titulo       string
	Descricao    string
	Icone        string
	Publicado    bool
	AutorID      string
	TotalNos     int
	TotalChefes  int
	XPTotal      int
	AtualizadoEm time.Time
}

// NoVinculado aponta um nó "concluído por artigo" de um roadmap publicado.
type NoVinculado struct {
	RoadmapSlug string
	NoID        string
}

type RoadmapRepository interface {
	// Save grava cabeçalho, nós e arestas de forma atômica (substitui o grafo).
	Save(ctx context.Context, r *entity.Roadmap) error
	// FindBySlug devolve o roadmap com o grafo completo (inclui rascunhos) ou NotFound.
	FindBySlug(ctx context.Context, slug string) (*entity.Roadmap, error)
	SlugExiste(ctx context.Context, slug string) (bool, error)
	ListarPublicados(ctx context.Context) ([]RoadmapResumo, error)
	ListarPorAutor(ctx context.Context, autorID string) ([]RoadmapResumo, error)
	// ListarNosPorArtigo devolve os nós com conclusão "artigo" ligados ao artigo, só de roadmaps publicados.
	ListarNosPorArtigo(ctx context.Context, artigoID string) ([]NoVinculado, error)
}

// ProgressoRoadmapItem é o progresso do aluno em um roadmap publicado.
type ProgressoRoadmapItem struct {
	Slug       string
	Concluidos int
	Total      int
}

type RoadmapProgressoRepository interface {
	ListarNosConcluidos(ctx context.Context, usuarioID, roadmapID string) ([]string, error)
	// MarcarNoConcluido é idempotente: devolve novo=false se o nó já estava concluído.
	MarcarNoConcluido(ctx context.Context, usuarioID, roadmapID, noID string) (novo bool, err error)
	// EstadoArtigos informa, para cada artigo com interação do aluno, se foi concluído.
	// Artigos sem registro não aparecem no mapa.
	EstadoArtigos(ctx context.Context, usuarioID string, artigoIDs []string) (map[string]bool, error)
	ContarConcluidos(ctx context.Context, usuarioID string) (nos, chefes int, err error)
	ListarProgressoPorRoadmap(ctx context.Context, usuarioID string) ([]ProgressoRoadmapItem, error)
}

type ConquistaObtida struct {
	Codigo string
	Em     time.Time
}

type GamificacaoRepository interface {
	Obter(ctx context.Context, usuarioID string) (entity.Gamificacao, error)
	// RegistrarAtividade atualiza a sequência diária sem conceder XP.
	RegistrarAtividade(ctx context.Context, usuarioID string, agora time.Time) (entity.Gamificacao, error)
	// ConcederXP é idempotente por (usuário, origem, ref). Só atualiza totais e sequência
	// quando o evento é novo (aplicado=true).
	ConcederXP(ctx context.Context, usuarioID, origem, refID string, xp int, agora time.Time) (aplicado bool, g entity.Gamificacao, err error)
	ListarConquistas(ctx context.Context, usuarioID string) ([]ConquistaObtida, error)
	// ConcederConquista é idempotente: novo=false se já estava conquistada.
	ConcederConquista(ctx context.Context, usuarioID, codigo string) (novo bool, err error)
}
