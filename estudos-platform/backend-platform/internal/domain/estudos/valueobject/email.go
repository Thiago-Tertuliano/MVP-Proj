package valueobject

import (
	"net/mail"
	"strings"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// Email é um Value Object imutável com validação na construção.
type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	addr, err := mail.ParseAddress(raw)
	if err != nil || addr.Address != raw {
		return Email{}, errors.ErrInvalidArgument("email inválido", "valueobject.NewEmail", nil)
	}
	return Email{value: raw}, nil
}

// ReconstructEmail reconstrói sem validar (dado já veio do banco).
func ReconstructEmail(value string) Email {
	return Email{value: value}
}

func (e Email) Value() string       { return e.value }
func (e Email) String() string      { return e.value }
func (e Email) Equals(o Email) bool { return e.value == o.value }
