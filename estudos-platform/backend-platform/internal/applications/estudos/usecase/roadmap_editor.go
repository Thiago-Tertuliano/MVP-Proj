package usecase

import (
	"context"
	"fmt"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// RoadmapEditor reúne a autoria de roadmaps. O papel "editor" é exigido pelo middleware;
// aqui garantimos também que cada editor só altera os próprios roadmaps.
type RoadmapEditor struct {
	roadmaps repository.RoadmapRepository
}

func NewRoadmapEditor(r repository.RoadmapRepository) *RoadmapEditor {
	return &RoadmapEditor{roadmaps: r}
}

func (uc *RoadmapEditor) carregarDoAutor(ctx context.Context, autorID, slug, op string) (*entity.Roadmap, error) {
	rm, err := uc.roadmaps.FindBySlug(ctx, slug)
	if err != nil {
		return nil, traduzir(err, "falha ao buscar roadmap", op)
	}
	if rm.AutorID != autorID {
		return nil, errors.ErrForbidden("você só pode editar os seus próprios roadmaps", op, nil)
	}
	return rm, nil
}

// Meus lista os roadmaps do editor (inclui rascunhos).
func (uc *RoadmapEditor) Meus(ctx context.Context, autorID string) (*dto.ListarRoadmapsResponse, error) {
	itens, err := uc.roadmaps.ListarPorAutor(ctx, autorID)
	if err != nil {
		return nil, traduzir(err, "falha ao listar roadmaps", "RoadmapEditor.Meus")
	}
	out := &dto.ListarRoadmapsResponse{Itens: make([]dto.RoadmapResumoResponse, 0, len(itens))}
	for _, it := range itens {
		out.Itens = append(out.Itens, resumoParaDTO(it))
	}
	return out, nil
}

// Obter devolve o roadmap completo, com gabarito, para edição.
func (uc *RoadmapEditor) Obter(ctx context.Context, autorID, slug string) (*dto.RoadmapResponse, error) {
	rm, err := uc.carregarDoAutor(ctx, autorID, slug, "RoadmapEditor.Obter")
	if err != nil {
		return nil, err
	}
	return roadmapParaDTO(rm, true), nil
}

// Criar cria um rascunho vazio com slug único derivado do título.
func (uc *RoadmapEditor) Criar(ctx context.Context, autorID string, req dto.CriarRoadmapRequest) (*dto.RoadmapResponse, error) {
	const op = "RoadmapEditor.Criar"
	base, err := valueobject.NewSlug(req.Titulo)
	if err != nil {
		return nil, errors.ErrInvalidArgument("título inválido para gerar o endereço do roadmap", op, err)
	}
	slug, err := uc.slugLivre(ctx, base.Value(), op)
	if err != nil {
		return nil, err
	}
	rm, err := entity.NovoRoadmap(slug, req.Titulo, req.Descricao, req.Icone, autorID)
	if err != nil {
		return nil, err
	}
	if err := uc.roadmaps.Save(ctx, rm); err != nil {
		return nil, traduzir(err, "falha ao salvar roadmap", op)
	}
	return roadmapParaDTO(rm, true), nil
}

func (uc *RoadmapEditor) slugLivre(ctx context.Context, base, op string) (string, error) {
	candidato := base
	if base == "meus" {
		// "meus" é rota reservada (/roadmaps/meus).
		candidato = "meus-roadmap"
		base = candidato
	}
	for i := 2; i < 50; i++ {
		existe, err := uc.roadmaps.SlugExiste(ctx, candidato)
		if err != nil {
			return "", traduzir(err, "falha ao verificar endereço", op)
		}
		if !existe {
			return candidato, nil
		}
		candidato = fmt.Sprintf("%s-%d", base, i)
	}
	return "", errors.ErrAlreadyExists("não foi possível gerar um endereço único", op, nil)
}

// Salvar substitui metadados e grafo. Rascunhos aceitam conteúdo incompleto; roadmaps publicados
// exigem o grafo íntegro para nunca expor um mapa quebrado aos alunos.
func (uc *RoadmapEditor) Salvar(ctx context.Context, autorID, slug string, req dto.SalvarRoadmapRequest) (*dto.RoadmapResponse, error) {
	const op = "RoadmapEditor.Salvar"
	rm, err := uc.carregarDoAutor(ctx, autorID, slug, op)
	if err != nil {
		return nil, err
	}
	nos, arestas := grafoDeDTO(req.Nos, req.Arestas)
	if rm.Publicado {
		err = entity.ValidarGrafo(nos, arestas)
	} else {
		err = entity.ValidarEstrutura(nos, arestas)
	}
	if err != nil {
		return nil, err
	}
	if err := rm.Substituir(req.Titulo, req.Descricao, req.Icone, nos, arestas); err != nil {
		return nil, err
	}
	if err := uc.roadmaps.Save(ctx, rm); err != nil {
		return nil, traduzir(err, "falha ao salvar roadmap", op)
	}
	// Recarrega para devolver os campos enriquecidos (slug/título do artigo vinculado).
	salvo, err := uc.roadmaps.FindBySlug(ctx, slug)
	if err != nil {
		return nil, traduzir(err, "falha ao buscar roadmap", op)
	}
	return roadmapParaDTO(salvo, true), nil
}

// Publicar valida o grafo e torna o roadmap visível; Despublicar volta a rascunho.
func (uc *RoadmapEditor) Publicar(ctx context.Context, autorID, slug string, publicar bool) (*dto.RoadmapResponse, error) {
	const op = "RoadmapEditor.Publicar"
	rm, err := uc.carregarDoAutor(ctx, autorID, slug, op)
	if err != nil {
		return nil, err
	}
	if publicar {
		if err := rm.Publicar(); err != nil {
			return nil, err
		}
	} else {
		rm.Despublicar()
	}
	if err := uc.roadmaps.Save(ctx, rm); err != nil {
		return nil, traduzir(err, "falha ao salvar roadmap", op)
	}
	return roadmapParaDTO(rm, true), nil
}
