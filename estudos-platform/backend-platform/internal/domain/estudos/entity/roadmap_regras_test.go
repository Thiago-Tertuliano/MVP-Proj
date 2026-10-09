package entity

import (
	"testing"

	"github.com/google/uuid"
)

func no(id string) RoadmapNo {
	return RoadmapNo{ID: id, Tipo: NoTopico, Titulo: "Nó " + id[:4], XP: 10, Conclusao: ConclusaoManual}
}

func quizValido() *Quiz {
	return &Quiz{Questoes: []QuestaoQuiz{{
		ID: "q1", Enunciado: "2+2?",
		Opcoes:  []OpcaoQuiz{{ID: "a", Texto: "3"}, {ID: "b", Texto: "4"}},
		Correta: "b",
	}}}
}

func TestValidarGrafo(t *testing.T) {
	a, b, c := uuid.NewString(), uuid.NewString(), uuid.NewString()

	casos := []struct {
		nome    string
		nos     []RoadmapNo
		arestas []RoadmapAresta
		erro    bool
	}{
		{"vazio é válido", nil, nil, false},
		{"cadeia simples", []RoadmapNo{no(a), no(b)}, []RoadmapAresta{{OrigemID: a, DestinoID: b, Tipo: ArestaRequer}}, false},
		{"nó duplicado", []RoadmapNo{no(a), no(a)}, nil, true},
		{"id inválido", []RoadmapNo{{ID: "x", Tipo: NoTopico, Titulo: "abc", Conclusao: ConclusaoManual}}, nil, true},
		{"aresta para nó inexistente", []RoadmapNo{no(a)}, []RoadmapAresta{{OrigemID: a, DestinoID: b, Tipo: ArestaRequer}}, true},
		{"auto-laço", []RoadmapNo{no(a)}, []RoadmapAresta{{OrigemID: a, DestinoID: a, Tipo: ArestaRequer}}, true},
		{"aresta duplicada", []RoadmapNo{no(a), no(b)}, []RoadmapAresta{
			{OrigemID: a, DestinoID: b, Tipo: ArestaRequer}, {OrigemID: a, DestinoID: b, Tipo: ArestaOpcional},
		}, true},
		{"ciclo em requer", []RoadmapNo{no(a), no(b), no(c)}, []RoadmapAresta{
			{OrigemID: a, DestinoID: b, Tipo: ArestaRequer},
			{OrigemID: b, DestinoID: c, Tipo: ArestaRequer},
			{OrigemID: c, DestinoID: a, Tipo: ArestaRequer},
		}, true},
		{"ciclo via opcional é ignorado", []RoadmapNo{no(a), no(b)}, []RoadmapAresta{
			{OrigemID: a, DestinoID: b, Tipo: ArestaRequer},
			{OrigemID: b, DestinoID: a, Tipo: ArestaOpcional},
		}, false},
		{"tipo de aresta inválido", []RoadmapNo{no(a), no(b)}, []RoadmapAresta{{OrigemID: a, DestinoID: b, Tipo: "x"}}, true},
		{"xp acima do limite", []RoadmapNo{{ID: a, Tipo: NoTopico, Titulo: "abc", XP: 5000, Conclusao: ConclusaoManual}}, nil, true},
		{"conclusão por artigo sem artigo", []RoadmapNo{{ID: a, Tipo: NoTopico, Titulo: "abc", Conclusao: ConclusaoArtigo}}, nil, true},
		{"conclusão por artigo com artigo", []RoadmapNo{{ID: a, Tipo: NoTopico, Titulo: "abc", Conclusao: ConclusaoArtigo, ArtigoID: uuid.NewString()}}, nil, false},
		{"chefe sem quiz", []RoadmapNo{{ID: a, Tipo: NoChefe, Titulo: "abc", Conclusao: ConclusaoQuiz}}, nil, true},
		{"chefe manual", []RoadmapNo{{ID: a, Tipo: NoChefe, Titulo: "abc", Conclusao: ConclusaoManual}}, nil, true},
		{"chefe com quiz", []RoadmapNo{{ID: a, Tipo: NoChefe, Titulo: "abc", Conclusao: ConclusaoQuiz, Quiz: quizValido()}}, nil, false},
		{"quiz em tópico", []RoadmapNo{{ID: a, Tipo: NoTopico, Titulo: "abc", Conclusao: ConclusaoQuiz, Quiz: quizValido()}}, nil, true},
		{"quiz com correta fora das opções", []RoadmapNo{{ID: a, Tipo: NoChefe, Titulo: "abc", Conclusao: ConclusaoQuiz, Quiz: &Quiz{Questoes: []QuestaoQuiz{{
			ID: "q1", Enunciado: "?", Opcoes: []OpcaoQuiz{{ID: "a", Texto: "x"}, {ID: "b", Texto: "y"}}, Correta: "z",
		}}}}}, nil, true},
		{"quiz com uma só opção", []RoadmapNo{{ID: a, Tipo: NoChefe, Titulo: "abc", Conclusao: ConclusaoQuiz, Quiz: &Quiz{Questoes: []QuestaoQuiz{{
			ID: "q1", Enunciado: "?", Opcoes: []OpcaoQuiz{{ID: "a", Texto: "x"}}, Correta: "a",
		}}}}}, nil, true},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			err := ValidarGrafo(c.nos, c.arestas)
			if (err != nil) != c.erro {
				t.Fatalf("erro = %v, esperado erro=%v", err, c.erro)
			}
		})
	}
}

func TestValidarEstrutura_AceitaRascunhoIncompleto(t *testing.T) {
	a := uuid.NewString()
	chefeSemQuiz := []RoadmapNo{{ID: a, Tipo: NoChefe, Titulo: "Chefe", Conclusao: ConclusaoQuiz}}
	if err := ValidarEstrutura(chefeSemQuiz, nil); err != nil {
		t.Fatalf("estrutura deveria aceitar: %v", err)
	}
	if err := ValidarGrafo(chefeSemQuiz, nil); err == nil {
		t.Fatal("grafo completo deveria rejeitar chefe sem quiz")
	}
	// estrutura ainda rejeita ciclo e id inválido
	b := uuid.NewString()
	ciclo := []RoadmapAresta{{OrigemID: a, DestinoID: b, Tipo: ArestaRequer}, {OrigemID: b, DestinoID: a, Tipo: ArestaRequer}}
	if err := ValidarEstrutura([]RoadmapNo{no(a), no(b)}, ciclo); err == nil {
		t.Fatal("ciclo deveria falhar mesmo em rascunho")
	}
	if err := ValidarEstrutura([]RoadmapNo{{ID: "x", Tipo: NoTopico, Titulo: "abc", Conclusao: ConclusaoManual}}, nil); err == nil {
		t.Fatal("id inválido deveria falhar")
	}
}

func TestValidarGrafo_LimiteDeNos(t *testing.T) {
	nos := make([]RoadmapNo, MaxNosPorRoadmap+1)
	for i := range nos {
		nos[i] = no(uuid.NewString())
	}
	if err := ValidarGrafo(nos, nil); err == nil {
		t.Fatal("esperava erro acima do limite de nós")
	}
}

func TestCalcularEstados(t *testing.T) {
	a, b, c, d := uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString()
	nos := []RoadmapNo{no(a), no(b), no(c), no(d)}
	arestas := []RoadmapAresta{
		{OrigemID: a, DestinoID: b, Tipo: ArestaRequer},
		{OrigemID: a, DestinoID: c, Tipo: ArestaRequer},
		{OrigemID: b, DestinoID: d, Tipo: ArestaRequer},
		{OrigemID: c, DestinoID: d, Tipo: ArestaOpcional},
	}

	t.Run("nada concluído", func(t *testing.T) {
		e := CalcularEstados(nos, arestas, nil, nil)
		assertEstado(t, e, a, EstadoDisponivel)
		assertEstado(t, e, b, EstadoBloqueado)
		assertEstado(t, e, c, EstadoBloqueado)
		assertEstado(t, e, d, EstadoBloqueado)
	})

	t.Run("raiz dominada libera filhos", func(t *testing.T) {
		e := CalcularEstados(nos, arestas, map[string]bool{a: true}, nil)
		assertEstado(t, e, a, EstadoDominado)
		assertEstado(t, e, b, EstadoDisponivel)
		assertEstado(t, e, c, EstadoDisponivel)
		assertEstado(t, e, d, EstadoBloqueado)
	})

	t.Run("aresta opcional não bloqueia", func(t *testing.T) {
		e := CalcularEstados(nos, arestas, map[string]bool{a: true, b: true}, nil)
		assertEstado(t, e, d, EstadoDisponivel)
	})

	t.Run("em curso só vale para nó disponível", func(t *testing.T) {
		e := CalcularEstados(nos, arestas, map[string]bool{a: true}, map[string]bool{b: true, d: true})
		assertEstado(t, e, b, EstadoEmCurso)
		assertEstado(t, e, d, EstadoBloqueado)
	})

	t.Run("dominado prevalece", func(t *testing.T) {
		e := CalcularEstados(nos, arestas, map[string]bool{a: true}, map[string]bool{a: true})
		assertEstado(t, e, a, EstadoDominado)
	})
}

func assertEstado(t *testing.T, e map[string]EstadoNo, id string, want EstadoNo) {
	t.Helper()
	if e[id] != want {
		t.Fatalf("estado de %s = %s, esperado %s", id[:4], e[id], want)
	}
}

func TestAvaliarQuiz(t *testing.T) {
	q := &Quiz{Questoes: []QuestaoQuiz{
		{ID: "1", Enunciado: "a", Opcoes: []OpcaoQuiz{{ID: "x", Texto: "x"}, {ID: "y", Texto: "y"}}, Correta: "x", Explicacao: "porque x"},
		{ID: "2", Enunciado: "b", Opcoes: []OpcaoQuiz{{ID: "x", Texto: "x"}, {ID: "y", Texto: "y"}}, Correta: "y"},
		{ID: "3", Enunciado: "c", Opcoes: []OpcaoQuiz{{ID: "x", Texto: "x"}, {ID: "y", Texto: "y"}}, Correta: "y"},
		{ID: "4", Enunciado: "d", Opcoes: []OpcaoQuiz{{ID: "x", Texto: "x"}, {ID: "y", Texto: "y"}}, Correta: "x"},
	}}

	t.Run("todas certas", func(t *testing.T) {
		r := AvaliarQuiz(q, map[string]string{"1": "x", "2": "y", "3": "y", "4": "x"})
		if !r.Aprovado || r.Acertos != 4 || r.Total != 4 {
			t.Fatalf("resultado inesperado: %+v", r)
		}
	})
	t.Run("3 de 4 (75%) aprova", func(t *testing.T) {
		r := AvaliarQuiz(q, map[string]string{"1": "x", "2": "y", "3": "y", "4": "y"})
		if !r.Aprovado || r.Acertos != 3 {
			t.Fatalf("resultado inesperado: %+v", r)
		}
	})
	t.Run("2 de 4 (50%) reprova e mantém explicação", func(t *testing.T) {
		r := AvaliarQuiz(q, map[string]string{"1": "y", "2": "y", "3": "y", "4": "y"})
		if r.Aprovado || r.Acertos != 2 {
			t.Fatalf("resultado inesperado: %+v", r)
		}
		if r.Correcao[0].Correta || r.Correcao[0].Explicacao != "porque x" {
			t.Fatalf("correção inesperada: %+v", r.Correcao[0])
		}
	})
	t.Run("sem respostas reprova", func(t *testing.T) {
		if r := AvaliarQuiz(q, nil); r.Aprovado || r.Acertos != 0 {
			t.Fatalf("resultado inesperado: %+v", r)
		}
	})
	t.Run("quiz nil nunca aprova", func(t *testing.T) {
		if r := AvaliarQuiz(nil, nil); r.Aprovado {
			t.Fatal("quiz nil não deve aprovar")
		}
	})
}

func TestRoadmap_PublicarExigeNosEGrafoValido(t *testing.T) {
	r, err := NovoRoadmap("go", "Go do zero", "", "", uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Publicar(); err == nil {
		t.Fatal("roadmap vazio não deve publicar")
	}
	a := uuid.NewString()
	r.Nos = []RoadmapNo{no(a)}
	if err := r.Publicar(); err != nil {
		t.Fatalf("publicar: %v", err)
	}
	if !r.Publicado {
		t.Fatal("deveria estar publicado")
	}
	r.Despublicar()
	if r.Publicado {
		t.Fatal("deveria voltar a rascunho")
	}
}

func TestNovoRoadmap_Validacoes(t *testing.T) {
	if _, err := NovoRoadmap("x", "ab", "", "", uuid.NewString()); err == nil {
		t.Fatal("título curto deveria falhar")
	}
	if _, err := NovoRoadmap("xyz", "Título ok", "", "", "nao-uuid"); err == nil {
		t.Fatal("autor inválido deveria falhar")
	}
}
