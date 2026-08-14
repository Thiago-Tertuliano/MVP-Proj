package valueobject

import (
	"regexp"
	"strings"

	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

var accentMap = map[rune]rune{
	'á': 'a', 'à': 'a', 'ã': 'a', 'â': 'a', 'ä': 'a',
	'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e',
	'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i',
	'ó': 'o', 'ò': 'o', 'õ': 'o', 'ô': 'o', 'ö': 'o',
	'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u',
	'ç': 'c', 'ñ': 'n',
}

// Slug é URL canônica (VO).
type Slug struct {
	value string
}

func NewSlug(raw string) (Slug, error) {
	normalized := normalizeSlug(raw)
	if len(normalized) < 2 || len(normalized) > 200 {
		return Slug{}, errors.ErrInvalidArgument("slug deve ter entre 2 e 200 caracteres", "valueobject.NewSlug", nil)
	}
	if !slugPattern.MatchString(normalized) {
		return Slug{}, errors.ErrInvalidArgument("slug inválido", "valueobject.NewSlug", nil)
	}
	return Slug{value: normalized}, nil
}

func ReconstructSlug(value string) Slug { return Slug{value: value} }

func (s Slug) Value() string      { return s.value }
func (s Slug) String() string     { return s.value }
func (s Slug) Equals(o Slug) bool { return s.value == o.value }

func normalizeSlug(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	var b strings.Builder
	prevHyphen := false
	for _, r := range raw {
		if mapped, ok := accentMap[r]; ok {
			r = mapped
		}
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevHyphen = false
		case r == ' ' || r == '_' || r == '-':
			if !prevHyphen && b.Len() > 0 {
				b.WriteByte('-')
				prevHyphen = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
