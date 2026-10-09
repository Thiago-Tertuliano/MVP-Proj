package usecase

import (
	"strings"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
)

// roadmapParaDTO converte para resposta. comGabarito=false remove respostas corretas e explicações
// (visão do aluno); true é usado apenas na rota do editor.
func roadmapParaDTO(rm *entity.Roadmap, comGabarito bool) *dto.RoadmapResponse {
	nos := make([]dto.RoadmapNoDTO, 0, len(rm.Nos))
	for _, n := range rm.Nos {
		d := dto.RoadmapNoDTO{
			ID: n.ID, Tipo: string(n.Tipo), Titulo: n.Titulo, Descricao: n.Descricao,
			ArtigoID: n.ArtigoID, ArtigoSlug: n.ArtigoSlug, ArtigoTitulo: n.ArtigoTitulo,
			PosX: n.PosX, PosY: n.PosY, XP: n.XP, Conclusao: string(n.Conclusao),
		}
		if n.Quiz != nil {
			d.Quiz = quizParaDTO(n.Quiz, comGabarito)
		}
		nos = append(nos, d)
	}
	arestas := make([]dto.RoadmapArestaDTO, 0, len(rm.Arestas))
	for _, a := range rm.Arestas {
		arestas = append(arestas, dto.RoadmapArestaDTO{Origem: a.OrigemID, Destino: a.DestinoID, Tipo: string(a.Tipo)})
	}
	return &dto.RoadmapResponse{
		ID: rm.ID, Slug: rm.Slug, Titulo: rm.Titulo, Descricao: rm.Descricao, Icone: rm.Icone,
		Publicado: rm.Publicado, Nos: nos, Arestas: arestas, XPTotal: rm.XPTotal(),
		AtualizadoEm: rm.UpdatedAt.Unix(),
	}
}

func quizParaDTO(q *entity.Quiz, comGabarito bool) *dto.QuizDTO {
	out := &dto.QuizDTO{Questoes: make([]dto.QuestaoQuizDTO, 0, len(q.Questoes))}
	for _, qs := range q.Questoes {
		d := dto.QuestaoQuizDTO{ID: qs.ID, Enunciado: qs.Enunciado, Opcoes: make([]dto.OpcaoQuizDTO, 0, len(qs.Opcoes))}
		for _, o := range qs.Opcoes {
			d.Opcoes = append(d.Opcoes, dto.OpcaoQuizDTO{ID: o.ID, Texto: o.Texto})
		}
		if comGabarito {
			d.Correta, d.Explicacao = qs.Correta, qs.Explicacao
		}
		out.Questoes = append(out.Questoes, d)
	}
	return out
}

func resumoParaDTO(r repository.RoadmapResumo) dto.RoadmapResumoResponse {
	return dto.RoadmapResumoResponse{
		ID: r.ID, Slug: r.Slug, Titulo: r.Titulo, Descricao: r.Descricao, Icone: r.Icone,
		Publicado: r.Publicado, TotalNos: r.TotalNos, TotalChefes: r.TotalChefes, XPTotal: r.XPTotal,
		AtualizadoEm: r.AtualizadoEm.Unix(),
	}
}

// grafoDeDTO converte o payload do editor em entidades, aplicando padrões (id, xp, conclusão).
func grafoDeDTO(nosDTO []dto.RoadmapNoDTO, arestasDTO []dto.RoadmapArestaDTO) ([]entity.RoadmapNo, []entity.RoadmapAresta) {
	nos := make([]entity.RoadmapNo, 0, len(nosDTO))
	for _, d := range nosDTO {
		id := strings.TrimSpace(d.ID)
		if id == "" {
			id = uuid.NewString()
		}
		tipo := entity.TipoNo(d.Tipo)
		if tipo == "" {
			tipo = entity.NoTopico
		}
		conclusao := entity.ConclusaoNo(d.Conclusao)
		if conclusao == "" {
			switch {
			case tipo == entity.NoChefe:
				conclusao = entity.ConclusaoQuiz
			case d.ArtigoID != "":
				conclusao = entity.ConclusaoArtigo
			default:
				conclusao = entity.ConclusaoManual
			}
		}
		n := entity.RoadmapNo{
			ID: id, Tipo: tipo, Titulo: strings.TrimSpace(d.Titulo), Descricao: strings.TrimSpace(d.Descricao),
			ArtigoID: strings.TrimSpace(d.ArtigoID), PosX: d.PosX, PosY: d.PosY, XP: d.XP, Conclusao: conclusao,
		}
		if d.Quiz != nil && conclusao == entity.ConclusaoQuiz {
			n.Quiz = quizDeDTO(d.Quiz)
		}
		nos = append(nos, n)
	}
	arestas := make([]entity.RoadmapAresta, 0, len(arestasDTO))
	for _, a := range arestasDTO {
		tipo := entity.TipoAresta(a.Tipo)
		if tipo == "" {
			tipo = entity.ArestaRequer
		}
		arestas = append(arestas, entity.RoadmapAresta{ID: uuid.NewString(), OrigemID: a.Origem, DestinoID: a.Destino, Tipo: tipo})
	}
	return nos, arestas
}

func quizDeDTO(q *dto.QuizDTO) *entity.Quiz {
	out := &entity.Quiz{Questoes: make([]entity.QuestaoQuiz, 0, len(q.Questoes))}
	for _, qs := range q.Questoes {
		e := entity.QuestaoQuiz{
			ID: strings.TrimSpace(qs.ID), Enunciado: strings.TrimSpace(qs.Enunciado),
			Correta: strings.TrimSpace(qs.Correta), Explicacao: strings.TrimSpace(qs.Explicacao),
		}
		for _, o := range qs.Opcoes {
			e.Opcoes = append(e.Opcoes, entity.OpcaoQuiz{ID: strings.TrimSpace(o.ID), Texto: strings.TrimSpace(o.Texto)})
		}
		out.Questoes = append(out.Questoes, e)
	}
	return out
}

func estadosParaString(m map[string]entity.EstadoNo) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = string(v)
	}
	return out
}
