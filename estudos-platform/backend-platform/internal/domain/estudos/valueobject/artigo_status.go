package valueobject

import "github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"

type ArtigoStatus string

const (
	ArtigoStatusRascunho  ArtigoStatus = "rascunho"
	ArtigoStatusRevisao   ArtigoStatus = "revisao"
	ArtigoStatusPublicado ArtigoStatus = "publicado"
	ArtigoStatusArquivado ArtigoStatus = "arquivado"
)

func ParseArtigoStatus(raw string) (ArtigoStatus, error) {
	s := ArtigoStatus(raw)
	switch s {
	case ArtigoStatusRascunho, ArtigoStatusRevisao, ArtigoStatusPublicado, ArtigoStatusArquivado:
		return s, nil
	default:
		return "", errors.ErrInvalidArgument("status de artigo inválido", "valueobject.ParseArtigoStatus", nil)
	}
}

func (s ArtigoStatus) String() string { return string(s) }
