package entity

import (
	"math"
	"time"
)

// XPBaseNivel define a curva: nível = floor(sqrt(xp/XPBaseNivel)) + 1.
const XPBaseNivel = 50

// Gamificacao é o resumo de XP/nível/streak de um usuário.
type Gamificacao struct {
	XPTotal     int
	StreakAtual int
	StreakMax   int
	UltimoDia   *time.Time // data (00:00 UTC) do último dia com atividade
}

// NivelPorXP calcula o nível a partir do XP total (mínimo 1).
func NivelPorXP(xp int) int {
	if xp < 0 {
		xp = 0
	}
	return int(math.Floor(math.Sqrt(float64(xp)/XPBaseNivel))) + 1
}

// XPInicioNivel é o XP total necessário para alcançar o nível n.
func XPInicioNivel(n int) int {
	if n <= 1 {
		return 0
	}
	return XPBaseNivel * (n - 1) * (n - 1)
}

// ProgressoNivel devolve (xp ganho dentro do nível, xp necessário para o próximo).
func ProgressoNivel(xp int) (noNivel, paraProximo int) {
	n := NivelPorXP(xp)
	inicio := XPInicioNivel(n)
	return xp - inicio, XPInicioNivel(n+1) - inicio
}

// fusoBrasil: o Brasil não adota horário de verão desde 2019 (UTC-3 fixo), o que
// evita depender do banco tzdata em imagens mínimas.
var fusoBrasil = time.FixedZone("America/Sao_Paulo", -3*60*60)

// DiaCalendario normaliza um instante para a data (00:00 UTC) no fuso do Brasil.
func DiaCalendario(t time.Time) time.Time {
	l := t.In(fusoBrasil)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, time.UTC)
}

// soData descarta a hora de uma data já normalizada (armazenada como DATE).
func soData(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// AtualizarStreak aplica a regra de sequência diária para uma atividade em `hoje`.
//   - mesmo dia: mantém (inicia em 1 se estava zerado);
//   - dia seguinte: +1;
//   - lacuna > 1 dia: reinicia em 1.
func AtualizarStreak(ultimo *time.Time, hoje time.Time, atual, max int) (novoAtual, novoMax int) {
	dia := DiaCalendario(hoje)
	switch {
	case ultimo == nil:
		novoAtual = 1
	default:
		diff := int(dia.Sub(soData(*ultimo)).Hours() / 24)
		switch {
		case diff < 0:
			novoAtual = atual
		case diff == 0:
			novoAtual = atual
			if novoAtual == 0 {
				novoAtual = 1
			}
		case diff == 1:
			novoAtual = atual + 1
		default:
			novoAtual = 1
		}
	}
	novoMax = max
	if novoAtual > novoMax {
		novoMax = novoAtual
	}
	return novoAtual, novoMax
}

// StreakVigente informa a sequência exibida hoje: se o último dia foi antes de ontem, a
// sequência foi quebrada (0) mesmo que ainda esteja gravada no banco.
func StreakVigente(g Gamificacao, agora time.Time) int {
	if g.UltimoDia == nil {
		return 0
	}
	diff := int(DiaCalendario(agora).Sub(soData(*g.UltimoDia)).Hours() / 24)
	if diff > 1 {
		return 0
	}
	return g.StreakAtual
}

// Conquista é uma medalha do catálogo.
type Conquista struct {
	Codigo    string
	Nome      string
	Descricao string
	Icone     string
}

// CatalogoConquistas lista as conquistas disponíveis (ordem de exibição).
var CatalogoConquistas = []Conquista{
	{"primeiro-passo", "Primeiro passo", "Dominou seu primeiro nó de um roadmap.", "footprints"},
	{"explorador", "Explorador", "Dominou 10 nós.", "compass"},
	{"mestre-dos-nos", "Mestre dos nós", "Dominou 50 nós.", "network"},
	{"cacador-de-chefes", "Caçador de chefes", "Derrotou seu primeiro chefe.", "swords"},
	{"lenda", "Lenda", "Derrotou 5 chefes.", "crown"},
	{"chama-acesa", "Chama acesa", "3 dias seguidos de estudo.", "flame"},
	{"semana-de-fogo", "Semana de fogo", "7 dias seguidos de estudo.", "flame"},
	{"mes-imparavel", "Mês imparável", "30 dias seguidos de estudo.", "flame"},
	{"nivel-5", "Nível 5", "Alcançou o nível 5.", "star"},
	{"nivel-10", "Nível 10", "Alcançou o nível 10.", "star"},
	{"cartografo", "Cartógrafo", "Completou um roadmap inteiro.", "map"},
}

// ConquistaPorCodigo busca no catálogo.
func ConquistaPorCodigo(codigo string) (Conquista, bool) {
	for _, c := range CatalogoConquistas {
		if c.Codigo == codigo {
			return c, true
		}
	}
	return Conquista{}, false
}

// SnapshotConquistas reúne os números usados para avaliar as conquistas.
type SnapshotConquistas struct {
	XPTotal          int
	StreakMax        int
	NosConcluidos    int
	ChefesConcluidos int
	RoadmapCompleto  bool // verdadeiro se o evento atual completou um roadmap
}

// ConquistasElegiveis devolve os códigos que o snapshot satisfaz (o chamador filtra as já obtidas).
func ConquistasElegiveis(s SnapshotConquistas) []string {
	var out []string
	add := func(cond bool, codigo string) {
		if cond {
			out = append(out, codigo)
		}
	}
	nivel := NivelPorXP(s.XPTotal)
	add(s.NosConcluidos >= 1, "primeiro-passo")
	add(s.NosConcluidos >= 10, "explorador")
	add(s.NosConcluidos >= 50, "mestre-dos-nos")
	add(s.ChefesConcluidos >= 1, "cacador-de-chefes")
	add(s.ChefesConcluidos >= 5, "lenda")
	add(s.StreakMax >= 3, "chama-acesa")
	add(s.StreakMax >= 7, "semana-de-fogo")
	add(s.StreakMax >= 30, "mes-imparavel")
	add(nivel >= 5, "nivel-5")
	add(nivel >= 10, "nivel-10")
	add(s.RoadmapCompleto, "cartografo")
	return out
}
