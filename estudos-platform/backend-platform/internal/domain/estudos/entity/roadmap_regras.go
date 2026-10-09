package entity

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// EstadoNo é o estado de um nó para um aluno (calculado no servidor).
type EstadoNo string

const (
	EstadoBloqueado  EstadoNo = "bloqueado"
	EstadoDisponivel EstadoNo = "disponivel"
	EstadoEmCurso    EstadoNo = "em_curso"
	EstadoDominado   EstadoNo = "dominado"
)

func invalido(msg string) error {
	return errors.ErrInvalidArgument(msg, "entity.ValidarGrafo", nil)
}

// ValidarGrafo aplica todas as invariantes (estrutura + conteúdo). Usado ao publicar e ao
// salvar um roadmap que já está publicado.
func ValidarGrafo(nos []RoadmapNo, arestas []RoadmapAresta) error {
	return validarGrafo(nos, arestas, true)
}

// ValidarEstrutura aplica só as invariantes estruturais (ids, conexões, ciclos, limites),
// permitindo salvar rascunhos incompletos (ex.: chefe ainda sem quiz).
func ValidarEstrutura(nos []RoadmapNo, arestas []RoadmapAresta) error {
	return validarGrafo(nos, arestas, false)
}

func validarGrafo(nos []RoadmapNo, arestas []RoadmapAresta, conteudo bool) error {
	if len(nos) > MaxNosPorRoadmap {
		return invalido(fmt.Sprintf("máximo de %d nós por roadmap", MaxNosPorRoadmap))
	}

	ids := make(map[string]struct{}, len(nos))
	for i := range nos {
		if err := validarNo(&nos[i], conteudo); err != nil {
			return err
		}
		if _, dup := ids[nos[i].ID]; dup {
			return invalido("nó duplicado no roadmap")
		}
		ids[nos[i].ID] = struct{}{}
	}

	vistas := make(map[string]struct{}, len(arestas))
	for _, a := range arestas {
		if a.Tipo != ArestaRequer && a.Tipo != ArestaOpcional {
			return invalido("tipo de conexão inválido")
		}
		if a.OrigemID == a.DestinoID {
			return invalido("um nó não pode se conectar a ele mesmo")
		}
		if _, ok := ids[a.OrigemID]; !ok {
			return invalido("conexão aponta para um nó inexistente")
		}
		if _, ok := ids[a.DestinoID]; !ok {
			return invalido("conexão aponta para um nó inexistente")
		}
		chave := a.OrigemID + ">" + a.DestinoID
		if _, dup := vistas[chave]; dup {
			return invalido("conexão duplicada")
		}
		vistas[chave] = struct{}{}
	}

	if TemCiclo(nos, arestas) {
		return invalido("o roadmap não pode ter ciclos de pré-requisito")
	}
	return nil
}

func validarNo(n *RoadmapNo, conteudo bool) error {
	if _, err := uuid.Parse(n.ID); err != nil {
		return invalido("id de nó inválido")
	}
	titulo := strings.TrimSpace(n.Titulo)
	if len(titulo) < 2 || len(titulo) > MaxTituloNo {
		return invalido("título do nó deve ter entre 2 e 120 caracteres")
	}
	if len(n.Descricao) > MaxDescricaoNo {
		return invalido("descrição do nó muito longa")
	}
	switch n.Tipo {
	case NoTopico, NoChefe, NoMarco:
	default:
		return invalido("tipo de nó inválido")
	}
	if n.XP < 0 || n.XP > MaxXPPorNo {
		return invalido("XP do nó deve estar entre 0 e 1000")
	}
	switch n.Conclusao {
	case ConclusaoManual, ConclusaoArtigo, ConclusaoQuiz:
	default:
		return invalido("forma de conclusão inválida")
	}
	if n.ArtigoID != "" {
		if _, err := uuid.Parse(n.ArtigoID); err != nil {
			return invalido("artigo vinculado inválido")
		}
	}
	if !conteudo {
		return nil
	}
	if n.Conclusao == ConclusaoArtigo && n.ArtigoID == "" {
		return invalido("nó com conclusão por artigo precisa de um artigo vinculado")
	}
	if n.Tipo == NoChefe && n.Conclusao != ConclusaoQuiz {
		return invalido("nó chefe deve ser concluído por quiz")
	}
	if n.Conclusao == ConclusaoQuiz && n.Tipo != NoChefe {
		return invalido("somente nó chefe usa quiz")
	}
	if n.Conclusao == ConclusaoQuiz {
		return validarQuiz(n.Quiz)
	}
	return nil
}

func validarQuiz(q *Quiz) error {
	if q == nil || len(q.Questoes) == 0 {
		return invalido("o quiz do chefe precisa de ao menos uma questão")
	}
	if len(q.Questoes) > MaxQuestoesPorQuiz {
		return invalido(fmt.Sprintf("máximo de %d questões por quiz", MaxQuestoesPorQuiz))
	}
	qIDs := map[string]struct{}{}
	for _, qs := range q.Questoes {
		if strings.TrimSpace(qs.ID) == "" {
			return invalido("questão sem identificador")
		}
		if _, dup := qIDs[qs.ID]; dup {
			return invalido("questão duplicada")
		}
		qIDs[qs.ID] = struct{}{}
		if e := strings.TrimSpace(qs.Enunciado); e == "" || len(e) > MaxEnunciadoQuestao {
			return invalido("enunciado da questão inválido")
		}
		if len(qs.Opcoes) < MinOpcoesPorQuestao || len(qs.Opcoes) > MaxOpcoesPorQuestao {
			return invalido("cada questão deve ter entre 2 e 6 opções")
		}
		oIDs := map[string]struct{}{}
		for _, o := range qs.Opcoes {
			if strings.TrimSpace(o.ID) == "" {
				return invalido("opção sem identificador")
			}
			if _, dup := oIDs[o.ID]; dup {
				return invalido("opção duplicada")
			}
			oIDs[o.ID] = struct{}{}
			if t := strings.TrimSpace(o.Texto); t == "" || len(t) > MaxTextoOpcaoQuestao {
				return invalido("texto da opção inválido")
			}
		}
		if _, ok := oIDs[qs.Correta]; !ok {
			return invalido("a resposta correta deve ser uma das opções")
		}
	}
	return nil
}

// TemCiclo detecta ciclos considerando apenas arestas "requer" (Kahn).
// Arestas opcionais não bloqueiam e, portanto, não participam da detecção.
func TemCiclo(nos []RoadmapNo, arestas []RoadmapAresta) bool {
	grau := make(map[string]int, len(nos))
	adj := make(map[string][]string, len(nos))
	for _, n := range nos {
		grau[n.ID] = 0
	}
	for _, a := range arestas {
		if a.Tipo != ArestaRequer {
			continue
		}
		adj[a.OrigemID] = append(adj[a.OrigemID], a.DestinoID)
		grau[a.DestinoID]++
	}
	fila := make([]string, 0, len(nos))
	for id, g := range grau {
		if g == 0 {
			fila = append(fila, id)
		}
	}
	visitados := 0
	for len(fila) > 0 {
		id := fila[0]
		fila = fila[1:]
		visitados++
		for _, d := range adj[id] {
			grau[d]--
			if grau[d] == 0 {
				fila = append(fila, d)
			}
		}
	}
	return visitados != len(grau)
}

// CalcularEstados determina o estado de cada nó para um aluno.
//   - dominado: nó concluído;
//   - disponível: todos os pré-requisitos "requer" dominados (ou sem pré-requisitos);
//   - em curso: disponível e com atividade iniciada (ex.: artigo aberto);
//   - bloqueado: algum pré-requisito "requer" ainda não dominado.
func CalcularEstados(nos []RoadmapNo, arestas []RoadmapAresta, concluidos, emCurso map[string]bool) map[string]EstadoNo {
	requisitos := make(map[string][]string, len(nos))
	for _, a := range arestas {
		if a.Tipo == ArestaRequer {
			requisitos[a.DestinoID] = append(requisitos[a.DestinoID], a.OrigemID)
		}
	}
	estados := make(map[string]EstadoNo, len(nos))
	for _, n := range nos {
		if concluidos[n.ID] {
			estados[n.ID] = EstadoDominado
			continue
		}
		livre := true
		for _, pre := range requisitos[n.ID] {
			if !concluidos[pre] {
				livre = false
				break
			}
		}
		switch {
		case !livre:
			estados[n.ID] = EstadoBloqueado
		case emCurso[n.ID]:
			estados[n.ID] = EstadoEmCurso
		default:
			estados[n.ID] = EstadoDisponivel
		}
	}
	return estados
}

// CorrecaoQuestao é o feedback de uma questão após o envio do quiz.
type CorrecaoQuestao struct {
	QuestaoID  string
	Correta    bool
	Explicacao string
}

// ResultadoQuiz é a avaliação de um envio.
type ResultadoQuiz struct {
	Acertos  int
	Total    int
	Aprovado bool
	Correcao []CorrecaoQuestao
}

// AvaliarQuiz corrige as respostas (questaoID → opcaoID). Sem resposta = erro.
// Aprovação: acertos >= NotaMinimaQuizPct% do total.
func AvaliarQuiz(q *Quiz, respostas map[string]string) ResultadoQuiz {
	res := ResultadoQuiz{}
	if q == nil {
		return res
	}
	res.Total = len(q.Questoes)
	for _, qs := range q.Questoes {
		ok := respostas[qs.ID] != "" && respostas[qs.ID] == qs.Correta
		if ok {
			res.Acertos++
		}
		res.Correcao = append(res.Correcao, CorrecaoQuestao{QuestaoID: qs.ID, Correta: ok, Explicacao: qs.Explicacao})
	}
	res.Aprovado = res.Total > 0 && res.Acertos*100 >= NotaMinimaQuizPct*res.Total
	return res
}
