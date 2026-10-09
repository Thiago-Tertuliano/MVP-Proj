package usecase

import (
	"context"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// RoadmapAluno agrupa as ações do aluno sobre um roadmap publicado.
type RoadmapAluno struct {
	roadmaps repository.RoadmapRepository
	prog     repository.RoadmapProgressoRepository
	gam      *GamificacaoRoadmap
}

func NewRoadmapAluno(r repository.RoadmapRepository, p repository.RoadmapProgressoRepository, g *GamificacaoRoadmap) *RoadmapAluno {
	return &RoadmapAluno{roadmaps: r, prog: p, gam: g}
}

func (uc *RoadmapAluno) carregarPublicado(ctx context.Context, slug, op string) (*entity.Roadmap, error) {
	rm, err := uc.roadmaps.FindBySlug(ctx, slug)
	if err != nil {
		return nil, traduzir(err, "falha ao buscar roadmap", op)
	}
	if !rm.Publicado {
		return nil, errors.ErrNotFound("roadmap não encontrado", op, nil)
	}
	return rm, nil
}

// situacao calcula os estados do aluno (incluindo "em curso" via artigos abertos).
func (uc *RoadmapAluno) situacao(ctx context.Context, usuarioID string, rm *entity.Roadmap, op string) (map[string]entity.EstadoNo, map[string]bool, error) {
	ids, err := uc.prog.ListarNosConcluidos(ctx, usuarioID, rm.ID)
	if err != nil {
		return nil, nil, traduzir(err, "falha ao obter progresso", op)
	}
	concluidos := make(map[string]bool, len(ids))
	for _, id := range ids {
		concluidos[id] = true
	}

	var artigos []string
	for _, n := range rm.Nos {
		if n.ArtigoID != "" && !concluidos[n.ID] {
			artigos = append(artigos, n.ArtigoID)
		}
	}
	tocados, err := uc.prog.EstadoArtigos(ctx, usuarioID, artigos)
	if err != nil {
		return nil, nil, traduzir(err, "falha ao consultar leituras", op)
	}
	emCurso := map[string]bool{}
	for _, n := range rm.Nos {
		if concl, ok := tocados[n.ArtigoID]; ok && n.ArtigoID != "" && !concl {
			emCurso[n.ID] = true
		}
	}
	return entity.CalcularEstados(rm.Nos, rm.Arestas, concluidos, emCurso), concluidos, nil
}

func (uc *RoadmapAluno) montarProgresso(rm *entity.Roadmap, estados map[string]entity.EstadoNo, concluidos map[string]bool) *dto.RoadmapProgressoResponse {
	xpGanho := 0
	for _, n := range rm.Nos {
		if concluidos[n.ID] {
			xpGanho += n.XP
		}
	}
	total := len(rm.Nos)
	pct := 0.0
	if total > 0 {
		pct = float64(len(concluidos)) * 100 / float64(total)
	}
	return &dto.RoadmapProgressoResponse{
		Slug: rm.Slug, Estados: estadosParaString(estados), Concluidos: len(concluidos), Total: total,
		Percentual: pct, XPGanho: xpGanho, XPTotal: rm.XPTotal(), Completo: total > 0 && len(concluidos) == total,
	}
}

// ObterProgresso devolve o estado de cada nó. Antes, sincroniza leituras antigas de artigos.
func (uc *RoadmapAluno) ObterProgresso(ctx context.Context, usuarioID, slug string) (*dto.RoadmapProgressoResponse, error) {
	const op = "RoadmapAluno.ObterProgresso"
	rm, err := uc.carregarPublicado(ctx, slug, op)
	if err != nil {
		return nil, err
	}
	if err := uc.gam.SincronizarArtigos(ctx, usuarioID, rm); err != nil {
		return nil, err
	}
	estados, concluidos, err := uc.situacao(ctx, usuarioID, rm, op)
	if err != nil {
		return nil, err
	}
	return uc.montarProgresso(rm, estados, concluidos), nil
}

func (uc *RoadmapAluno) respostaConclusao(ctx context.Context, usuarioID string, rm *entity.Roadmap, noID string, ja bool, res *dto.ResultadoGamificacao, op string) (*dto.ConclusaoNoResponse, error) {
	estados, concluidos, err := uc.situacao(ctx, usuarioID, rm, op)
	if err != nil {
		return nil, err
	}
	return &dto.ConclusaoNoResponse{
		NoID: noID, JaConcluido: ja, Resultado: *res, Estados: estadosParaString(estados),
		Concluidos: len(concluidos), Total: len(rm.Nos),
	}, nil
}

// prepararConclusao valida existência, forma de conclusão e bloqueio do nó.
func (uc *RoadmapAluno) prepararConclusao(ctx context.Context, usuarioID, slug, noID string, esperada entity.ConclusaoNo, op string) (*entity.Roadmap, *entity.RoadmapNo, entity.EstadoNo, error) {
	rm, err := uc.carregarPublicado(ctx, slug, op)
	if err != nil {
		return nil, nil, "", err
	}
	no := rm.NoPorID(noID)
	if no == nil {
		return nil, nil, "", errors.ErrNotFound("nó não encontrado", op, nil)
	}
	if no.Conclusao != esperada {
		return nil, nil, "", errors.ErrInvalidArgument(mensagemConclusao(no.Conclusao), op, nil)
	}
	estados, _, err := uc.situacao(ctx, usuarioID, rm, op)
	if err != nil {
		return nil, nil, "", err
	}
	estado := estados[no.ID]
	if estado == entity.EstadoBloqueado {
		return nil, nil, "", errors.ErrInvalidState("nó bloqueado: conclua os pré-requisitos primeiro", op, nil)
	}
	return rm, no, estado, nil
}

func mensagemConclusao(c entity.ConclusaoNo) string {
	switch c {
	case entity.ConclusaoArtigo:
		return "este nó é concluído lendo o artigo vinculado"
	case entity.ConclusaoQuiz:
		return "este nó é concluído acertando o quiz do chefe"
	default:
		return "este nó é concluído manualmente"
	}
}

func resultadoVazio(g *GamificacaoRoadmap, ctx context.Context, usuarioID string) (*dto.ResultadoGamificacao, error) {
	resumo, err := g.Resumo(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	return &dto.ResultadoGamificacao{
		XPTotal: resumo.XPTotal, Nivel: resumo.Nivel, NivelAnterior: resumo.Nivel, Streak: resumo.StreakAtual,
		NosConcluidos: []string{}, NosDesbloqueados: []string{}, RoadmapsCompletos: []string{}, ConquistasNovas: []dto.ConquistaDTO{},
	}, nil
}

// Concluir marca um nó de conclusão manual ("dominei") como dominado.
func (uc *RoadmapAluno) Concluir(ctx context.Context, usuarioID, slug, noID string) (*dto.ConclusaoNoResponse, error) {
	const op = "RoadmapAluno.Concluir"
	rm, no, estado, err := uc.prepararConclusao(ctx, usuarioID, slug, noID, entity.ConclusaoManual, op)
	if err != nil {
		return nil, err
	}
	if estado == entity.EstadoDominado {
		res, err := resultadoVazio(uc.gam, ctx, usuarioID)
		if err != nil {
			return nil, err
		}
		return uc.respostaConclusao(ctx, usuarioID, rm, no.ID, true, res, op)
	}
	res, err := uc.gam.ConcluirNo(ctx, usuarioID, rm, no.ID)
	if err != nil {
		return nil, err
	}
	return uc.respostaConclusao(ctx, usuarioID, rm, no.ID, false, res, op)
}

// ResponderQuiz corrige o quiz do chefe; só conclui o nó com nota suficiente.
func (uc *RoadmapAluno) ResponderQuiz(ctx context.Context, usuarioID, slug, noID string, respostas map[string]string) (*dto.QuizResultadoResponse, error) {
	const op = "RoadmapAluno.ResponderQuiz"
	rm, no, estado, err := uc.prepararConclusao(ctx, usuarioID, slug, noID, entity.ConclusaoQuiz, op)
	if err != nil {
		return nil, err
	}

	av := entity.AvaliarQuiz(no.Quiz, respostas)
	resp := &dto.QuizResultadoResponse{
		Aprovado: av.Aprovado, Acertos: av.Acertos, Total: av.Total, NotaMinimaPct: entity.NotaMinimaQuizPct,
	}
	if !av.Aprovado {
		return resp, nil
	}
	for _, c := range av.Correcao {
		resp.Correcao = append(resp.Correcao, dto.CorrecaoQuestaoDTO{QuestaoID: c.QuestaoID, Correta: c.Correta, Explicacao: c.Explicacao})
	}

	if estado == entity.EstadoDominado {
		res, err := resultadoVazio(uc.gam, ctx, usuarioID)
		if err != nil {
			return nil, err
		}
		resp.Conclusao, err = uc.respostaConclusao(ctx, usuarioID, rm, no.ID, true, res, op)
		return resp, err
	}
	res, err := uc.gam.ConcluirNo(ctx, usuarioID, rm, no.ID)
	if err != nil {
		return nil, err
	}
	resp.Conclusao, err = uc.respostaConclusao(ctx, usuarioID, rm, no.ID, false, res, op)
	return resp, err
}

// ListarProgresso resume o avanço do aluno em cada roadmap publicado (catálogo).
func (uc *RoadmapAluno) ListarProgresso(ctx context.Context, usuarioID string) (*dto.ListarProgressoRoadmapsResponse, error) {
	itens, err := uc.prog.ListarProgressoPorRoadmap(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao obter progresso", "RoadmapAluno.ListarProgresso")
	}
	out := &dto.ListarProgressoRoadmapsResponse{Itens: make([]dto.RoadmapProgressoItem, 0, len(itens))}
	for _, it := range itens {
		pct := 0.0
		if it.Total > 0 {
			pct = float64(it.Concluidos) * 100 / float64(it.Total)
		}
		out.Itens = append(out.Itens, dto.RoadmapProgressoItem{Slug: it.Slug, Concluidos: it.Concluidos, Total: it.Total, Percentual: pct})
	}
	return out, nil
}

// Gamificacao devolve o painel do aluno (XP, nível, sequência, conquistas).
func (uc *RoadmapAluno) Gamificacao(ctx context.Context, usuarioID string) (*dto.GamificacaoResponse, error) {
	return uc.gam.Resumo(ctx, usuarioID)
}
