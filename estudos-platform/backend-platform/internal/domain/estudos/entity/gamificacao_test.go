package entity

import (
	"testing"
	"time"
)

func TestNivelPorXP(t *testing.T) {
	casos := []struct{ xp, nivel int }{
		{-5, 1}, {0, 1}, {49, 1}, {50, 2}, {199, 2}, {200, 3}, {449, 3}, {450, 4}, {1250, 6}, {4050, 10},
	}
	for _, c := range casos {
		if got := NivelPorXP(c.xp); got != c.nivel {
			t.Errorf("NivelPorXP(%d) = %d, esperado %d", c.xp, got, c.nivel)
		}
	}
}

func TestXPInicioNivel_Consistente(t *testing.T) {
	for n := 1; n <= 30; n++ {
		inicio := XPInicioNivel(n)
		if NivelPorXP(inicio) != n {
			t.Fatalf("início do nível %d (%d xp) cai no nível %d", n, inicio, NivelPorXP(inicio))
		}
		if n > 1 && NivelPorXP(inicio-1) != n-1 {
			t.Fatalf("xp %d deveria ser nível %d", inicio-1, n-1)
		}
	}
}

func TestProgressoNivel(t *testing.T) {
	no, prox := ProgressoNivel(120) // nível 2: 50..200
	if no != 70 || prox != 150 {
		t.Fatalf("ProgressoNivel(120) = %d/%d, esperado 70/150", no, prox)
	}
}

func dia(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestAtualizarStreak(t *testing.T) {
	d := func(s string) *time.Time { v := DiaCalendario(dia(s)); return &v }

	t.Run("primeira atividade", func(t *testing.T) {
		a, m := AtualizarStreak(nil, dia("2026-10-08 12:00"), 0, 0)
		if a != 1 || m != 1 {
			t.Fatalf("got %d/%d", a, m)
		}
	})
	t.Run("mesmo dia mantém", func(t *testing.T) {
		a, m := AtualizarStreak(d("2026-10-08 12:00"), dia("2026-10-08 20:00"), 4, 6)
		if a != 4 || m != 6 {
			t.Fatalf("got %d/%d", a, m)
		}
	})
	t.Run("dia seguinte soma e atualiza máximo", func(t *testing.T) {
		a, m := AtualizarStreak(d("2026-10-07 12:00"), dia("2026-10-08 12:00"), 6, 6)
		if a != 7 || m != 7 {
			t.Fatalf("got %d/%d", a, m)
		}
	})
	t.Run("lacuna reinicia mas preserva máximo", func(t *testing.T) {
		a, m := AtualizarStreak(d("2026-10-05 12:00"), dia("2026-10-08 12:00"), 9, 9)
		if a != 1 || m != 9 {
			t.Fatalf("got %d/%d", a, m)
		}
	})
	t.Run("relógio retroativo não quebra", func(t *testing.T) {
		a, m := AtualizarStreak(d("2026-10-08 12:00"), dia("2026-10-07 12:00"), 3, 3)
		if a != 3 || m != 3 {
			t.Fatalf("got %d/%d", a, m)
		}
	})
}

func TestDiaCalendario_FusoBrasil(t *testing.T) {
	// 02:00 UTC de 09/10 ainda é 23:00 de 08/10 em Brasília.
	got := DiaCalendario(time.Date(2026, 10, 9, 2, 0, 0, 0, time.UTC))
	want := time.Date(2026, 10, 8, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("got %v, esperado %v", got, want)
	}
}

func TestStreakVigente(t *testing.T) {
	ult := DiaCalendario(dia("2026-10-06 12:00"))
	g := Gamificacao{StreakAtual: 5, UltimoDia: &ult}
	if v := StreakVigente(g, dia("2026-10-07 12:00")); v != 5 {
		t.Fatalf("ontem ainda vale: %d", v)
	}
	if v := StreakVigente(g, dia("2026-10-08 12:00")); v != 0 {
		t.Fatalf("quebrada deveria ser 0: %d", v)
	}
	if v := StreakVigente(Gamificacao{}, dia("2026-10-08 12:00")); v != 0 {
		t.Fatalf("sem atividade deveria ser 0: %d", v)
	}
}

func TestConquistasElegiveis(t *testing.T) {
	has := func(l []string, c string) bool {
		for _, x := range l {
			if x == c {
				return true
			}
		}
		return false
	}
	if got := ConquistasElegiveis(SnapshotConquistas{}); len(got) != 0 {
		t.Fatalf("snapshot zerado não deveria render nada: %v", got)
	}
	got := ConquistasElegiveis(SnapshotConquistas{XPTotal: 1250, StreakMax: 7, NosConcluidos: 10, ChefesConcluidos: 1, RoadmapCompleto: true})
	for _, c := range []string{"primeiro-passo", "explorador", "cacador-de-chefes", "chama-acesa", "semana-de-fogo", "nivel-5", "cartografo"} {
		if !has(got, c) {
			t.Errorf("faltou conquista %s em %v", c, got)
		}
	}
	for _, c := range []string{"mestre-dos-nos", "lenda", "mes-imparavel", "nivel-10"} {
		if has(got, c) {
			t.Errorf("não deveria ter %s", c)
		}
	}
}

func TestCatalogoConquistas_CodigosUnicos(t *testing.T) {
	vistos := map[string]bool{}
	for _, c := range CatalogoConquistas {
		if vistos[c.Codigo] {
			t.Fatalf("código duplicado: %s", c.Codigo)
		}
		vistos[c.Codigo] = true
		if _, ok := ConquistaPorCodigo(c.Codigo); !ok {
			t.Fatalf("ConquistaPorCodigo falhou para %s", c.Codigo)
		}
	}
}
