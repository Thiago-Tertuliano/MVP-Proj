package external

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/shared/errors"
)

// JWTService implementa port.TokenGerador usando JWT HS256.
type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

// claims locais: embed do payload custom + claims registradas (sub, iat, exp)
type claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (s *JWTService) Gerar(c port.Claims, accessTTL, refreshTTL time.Duration) (*port.TokenPar, error) {
	now := time.Now()

	// access token: carrega email + subject (id)
	accessExp := now.Add(accessTTL)
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Email: c.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   c.UsuarioID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}).SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	// refresh token: apenas subject (sem claims sensíveis)
	refreshExp := now.Add(refreshTTL)
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   c.UsuarioID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(refreshExp),
	}).SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &port.TokenPar{
		AccessToken:  access,
		RefreshToken: refresh,
		AccessExp:    accessExp,
		RefreshExp:   refreshExp,
	}, nil
}

func (s *JWTService) ValidarAccessToken(token string) (*port.Claims, error) {
	var c claims
	parsed, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !parsed.Valid {
		return nil, errors.ErrUnauthorized("token inválido", "JWTService.ValidarAccessToken", nil)
	}
	return &port.Claims{UsuarioID: c.Subject, Email: c.Email}, nil
}
