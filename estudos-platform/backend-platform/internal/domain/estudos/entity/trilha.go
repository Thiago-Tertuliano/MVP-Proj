package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/kernel"
)

// Modulo é entidade filha de Trilha (sem root próprio).
type Modulo struct {
	id        uuid.UUID
	slug      valueobject.Slug
	titulo    string
	descricao string
	ordem     int
	createdAt time.Time
}

func NovoModulo(slug valueobject.Slug, titulo, descricao string, ordem int) (*Modulo, error) {
	if len(titulo) < 2 || len(titulo) > 200 {
		return nil, errors.ErrInvalidArgument("título do módulo inválido", "entity.NovoModulo", nil)
	}
	return &Modulo{
		id:        uuid.New(),
		slug:      slug,
		titulo:    titulo,
		descricao: descricao,
		ordem:     ordem,
		createdAt: time.Now().UTC(),
	}, nil
}

func ReconstruirModulo(id uuid.UUID, slug valueobject.Slug, titulo, descricao string, ordem int, createdAt time.Time) *Modulo {
	return &Modulo{id: id, slug: slug, titulo: titulo, descricao: descricao, ordem: ordem, createdAt: createdAt}
}

func (m *Modulo) ID() uuid.UUID          { return m.id }
func (m *Modulo) Slug() valueobject.Slug { return m.slug }
func (m *Modulo) Titulo() string         { return m.titulo }
func (m *Modulo) Descricao() string      { return m.descricao }
func (m *Modulo) Ordem() int             { return m.ordem }
func (m *Modulo) CreatedAt() time.Time   { return m.createdAt }

// Trilha é Aggregate Root de trilhas de estudo.
type Trilha struct {
	kernel.BaseEntity
	slug      valueobject.Slug
	titulo    string
	descricao string
	capaURL   string
	ordem     int
	publicada bool
	modulos   []*Modulo
}

type NovaTrilhaInput struct {
	Slug      valueobject.Slug
	Titulo    string
	Descricao string
	CapaURL   string
	Ordem     int
}

func NovaTrilha(in NovaTrilhaInput) (*Trilha, error) {
	if len(in.Titulo) < 3 || len(in.Titulo) > 200 {
		return nil, errors.ErrInvalidArgument("título deve ter entre 3 e 200 caracteres", "entity.NovaTrilha", nil)
	}
	return &Trilha{
		BaseEntity: kernel.NewBaseEntity(),
		slug:       in.Slug,
		titulo:     in.Titulo,
		descricao:  in.Descricao,
		capaURL:    in.CapaURL,
		ordem:      in.Ordem,
		publicada:  false,
		modulos:    []*Modulo{},
	}, nil
}

func ReconstruirTrilha(
	id uuid.UUID,
	slug valueobject.Slug,
	titulo, descricao, capaURL string,
	ordem int,
	publicada bool,
	modulos []*Modulo,
	createdAt, updatedAt time.Time,
) *Trilha {
	if modulos == nil {
		modulos = []*Modulo{}
	}
	return &Trilha{
		BaseEntity: kernel.NewBaseEntityWithID(id, createdAt, updatedAt),
		slug:       slug,
		titulo:     titulo,
		descricao:  descricao,
		capaURL:    capaURL,
		ordem:      ordem,
		publicada:  publicada,
		modulos:    modulos,
	}
}

func (t *Trilha) Slug() valueobject.Slug { return t.slug }
func (t *Trilha) Titulo() string         { return t.titulo }
func (t *Trilha) Descricao() string      { return t.descricao }
func (t *Trilha) CapaURL() string        { return t.capaURL }
func (t *Trilha) Ordem() int             { return t.ordem }
func (t *Trilha) Publicada() bool        { return t.publicada }
func (t *Trilha) Modulos() []*Modulo     { return t.modulos }

func (t *Trilha) AdicionarModulo(slug valueobject.Slug, titulo, descricao string) (*Modulo, error) {
	for _, m := range t.modulos {
		if m.slug.Equals(slug) {
			return nil, errors.ErrAlreadyExists("slug de módulo já existe nesta trilha", "Trilha.AdicionarModulo", nil)
		}
	}
	ordem := len(t.modulos)
	mod, err := NovoModulo(slug, titulo, descricao, ordem)
	if err != nil {
		return nil, err
	}
	t.modulos = append(t.modulos, mod)
	t.Touch()
	return mod, nil
}

func (t *Trilha) Publicar() error {
	if len(t.modulos) == 0 {
		return errors.ErrInvalidState("trilha sem módulos não pode ser publicada", "Trilha.Publicar", nil)
	}
	t.publicada = true
	t.Touch()
	return nil
}
