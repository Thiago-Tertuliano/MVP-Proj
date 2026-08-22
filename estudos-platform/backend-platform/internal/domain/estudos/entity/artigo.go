package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/kernel"
)

// Artigo é o Aggregate Root de conteúdo.
type Artigo struct {
	kernel.BaseEntity
	slug        valueobject.Slug
	titulo      string
	subtitulo   string
	capaURL     string
	conteudo    json.RawMessage
	metadados   json.RawMessage
	autorID     uuid.UUID
	status      valueobject.ArtigoStatus
	publicadoEm *time.Time
	trilhaID    *uuid.UUID // Adcionado (Sprint A1)
	moduloID    *uuid.UUID // Adcionado (Sprint A1)
}

type NovoArtigoInput struct {
	Titulo    string
	Subtitulo string
	CapaURL   string
	Conteudo  json.RawMessage
	Metadados json.RawMessage
	Slug      valueobject.Slug
	AutorID   uuid.UUID
	TrilhaID  *uuid.UUID // Adcionado (Sprint A1)
	ModuloID  *uuid.UUID // Adcionado (Sprint A1)
}

func NovoArtigo(in NovoArtigoInput) (*Artigo, error) {
	if len(in.Titulo) < 3 || len(in.Titulo) > 300 {
		return nil, errors.ErrInvalidArgument("título deve ter entre 3 e 300 caracteres", "entity.NovoArtigo", nil)
	}
	if len(in.Conteudo) == 0 {
		in.Conteudo = json.RawMessage(`{}`)
	}
	if !json.Valid(in.Conteudo) {
		return nil, errors.ErrInvalidArgument("conteúdo JSON inválido", "entity.NovoArtigo", nil)
	}
	if len(in.Metadados) == 0 {
		in.Metadados = json.RawMessage(`{}`)
	}
	if !json.Valid(in.Metadados) {
		return nil, errors.ErrInvalidArgument("metadados JSON inválidos", "entity.NovoArtigo", nil)
	}
	if in.AutorID == uuid.Nil {
		return nil, errors.ErrInvalidArgument("autor_id obrigatório", "entity.NovoArtigo", nil)
	}

	if in.ModuloID != nil && in.TrilhaID == nil {
		return nil, errors.ErrInvalidArgument("módulo órfão sem trilha", "entity.NovoArtigo", nil)
	}

	return &Artigo{
		BaseEntity: kernel.NewBaseEntity(),
		slug:       in.Slug,
		titulo:     in.Titulo,
		subtitulo:  in.Subtitulo,
		capaURL:    in.CapaURL,
		conteudo:   in.Conteudo,
		metadados:  in.Metadados,
		autorID:    in.AutorID,
		status:     valueobject.ArtigoStatusRascunho,
		trilhaID:   in.TrilhaID,
		moduloID:   in.ModuloID,
	}, nil
}

func ReconstruirArtigo(
	id uuid.UUID,
	slug valueobject.Slug,
	titulo, subtitulo, capaURL string,
	conteudo, metadados json.RawMessage,
	autorID uuid.UUID,
	status valueobject.ArtigoStatus,
	publicadoEm *time.Time,
	createdAt, updatedAt time.Time,
	trilhaID, moduloID *uuid.UUID, // Adicionado (Sprint A1)
) *Artigo {
	return &Artigo{
		BaseEntity:  kernel.NewBaseEntityWithID(id, createdAt, updatedAt),
		slug:        slug,
		titulo:      titulo,
		subtitulo:   subtitulo,
		capaURL:     capaURL,
		conteudo:    conteudo,
		metadados:   metadados,
		autorID:     autorID,
		status:      status,
		publicadoEm: publicadoEm,
		trilhaID:    trilhaID, // Adicionado (Sprint A1)
		moduloID:    moduloID, // Adicionado (Sprint A1)
	}
}

func (a *Artigo) Slug() valueobject.Slug           { return a.slug }
func (a *Artigo) Titulo() string                   { return a.titulo }
func (a *Artigo) Subtitulo() string                { return a.subtitulo }
func (a *Artigo) CapaURL() string                  { return a.capaURL }
func (a *Artigo) Conteudo() json.RawMessage        { return a.conteudo }
func (a *Artigo) Metadados() json.RawMessage       { return a.metadados }
func (a *Artigo) AutorID() uuid.UUID               { return a.autorID }
func (a *Artigo) Status() valueobject.ArtigoStatus { return a.status }
func (a *Artigo) PublicadoEm() *time.Time          { return a.publicadoEm }
func (a *Artigo) TrilhaID() *uuid.UUID             { return a.trilhaID } // Adicionado (Sprint A1)
func (a *Artigo) ModuloID() *uuid.UUID             { return a.moduloID } // Adicionado (Sprint A1)

func (a *Artigo) AtualizarConteudo(titulo, subtitulo, capaURL string, conteudo, metadados json.RawMessage) error {
	if a.status == valueobject.ArtigoStatusArquivado {
		return errors.ErrInvalidState("artigo arquivado não pode ser editado", "Artigo.AtualizarConteudo", nil)
	}
	if len(titulo) < 3 || len(titulo) > 300 {
		return errors.ErrInvalidArgument("título deve ter entre 3 e 300 caracteres", "Artigo.AtualizarConteudo", nil)
	}
	if len(conteudo) == 0 {
		conteudo = json.RawMessage(`{}`)
	}
	if !json.Valid(conteudo) {
		return errors.ErrInvalidArgument("conteúdo JSON inválido", "Artigo.AtualizarConteudo", nil)
	}
	if len(metadados) == 0 {
		metadados = json.RawMessage(`{}`)
	}
	if !json.Valid(metadados) {
		return errors.ErrInvalidArgument("metadados JSON inválidos", "Artigo.AtualizarConteudo", nil)
	}
	a.titulo = titulo
	a.subtitulo = subtitulo
	a.capaURL = capaURL
	a.conteudo = conteudo
	a.metadados = metadados
	a.Touch()
	return nil
}

func (a *Artigo) VincularTrilhaEModulo(trilhaID, moduloID *uuid.UUID) error {
	if moduloID != nil && trilhaID == nil {
		return errors.ErrInvalidArgument("módulo órfão sem trilha", "Artigo.VincularTrilhaEModulo", nil)
	}
	a.trilhaID = trilhaID
	a.moduloID = moduloID
	a.Touch()
	return nil
}

func (a *Artigo) EnviarParaRevisao() error {
	if a.status != valueobject.ArtigoStatusRascunho {
		return errors.ErrInvalidState("só rascunho pode ir para revisão", "Artigo.EnviarParaRevisao", nil)
	}
	if len(a.conteudo) == 0 || string(a.conteudo) == "{}" {
		return errors.ErrInvalidArgument("conteúdo vazio", "Artigo.EnviarParaRevisao", nil)
	}
	a.status = valueobject.ArtigoStatusRevisao
	a.Touch()
	return nil
}

func (a *Artigo) Publicar() error {
	if a.status != valueobject.ArtigoStatusRevisao && a.status != valueobject.ArtigoStatusRascunho {
		return errors.ErrInvalidState("só rascunho ou revisão podem ser publicados", "Artigo.Publicar", nil)
	}
	if len(a.conteudo) == 0 || string(a.conteudo) == "{}" {
		return errors.ErrInvalidArgument("conteúdo vazio", "Artigo.Publicar", nil)
	}
	now := time.Now().UTC()
	a.status = valueobject.ArtigoStatusPublicado
	a.publicadoEm = &now
	a.Touch()
	return nil
}

func (a *Artigo) EhAutor(usuarioID uuid.UUID) bool {
	return a.autorID == usuarioID
}
