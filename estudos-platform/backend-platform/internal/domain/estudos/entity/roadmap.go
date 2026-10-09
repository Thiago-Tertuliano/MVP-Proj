package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

type TipoNo string

const (
	NoTopico TipoNo = "topico"
	NoChefe  TipoNo = "chefe"
	NoMarco  TipoNo = "marco"
)

type ConclusaoNo string

const (
	ConclusaoManual ConclusaoNo = "manual"
	ConclusaoArtigo ConclusaoNo = "artigo"
	ConclusaoQuiz   ConclusaoNo = "quiz"
)

type TipoAresta string

const (
	ArestaRequer   TipoAresta = "requer"
	ArestaOpcional TipoAresta = "opcional"
)

const (
	MaxNosPorRoadmap     = 200
	MaxXPPorNo           = 1000
	XPPadraoNo           = 10
	XPPadraoChefe        = 100
	XPBonusRoadmap       = 50
	NotaMinimaQuizPct    = 70
	MaxQuestoesPorQuiz   = 20
	MaxOpcoesPorQuestao  = 6
	MinOpcoesPorQuestao  = 2
	MaxTituloRoadmap     = 200
	MaxTituloNo          = 120
	MaxDescricaoRoadmap  = 2000
	MaxDescricaoNo       = 2000
	MaxEnunciadoQuestao  = 500
	MaxTextoOpcaoQuestao = 300
)

// OpcaoQuiz é uma alternativa de uma questão.
type OpcaoQuiz struct {
	ID    string `json:"id"`
	Texto string `json:"texto"`
}

// QuestaoQuiz é uma questão de múltipla escolha. Correta nunca sai para o aluno.
type QuestaoQuiz struct {
	ID         string      `json:"id"`
	Enunciado  string      `json:"enunciado"`
	Opcoes     []OpcaoQuiz `json:"opcoes"`
	Correta    string      `json:"correta"`
	Explicacao string      `json:"explicacao,omitempty"`
}

// Quiz é o desafio de um nó chefe.
type Quiz struct {
	Questoes []QuestaoQuiz `json:"questoes"`
}

// RoadmapNo é um nó do grafo.
type RoadmapNo struct {
	ID        string
	Tipo      TipoNo
	Titulo    string
	Descricao string
	ArtigoID  string // vazio = sem artigo vinculado
	// ArtigoSlug/ArtigoTitulo são enriquecimentos somente-leitura preenchidos na leitura.
	ArtigoSlug   string
	ArtigoTitulo string
	PosX         float64
	PosY         float64
	XP           int
	Conclusao    ConclusaoNo
	Quiz         *Quiz
}

// RoadmapAresta liga dois nós (origem → destino).
type RoadmapAresta struct {
	ID        string
	OrigemID  string
	DestinoID string
	Tipo      TipoAresta
}

// Roadmap é o Aggregate Root do mapa gamificado.
type Roadmap struct {
	ID        string
	Slug      string
	Titulo    string
	Descricao string
	Icone     string
	Publicado bool
	AutorID   string
	Nos       []RoadmapNo
	Arestas   []RoadmapAresta
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NovoRoadmap cria um rascunho vazio.
func NovoRoadmap(slug, titulo, descricao, icone, autorID string) (*Roadmap, error) {
	titulo = strings.TrimSpace(titulo)
	if len(titulo) < 3 || len(titulo) > MaxTituloRoadmap {
		return nil, errors.ErrInvalidArgument("título deve ter entre 3 e 200 caracteres", "entity.NovoRoadmap", nil)
	}
	if _, err := uuid.Parse(autorID); err != nil {
		return nil, errors.ErrInvalidArgument("autor inválido", "entity.NovoRoadmap", err)
	}
	if len(descricao) > MaxDescricaoRoadmap {
		return nil, errors.ErrInvalidArgument("descrição muito longa", "entity.NovoRoadmap", nil)
	}
	now := time.Now().UTC()
	return &Roadmap{
		ID:        uuid.NewString(),
		Slug:      slug,
		Titulo:    titulo,
		Descricao: strings.TrimSpace(descricao),
		Icone:     strings.TrimSpace(icone),
		AutorID:   autorID,
		Nos:       []RoadmapNo{},
		Arestas:   []RoadmapAresta{},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NoPorID devolve o nó ou nil.
func (r *Roadmap) NoPorID(id string) *RoadmapNo {
	for i := range r.Nos {
		if r.Nos[i].ID == id {
			return &r.Nos[i]
		}
	}
	return nil
}

// XPTotal soma o XP de todos os nós.
func (r *Roadmap) XPTotal() int {
	total := 0
	for _, n := range r.Nos {
		total += n.XP
	}
	return total
}

// Substituir troca metadados e grafo (edição completa pelo editor).
func (r *Roadmap) Substituir(titulo, descricao, icone string, nos []RoadmapNo, arestas []RoadmapAresta) error {
	titulo = strings.TrimSpace(titulo)
	if len(titulo) < 3 || len(titulo) > MaxTituloRoadmap {
		return errors.ErrInvalidArgument("título deve ter entre 3 e 200 caracteres", "Roadmap.Substituir", nil)
	}
	if len(descricao) > MaxDescricaoRoadmap {
		return errors.ErrInvalidArgument("descrição muito longa", "Roadmap.Substituir", nil)
	}
	r.Titulo = titulo
	r.Descricao = strings.TrimSpace(descricao)
	r.Icone = strings.TrimSpace(icone)
	r.Nos = nos
	r.Arestas = arestas
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// Publicar exige um grafo válido e ao menos um nó.
func (r *Roadmap) Publicar() error {
	if len(r.Nos) == 0 {
		return errors.ErrInvalidState("roadmap sem nós não pode ser publicado", "Roadmap.Publicar", nil)
	}
	if err := ValidarGrafo(r.Nos, r.Arestas); err != nil {
		return err
	}
	r.Publicado = true
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *Roadmap) Despublicar() {
	r.Publicado = false
	r.UpdatedAt = time.Now().UTC()
}
