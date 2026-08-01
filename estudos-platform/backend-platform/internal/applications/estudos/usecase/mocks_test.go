package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/entity"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/repository"
	"github.com/thiago-tertuliano/estudos-platform/internal/domain/estudos/valueobject"
)

// ---- MockUsuarioRepository ----
type MockUsuarioRepository struct {
	SaveFn        func(ctx context.Context, u *entity.Usuario) error
	FindByEmailFn func(ctx context.Context, email valueobject.Email) (*entity.Usuario, error)
	FindByIDFn    func(ctx context.Context, id string) (*entity.Usuario, error)
	EmailExisteFn func(ctx context.Context, email valueobject.Email) (bool, error)
}

func (m *MockUsuarioRepository) Save(ctx context.Context, u *entity.Usuario) error {
	return m.SaveFn(ctx, u)
}
func (m *MockUsuarioRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Usuario, error) {
	return m.FindByEmailFn(ctx, email)
}
func (m *MockUsuarioRepository) FindByID(ctx context.Context, id string) (*entity.Usuario, error) {
	return m.FindByIDFn(ctx, id)
}
func (m *MockUsuarioRepository) EmailExiste(ctx context.Context, email valueobject.Email) (bool, error) {
	return m.EmailExisteFn(ctx, email)
}

// ---- MockRefreshTokenRepository ----
type MockRefreshTokenRepository struct {
	SaveFn            func(ctx context.Context, t *repository.RefreshToken) error
	FindByHashFn      func(ctx context.Context, hash string) (*repository.RefreshToken, error)
	RevokeFn          func(ctx context.Context, id string) error
	RevokeAllByUserFn func(ctx context.Context, usuarioID uuid.UUID) error
}

func (m *MockRefreshTokenRepository) Save(ctx context.Context, t *repository.RefreshToken) error {
	return m.SaveFn(ctx, t)
}
func (m *MockRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*repository.RefreshToken, error) {
	return m.FindByHashFn(ctx, hash)
}
func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	return m.RevokeFn(ctx, id)
}
func (m *MockRefreshTokenRepository) RevokeAllByUser(ctx context.Context, usuarioID uuid.UUID) error {
	return m.RevokeAllByUserFn(ctx, usuarioID)
}

// ---- MockSenhaHasher ----
type MockSenhaHasher struct {
	HashFn     func(plain string) (string, error)
	CompararFn func(hash, plain string) bool
}

func (m *MockSenhaHasher) Hash(plain string) (string, error) { return m.HashFn(plain) }
func (m *MockSenhaHasher) Comparar(hash, plain string) bool  { return m.CompararFn(hash, plain) }

// ---- MockTokenGerador ----
type MockTokenGerador struct {
	GerarFn              func(claims port.Claims, accessTTL, refreshTTL time.Duration) (*port.TokenPar, error)
	ValidarAccessTokenFn func(token string) (*port.Claims, error)
}

func (m *MockTokenGerador) Gerar(c port.Claims, accessTTL, refreshTTL time.Duration) (*port.TokenPar, error) {
	return m.GerarFn(c, accessTTL, refreshTTL)
}
func (m *MockTokenGerador) ValidarAccessToken(token string) (*port.Claims, error) {
	return m.ValidarAccessTokenFn(token)
}
