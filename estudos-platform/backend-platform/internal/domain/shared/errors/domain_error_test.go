package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestConstructors_KindCorreto(t *testing.T) {
	casos := []struct {
		name string
		got  *DomainError
		kind Kind
	}{
		{"InvalidArgument", ErrInvalidArgument("m", "op", nil), InvalidArgument},
		{"NotFound", ErrNotFound("m", "op", nil), NotFound},
		{"AlreadyExists", ErrAlreadyExists("m", "op", nil), AlreadyExists},
		{"InvalidState", ErrInvalidState("m", "op", nil), InvalidState},
		{"Unauthorized", ErrUnauthorized("m", "op", nil), Unauthorized},
		{"Internal", ErrInternal("m", "op", nil), Internal},
	}
	for _, c := range casos {
		t.Run(c.name, func(t *testing.T) {
			if c.got.Kind != c.kind {
				t.Errorf("esperava %s, got %s", c.kind, c.got.Kind)
			}
		})
	}
}

func TestError_FormatoComWrap(t *testing.T) {
	base := errors.New("causa raiz")
	de := ErrInternal("falhou", "op.Test", base)
	s := de.Error()
	if !strings.Contains(s, "INTERNAL") || !strings.Contains(s, "falhou") || !strings.Contains(s, "causa raiz") {
		t.Errorf("formato inesperado: %s", s)
	}
}

func TestError_Unwrap(t *testing.T) {
	base := errors.New("causa")
	de := ErrInternal("m", "op", base)
	if !errors.Is(de, base) {
		t.Error("errors.Is deveria encontrar a causa raiz")
	}
	if de.Unwrap() != base {
		t.Error("Unwrap deveria retornar a causa")
	}
}