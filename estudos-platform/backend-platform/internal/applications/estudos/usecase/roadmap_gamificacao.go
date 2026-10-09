package usecase

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/dto"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

const (
	origemXPNo      = "roadmap_no"
	origemXPRoadmap = "roadmap_completo"
)

// traduzir mantém DomainError e converte o restante em erro interno.
func traduzir(err error, msg, op string) error {
	if err == nil {
		return nil
	}
	var de *errors.DomainError
	if stderrors.As(err, &de) {
		return err
	}
	return errors.ErrInternal(msg, op, err)
}

// GamificacaoRoadmap concentra as regras de XP, nível, sequência e conquistas.
// Toda concessão é idempotente e decidida no servidor.
type GamificacaoRoadmap struct {
	roadmaps repository.RoadmapRepository
	prog     repository.RoadmapProgressoRepository
	gam      repository.GamificacaoRepository
	agora    func() time.Time
}

func NewGamificacaoRoadmap(
	roadmaps repository.RoadmapRepository,
	prog repository.RoadmapProgressoRepository,
	gam repository.GamificacaoRepository,
) *GamificacaoRoadmap {
	return &GamificacaoRoadmap{roadmaps: roadmaps, prog: prog, gam: gam, agora: func() time.Time { return time.Now().UTC() }}
}

// ComRelogio permite fixar o tempo nos testes.
func (s *GamificacaoRoadmap) ComRelogio(f func() time.Time) *GamificacaoRoadmap {
	s.agora = f
	return s
}

type alvoNo struct {
	rm   *entity.Roadmap
	noID string
}

// ConcluirNo marca o nó como dominado, concede XP e avalia conquistas.
func (s *GamificacaoRoadmap) ConcluirNo(ctx context.Context, usuarioID string, rm *entity.Roadmap, noID string) (*dto.ResultadoGamificacao, error) {
	return s.concluirNos(ctx, usuarioID, []alvoNo{{rm: rm, noID: noID}})
}

func (s *GamificacaoRoadmap) concluirNos(ctx context.Context, usuarioID string, alvos []alvoNo) (*dto.ResultadoGamificacao, error) {
	const op = "GamificacaoRoadmap.concluirNos"
	agora := s.agora()

	antes, err := s.gam.Obter(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao obter gamificação", op)
	}
	acc := &dto.ResultadoGamificacao{
		NivelAnterior:     entity.NivelPorXP(antes.XPTotal),
		NosConcluidos:     []string{},
		NosDesbloqueados:  []string{},
		RoadmapsCompletos: []string{},
		ConquistasNovas:   []dto.ConquistaDTO{},
	}

	cache := map[string]map[string]bool{}
	completoAgora := false

	for _, alvo := range alvos {
		rm := alvo.rm
		no := rm.NoPorID(alvo.noID)
		if no == nil {
			continue
		}
		concluidos, ok := cache[rm.ID]
		if !ok {
			ids, err := s.prog.ListarNosConcluidos(ctx, usuarioID, rm.ID)
			if err != nil {
				return nil, traduzir(err, "falha ao obter progresso", op)
			}
			concluidos = make(map[string]bool, len(ids))
			for _, id := range ids {
				concluidos[id] = true
			}
			cache[rm.ID] = concluidos
		}

		estadosAntes := entity.CalcularEstados(rm.Nos, rm.Arestas, concluidos, nil)
		novo, err := s.prog.MarcarNoConcluido(ctx, usuarioID, rm.ID, no.ID)
		if err != nil {
			return nil, traduzir(err, "falha ao salvar progresso", op)
		}
		concluidos[no.ID] = true
		if !novo {
			continue
		}
		acc.NosConcluidos = append(acc.NosConcluidos, no.ID)

		aplicado, _, err := s.gam.ConcederXP(ctx, usuarioID, origemXPNo, no.ID, no.XP, agora)
		if err != nil {
			return nil, traduzir(err, "falha ao conceder XP", op)
		}
		if aplicado {
			acc.XPGanho += no.XP
		}

		estadosDepois := entity.CalcularEstados(rm.Nos, rm.Arestas, concluidos, nil)
		for _, n := range rm.Nos {
			if estadosAntes[n.ID] == entity.EstadoBloqueado && estadosDepois[n.ID] != entity.EstadoBloqueado {
				acc.NosDesbloqueados = append(acc.NosDesbloqueados, n.ID)
			}
		}

		if todosConcluidos(rm, concluidos) {
			completoAgora = true
			acc.RoadmapsCompletos = append(acc.RoadmapsCompletos, rm.Slug)
			bonus, _, err := s.gam.ConcederXP(ctx, usuarioID, origemXPRoadmap, rm.ID, entity.XPBonusRoadmap, agora)
			if err != nil {
				return nil, traduzir(err, "falha ao conceder XP", op)
			}
			if bonus {
				acc.XPGanho += entity.XPBonusRoadmap
			}
		}
	}

	return s.finalizar(ctx, usuarioID, acc, completoAgora, agora)
}

// finalizar lê os totais finais, avalia conquistas e preenche nível/streak.
func (s *GamificacaoRoadmap) finalizar(ctx context.Context, usuarioID string, acc *dto.ResultadoGamificacao, completoAgora bool, agora time.Time) (*dto.ResultadoGamificacao, error) {
	const op = "GamificacaoRoadmap.finalizar"
	g, err := s.gam.Obter(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao obter gamificação", op)
	}
	nos, chefes, err := s.prog.ContarConcluidos(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao contar progresso", op)
	}
	obtidas, err := s.gam.ListarConquistas(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao listar conquistas", op)
	}
	ja := make(map[string]bool, len(obtidas))
	for _, c := range obtidas {
		ja[c.Codigo] = true
	}

	elegiveis := entity.ConquistasElegiveis(entity.SnapshotConquistas{
		XPTotal: g.XPTotal, StreakMax: g.StreakMax, NosConcluidos: nos, ChefesConcluidos: chefes, RoadmapCompleto: completoAgora,
	})
	for _, codigo := range elegiveis {
		if ja[codigo] {
			continue
		}
		novo, err := s.gam.ConcederConquista(ctx, usuarioID, codigo)
		if err != nil {
			return nil, traduzir(err, "falha ao conceder conquista", op)
		}
		if !novo {
			continue
		}
		if c, ok := entity.ConquistaPorCodigo(codigo); ok {
			em := agora.Unix()
			acc.ConquistasNovas = append(acc.ConquistasNovas, dto.ConquistaDTO{
				Codigo: c.Codigo, Nome: c.Nome, Descricao: c.Descricao, Icone: c.Icone, Conquistada: true, Em: &em,
			})
		}
	}

	acc.XPTotal = g.XPTotal
	acc.Nivel = entity.NivelPorXP(g.XPTotal)
	acc.SubiuDeNivel = acc.Nivel > acc.NivelAnterior
	acc.Streak = entity.StreakVigente(g, agora)
	return acc, nil
}

func todosConcluidos(rm *entity.Roadmap, concluidos map[string]bool) bool {
	if len(rm.Nos) == 0 {
		return false
	}
	for _, n := range rm.Nos {
		if !concluidos[n.ID] {
			return false
		}
	}
	return true
}

// AoConcluirArtigo registra atividade (sequência diária) e conclui os nós ligados ao artigo.
func (s *GamificacaoRoadmap) AoConcluirArtigo(ctx context.Context, usuarioID, artigoID string) (*dto.ResultadoGamificacao, error) {
	const op = "GamificacaoRoadmap.AoConcluirArtigo"
	agora := s.agora()
	if _, err := s.gam.RegistrarAtividade(ctx, usuarioID, agora); err != nil {
		return nil, traduzir(err, "falha ao registrar atividade", op)
	}
	refs, err := s.roadmaps.ListarNosPorArtigo(ctx, artigoID)
	if err != nil {
		return nil, traduzir(err, "falha ao buscar nós do artigo", op)
	}

	roadmaps := map[string]*entity.Roadmap{}
	var alvos []alvoNo
	for _, ref := range refs {
		rm, ok := roadmaps[ref.RoadmapSlug]
		if !ok {
			rm, err = s.roadmaps.FindBySlug(ctx, ref.RoadmapSlug)
			if err != nil {
				return nil, traduzir(err, "falha ao buscar roadmap", op)
			}
			roadmaps[ref.RoadmapSlug] = rm
		}
		// Só conclui nós já liberados: pré-requisitos pendentes impedem o avanço.
		alvos = append(alvos, alvoNo{rm: rm, noID: ref.NoID})
	}
	return s.concluirNosLiberados(ctx, usuarioID, alvos)
}

// concluirNosLiberados conclui apenas nós cujos pré-requisitos já foram dominados; repete
// até estabilizar, para que uma cadeia de artigos já lidos avance em sequência.
func (s *GamificacaoRoadmap) concluirNosLiberados(ctx context.Context, usuarioID string, alvos []alvoNo) (*dto.ResultadoGamificacao, error) {
	const op = "GamificacaoRoadmap.concluirNosLiberados"
	agora := s.agora()
	antes, err := s.gam.Obter(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao obter gamificação", op)
	}
	total := &dto.ResultadoGamificacao{
		NivelAnterior:     entity.NivelPorXP(antes.XPTotal),
		NosConcluidos:     []string{},
		NosDesbloqueados:  []string{},
		RoadmapsCompletos: []string{},
		ConquistasNovas:   []dto.ConquistaDTO{},
	}

	pendentes := alvos
	for rodada := 0; rodada < entity.MaxNosPorRoadmap && len(pendentes) > 0; rodada++ {
		var liberados, restantes []alvoNo
		for _, a := range pendentes {
			ids, err := s.prog.ListarNosConcluidos(ctx, usuarioID, a.rm.ID)
			if err != nil {
				return nil, traduzir(err, "falha ao obter progresso", op)
			}
			c := make(map[string]bool, len(ids))
			for _, id := range ids {
				c[id] = true
			}
			if c[a.noID] {
				continue
			}
			if entity.CalcularEstados(a.rm.Nos, a.rm.Arestas, c, nil)[a.noID] == entity.EstadoBloqueado {
				restantes = append(restantes, a)
				continue
			}
			liberados = append(liberados, a)
		}
		if len(liberados) == 0 {
			break
		}
		parcial, err := s.concluirNos(ctx, usuarioID, liberados)
		if err != nil {
			return nil, err
		}
		total.XPGanho += parcial.XPGanho
		total.NosConcluidos = append(total.NosConcluidos, parcial.NosConcluidos...)
		total.NosDesbloqueados = append(total.NosDesbloqueados, parcial.NosDesbloqueados...)
		total.RoadmapsCompletos = append(total.RoadmapsCompletos, parcial.RoadmapsCompletos...)
		total.ConquistasNovas = append(total.ConquistasNovas, parcial.ConquistasNovas...)
		pendentes = restantes
	}

	completo := len(total.RoadmapsCompletos) > 0
	return s.finalizar(ctx, usuarioID, total, completo, agora)
}

// SincronizarArtigos conclui nós "por artigo" cujo artigo o aluno já leu antes do nó existir.
func (s *GamificacaoRoadmap) SincronizarArtigos(ctx context.Context, usuarioID string, rm *entity.Roadmap) error {
	const op = "GamificacaoRoadmap.SincronizarArtigos"
	var artigos []string
	var alvos []alvoNo
	for _, n := range rm.Nos {
		if n.Conclusao == entity.ConclusaoArtigo && n.ArtigoID != "" {
			artigos = append(artigos, n.ArtigoID)
			alvos = append(alvos, alvoNo{rm: rm, noID: n.ID})
		}
	}
	if len(artigos) == 0 {
		return nil
	}
	estado, err := s.prog.EstadoArtigos(ctx, usuarioID, artigos)
	if err != nil {
		return traduzir(err, "falha ao consultar leituras", op)
	}
	var lidos []alvoNo
	for i, a := range alvos {
		if estado[artigos[i]] {
			lidos = append(lidos, a)
		}
	}
	if len(lidos) == 0 {
		return nil
	}
	_, err = s.concluirNosLiberados(ctx, usuarioID, lidos)
	return err
}

// Resumo devolve o painel de gamificação do aluno (XP, nível, sequência e catálogo de conquistas).
func (s *GamificacaoRoadmap) Resumo(ctx context.Context, usuarioID string) (*dto.GamificacaoResponse, error) {
	const op = "GamificacaoRoadmap.Resumo"
	agora := s.agora()
	g, err := s.gam.Obter(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao obter gamificação", op)
	}
	obtidas, err := s.gam.ListarConquistas(ctx, usuarioID)
	if err != nil {
		return nil, traduzir(err, "falha ao listar conquistas", op)
	}
	quando := make(map[string]int64, len(obtidas))
	for _, c := range obtidas {
		quando[c.Codigo] = c.Em.Unix()
	}

	conquistas := make([]dto.ConquistaDTO, 0, len(entity.CatalogoConquistas))
	for _, c := range entity.CatalogoConquistas {
		d := dto.ConquistaDTO{Codigo: c.Codigo, Nome: c.Nome, Descricao: c.Descricao, Icone: c.Icone}
		if em, ok := quando[c.Codigo]; ok {
			d.Conquistada = true
			d.Em = &em
		}
		conquistas = append(conquistas, d)
	}

	noNivel, proximo := entity.ProgressoNivel(g.XPTotal)
	ativoHoje := g.UltimoDia != nil && entity.DiaCalendario(agora).Equal(time.Date(g.UltimoDia.Year(), g.UltimoDia.Month(), g.UltimoDia.Day(), 0, 0, 0, 0, time.UTC))
	return &dto.GamificacaoResponse{
		XPTotal: g.XPTotal, Nivel: entity.NivelPorXP(g.XPTotal), XPNoNivel: noNivel, XPParaProximo: proximo,
		StreakAtual: entity.StreakVigente(g, agora), StreakMax: g.StreakMax, AtivoHoje: ativoHoje,
		Conquistas: conquistas,
	}, nil
}
