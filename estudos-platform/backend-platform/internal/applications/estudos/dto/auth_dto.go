package dto

// DTOs são contratos de borda: o que entra e o que sai da API.

type RegistrarRequest struct {
	Nome  string `json:"nome" validate:"required,min=2,max=150"`
	Email string `json:"email" validate:"required,email,max=255"`
	Senha string `json:"senha" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiracaoEm  int64  `json:"expiracao_em"` // unix seconds
}

type UsuarioResponse struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Tokens  TokenResponse   `json:"tokens"`
	Usuario UsuarioResponse `json:"usuario"`
}
