package dto

type CriarTrilhaRequest struct {
	Titulo    string `json:"titulo" validate:"required,min=3,max=200"`
	Descricao string `json:"descricao" validate:"omitempty,max=2000"`
	CapaURL   string `json:"capa_url" validate:"omitempty,url"`
	Slug      string `json:"slug" validate:"omitempty,min=2,max=200"`
	Ordem     int    `json:"ordem" validate:"omitempty,min=0"`
}

type AdicionarModuloRequest struct {
	Titulo    string `json:"titulo" validate:"required,min=2,max=200"`
	Descricao string `json:"descricao" validate:"omitempty,max=2000"`
	Slug      string `json:"slug" validate:"omitempty,min=2,max=200"`
}

type ModuloResponse struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao,omitempty"`
	Ordem     int    `json:"ordem"`
}

type TrilhaResponse struct {
	ID        string           `json:"id"`
	Slug      string           `json:"slug"`
	Titulo    string           `json:"titulo"`
	Descricao string           `json:"descricao,omitempty"`
	CapaURL   string           `json:"capa_url,omitempty"`
	Ordem     int              `json:"ordem"`
	Publicada bool             `json:"publicada"`
	Modulos   []ModuloResponse `json:"modulos"`
	CreatedAt int64            `json:"created_at"`
	UpdatedAt int64            `json:"updated_at"`
}

type ListarTrilhasResponse struct {
	Itens []*TrilhaResponse `json:"itens"`
}
