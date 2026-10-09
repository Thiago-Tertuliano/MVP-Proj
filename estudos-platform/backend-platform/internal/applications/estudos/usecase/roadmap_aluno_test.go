package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	domainErros "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

func kind(t *testing.T, err error) domainErros.Kind {
	t.Helper()
	var de *domainErros.DomainError
	if !errors.As(err, &de) {
		t.Fatalf("esperava DomainError, got %v", err)
	}
	return de.Kind
}

func TestRoadmapAluno_ConcluirManual_XPIdempotenteEDesbloqueio(t *testing.T) {
	c := novoCenario(true)
	ctx := context.Background()

	resp, err := c.aluno.Concluir(ctx, c.user, "go-basico", c.nos["a"])
	if err != nil {
		t.Fatal(err)
	}
	if resp.JaConcluido || resp.Resultado.XPGanho != 10 || resp.Resultado.XPTotal != 10 {
		t.Fatalf("resultado: %+v", resp.Resultado)
	}
	if len(resp.Resultado.NosDesbloqueados) != 1 || resp.Resultado.NosDesbloqueados[0] != c.nos["b"] {
		t.Fatalf("desbloqueados: %v", resp.Resultado.NosDesbloqueados)
	}
	if resp.Estados[c.nos["a"]] != "dominado" || resp.Estados[c.nos["b"]] != "disponivel" || resp.Estados[c.nos["chefe"]] != "bloqueado" {
		t.Fatalf("estados: %v", resp.Estados)
	}
	if resp.Resultado.Streak != 1 {
		t.Fatalf("streak = %d", resp.Resultado.Streak)
	}
	found := false
	for _, q := range resp.Resultado.ConquistasNovas {
		if q.Codigo == "primeiro-passo" {
			found = true
		}
	}
	if !found {
		t.Fatalf("faltou 'primeiro-passo': %+v", resp.Resultado.ConquistasNovas)
	}

	// repetir não duplica XP
	resp, err = c.aluno.Concluir(ctx, c.user, "go-basico", c.nos["a"])
	if err != nil {
		t.Fatal(err)
	}
	if !resp.JaConcluido || resp.Resultado.XPGanho != 0 || resp.Resultado.XPTotal != 10 {
		t.Fatalf("repetição: %+v", resp)
	}
}

func TestRoadmapAluno_NoBloqueadoNaoConclui(t *testing.T) {
	c := novoCenario(true)
	_, err := c.aluno.Concluir(context.Background(), c.user, "go-basico", c.nos["b"])
	if kind(t, err) != domainErros.InvalidState {
		t.Fatalf("esperava InvalidState, got %v", err)
	}
	if c.gam.get(c.user).XPTotal != 0 {
		t.Fatal("não deveria ter XP")
	}
}

func TestRoadmapAluno_FormaDeConclusaoIncorreta(t *testing.T) {
	c := novoCenario(true)
	ctx := context.Background()
	// nó por artigo e nó chefe não aceitam "dominei"
	if _, err := c.aluno.Concluir(ctx, c.user, "go-basico", c.nos["c"]); kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("artigo: %v", err)
	}
	if _, err := c.aluno.Concluir(ctx, c.user, "go-basico", c.nos["chefe"]); kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("chefe: %v", err)
	}
	// e o quiz só vale para nó de quiz
	if _, err := c.aluno.ResponderQuiz(ctx, c.user, "go-basico", c.nos["a"], nil); kind(t, err) != domainErros.InvalidArgument {
		t.Fatalf("quiz em manual: %v", err)
	}
}

func TestRoadmapAluno_RascunhoEInexistente(t *testing.T) {
	ctx := context.Background()
	c := novoCenario(false)
	if _, err := c.aluno.Concluir(ctx, c.user, "go-basico", c.nos["a"]); kind(t, err) != domainErros.NotFound {
		t.Fatalf("rascunho: %v", err)
	}
	c = novoCenario(true)
	if _, err := c.aluno.Concluir(ctx, c.user, "outro", c.nos["a"]); kind(t, err) != domainErros.NotFound {
		t.Fatalf("slug inexistente: %v", err)
	}
	if _, err := c.aluno.Concluir(ctx, c.user, "go-basico", "no-que-nao-existe"); kind(t, err) != domainErros.NotFound {
		t.Fatalf("nó inexistente: %v", err)
	}
}

func prepararChefe(t *testing.T, c *cenario) {
	t.Helper()
	ctx := context.Background()
	for _, k := range []string{"a", "b"} {
		if _, err := c.aluno.Concluir(ctx, c.user, "go-basico", c.nos[k]); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRoadmapAluno_QuizReprovadoNaoConclui(t *testing.T) {
	c := novoCenario(true)
	prepararChefe(t, c)
	xpAntes := c.gam.get(c.user).XPTotal

	resp, err := c.aluno.ResponderQuiz(context.Background(), c.user, "go-basico", c.nos["chefe"], map[string]string{"q1": "y", "q2": "x"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Aprovado || resp.Conclusao != nil || len(resp.Correcao) != 0 {
		t.Fatalf("reprovado não deve concluir nem vazar gabarito: %+v", resp)
	}
	if resp.Acertos != 0 || resp.Total != 2 || resp.NotaMinimaPct != 70 {
		t.Fatalf("contagem: %+v", resp)
	}
	if c.gam.get(c.user).XPTotal != xpAntes {
		t.Fatal("XP não deveria mudar")
	}
}

func TestRoadmapAluno_QuizAprovadoConcedeXPChefeConquistasEBonus(t *testing.T) {
	c := novoCenario(true)
	prepararChefe(t, c)

	resp, err := c.aluno.ResponderQuiz(context.Background(), c.user, "go-basico", c.nos["chefe"], map[string]string{"q1": "x", "q2": "y"})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Aprovado || resp.Conclusao == nil || len(resp.Correcao) != 2 {
		t.Fatalf("aprovado deveria concluir e trazer correção: %+v", resp)
	}
	if resp.Correcao[0].Explicacao != "porque x" {
		t.Fatalf("correção: %+v", resp.Correcao)
	}
	if resp.Conclusao.Resultado.XPGanho != 100 {
		t.Fatalf("xp do chefe: %+v", resp.Conclusao.Resultado)
	}
	temConquista := func(codigo string) bool {
		for _, q := range resp.Conclusao.Resultado.ConquistasNovas {
			if q.Codigo == codigo {
				return true
			}
		}
		return false
	}
	if !temConquista("cacador-de-chefes") {
		t.Fatalf("faltou caçador de chefes: %+v", resp.Conclusao.Resultado.ConquistasNovas)
	}
	// 10 + 10 = 20 XP (nível 1) → com os 100 do chefe, 120 XP = nível 2.
	if !resp.Conclusao.Resultado.SubiuDeNivel || resp.Conclusao.Resultado.Nivel != 2 {
		t.Fatalf("nível: %+v", resp.Conclusao.Resultado)
	}
}

func TestRoadmapAluno_Progresso_SincronizaArtigoJaLido(t *testing.T) {
	c := novoCenario(true)
	ctx := context.Background()
	c.prog.artigos[c.rm.Nos[3].ArtigoID] = true // leu o artigo antes de abrir o roadmap

	p, err := c.aluno.ObterProgresso(ctx, c.user, "go-basico")
	if err != nil {
		t.Fatal(err)
	}
	if p.Estados[c.nos["c"]] != "dominado" || p.Concluidos != 1 || p.Total != 4 {
		t.Fatalf("progresso: %+v", p)
	}
	if c.gam.get(c.user).XPTotal != 20 {
		t.Fatalf("xp = %d", c.gam.get(c.user).XPTotal)
	}
	// chamar de novo não duplica
	if _, err := c.aluno.ObterProgresso(ctx, c.user, "go-basico"); err != nil {
		t.Fatal(err)
	}
	if c.gam.get(c.user).XPTotal != 20 {
		t.Fatalf("xp duplicado = %d", c.gam.get(c.user).XPTotal)
	}
}

func TestRoadmapAluno_Progresso_EmCurso(t *testing.T) {
	c := novoCenario(true)
	c.prog.artigos[c.rm.Nos[3].ArtigoID] = false // abriu, mas não concluiu
	p, err := c.aluno.ObterProgresso(context.Background(), c.user, "go-basico")
	if err != nil {
		t.Fatal(err)
	}
	if p.Estados[c.nos["c"]] != "em_curso" {
		t.Fatalf("estados: %v", p.Estados)
	}
	if p.Estados[c.nos["a"]] != "disponivel" || p.Estados[c.nos["b"]] != "bloqueado" {
		t.Fatalf("estados: %v", p.Estados)
	}
}

func TestGamificacao_AoConcluirArtigo_ConcluiNoVinculadoERegistraStreak(t *testing.T) {
	c := novoCenario(true)
	res, err := c.svc.AoConcluirArtigo(context.Background(), c.user, c.rm.Nos[3].ArtigoID)
	if err != nil {
		t.Fatal(err)
	}
	if res.XPGanho != 20 || len(res.NosConcluidos) != 1 || res.NosConcluidos[0] != c.nos["c"] {
		t.Fatalf("resultado: %+v", res)
	}
	if res.Streak != 1 {
		t.Fatalf("streak = %d", res.Streak)
	}
}

func TestGamificacao_ArtigoSemNoApenasRegistraAtividade(t *testing.T) {
	c := novoCenario(true)
	res, err := c.svc.AoConcluirArtigo(context.Background(), c.user, "00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if res.XPGanho != 0 || len(res.NosConcluidos) != 0 || res.Streak != 1 {
		t.Fatalf("resultado: %+v", res)
	}
}

func TestGamificacao_ArtigoDeNoBloqueadoNaoAvanca(t *testing.T) {
	c := novoCenario(true)
	// faz o nó "c" depender de "a" (ainda não dominado)
	c.rm.Arestas = append(c.rm.Arestas, entity.RoadmapAresta{OrigemID: c.nos["a"], DestinoID: c.nos["c"], Tipo: entity.ArestaRequer})
	res, err := c.svc.AoConcluirArtigo(context.Background(), c.user, c.rm.Nos[3].ArtigoID)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.NosConcluidos) != 0 || res.XPGanho != 0 {
		t.Fatalf("nó bloqueado não deveria concluir: %+v", res)
	}
}

func TestGamificacao_RoadmapCompletoDaBonusECartografo(t *testing.T) {
	c := novoCenario(true)
	ctx := context.Background()
	prepararChefe(t, c)
	if _, err := c.aluno.ResponderQuiz(ctx, c.user, "go-basico", c.nos["chefe"], map[string]string{"q1": "x", "q2": "y"}); err != nil {
		t.Fatal(err)
	}
	res, err := c.svc.AoConcluirArtigo(ctx, c.user, c.rm.Nos[3].ArtigoID)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.RoadmapsCompletos) != 1 || res.XPGanho != 20+entity.XPBonusRoadmap {
		t.Fatalf("resultado: %+v", res)
	}
	achou := false
	for _, q := range res.ConquistasNovas {
		if q.Codigo == "cartografo" {
			achou = true
		}
	}
	if !achou {
		t.Fatalf("faltou cartógrafo: %+v", res.ConquistasNovas)
	}
}

func TestGamificacao_Resumo(t *testing.T) {
	c := novoCenario(true)
	if _, err := c.aluno.Concluir(context.Background(), c.user, "go-basico", c.nos["a"]); err != nil {
		t.Fatal(err)
	}
	r, err := c.aluno.Gamificacao(context.Background(), c.user)
	if err != nil {
		t.Fatal(err)
	}
	if r.XPTotal != 10 || r.Nivel != 1 || r.XPNoNivel != 10 || r.XPParaProximo != 50 || r.StreakAtual != 1 || !r.AtivoHoje {
		t.Fatalf("resumo: %+v", r)
	}
	if len(r.Conquistas) != len(entity.CatalogoConquistas) {
		t.Fatalf("catálogo: %d", len(r.Conquistas))
	}
	obtidas := 0
	for _, q := range r.Conquistas {
		if q.Conquistada {
			obtidas++
		}
	}
	if obtidas != 1 {
		t.Fatalf("conquistadas = %d", obtidas)
	}
}

func TestObterRoadmap_PublicoNaoVazaGabaritoERascunho(t *testing.T) {
	c := novoCenario(true)
	resp, err := NewObterRoadmap(c.repo).Execute(context.Background(), "go-basico")
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range resp.Nos {
		if n.Quiz == nil {
			continue
		}
		for _, q := range n.Quiz.Questoes {
			if q.Correta != "" || q.Explicacao != "" {
				t.Fatalf("gabarito vazou: %+v", q)
			}
		}
	}
	rascunho := novoCenario(false)
	if _, err := NewObterRoadmap(rascunho.repo).Execute(context.Background(), "go-basico"); kind(t, err) != domainErros.NotFound {
		t.Fatalf("rascunho deveria ser 404: %v", err)
	}
}

func TestListarRoadmaps_SoPublicados(t *testing.T) {
	pub := novoCenario(true)
	resp, err := NewListarRoadmaps(pub.repo).Execute(context.Background())
	if err != nil || len(resp.Itens) != 1 || resp.Itens[0].TotalChefes != 1 || resp.Itens[0].TotalNos != 4 || resp.Itens[0].XPTotal != 140 {
		t.Fatalf("itens: %+v err=%v", resp, err)
	}
	rasc := novoCenario(false)
	resp, _ = NewListarRoadmaps(rasc.repo).Execute(context.Background())
	if len(resp.Itens) != 0 {
		t.Fatalf("rascunho no catálogo: %+v", resp)
	}
}
