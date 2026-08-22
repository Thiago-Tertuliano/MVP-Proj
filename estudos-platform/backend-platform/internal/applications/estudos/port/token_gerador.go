package port

import "time"

// Claims são os dados embutidos no token JWT.
type Claims struct {
	UsuarioID string
	Email     string
}

type TokenPar struct {
	AccessToken  string
	RefreshToken string
	AccessExp    time.Time
	RefreshExp   time.Time
}

// TokenGerador abstrai a emissão/validação de JWT para testar use cases com mock.
type TokenGerador interface {
	Gerar(claims Claims, accessTTL, refreshTTL time.Duration) (*TokenPar, error)
	ValidarAccessToken(token string) (*Claims, error)
}
