package external

import (
	"testing"
	"time"

	"github.com/thiago-tertuliano/estudos-platform/internal/applications/estudos/port"
)

func TestJWTService_GerarEValidar(t *testing.T) {
	svc := NewJWTService("segredo-teste")
	par, err := svc.Gerar(port.Claims{UsuarioID: "u-1", Email: "t@ex.com"}, 15*time.Minute, 168*time.Hour)
	if err != nil {
		t.Fatalf("erro ao gerar: %v", err)
	}
	if par.AccessToken == "" || par.RefreshToken == "" {
		t.Error("tokens não deveriam ser vazios")
	}

	claims, err := svc.ValidarAccessToken(par.AccessToken)
	if err != nil {
		t.Fatalf("token deveria ser válido: %v", err)
	}
	if claims.UsuarioID != "u-1" || claims.Email != "t@ex.com" {
		t.Errorf("claims incorretos: %+v", claims)
	}
}

func TestJWTService_AssinaturaErrada(t *testing.T) {
	svc := NewJWTService("segredo")
	par, _ := svc.Gerar(port.Claims{UsuarioID: "u-1"}, 15*time.Minute, time.Hour)

	svc2 := NewJWTService("outro-segredo")
	if _, err := svc2.ValidarAccessToken(par.AccessToken); err == nil {
		t.Error("token assinado com outra chave deveria falhar")
	}
}

func TestJWTService_TokenExpirado(t *testing.T) {
	svc := NewJWTService("segredo")
	par, _ := svc.Gerar(port.Claims{UsuarioID: "u-1"}, -time.Minute, -time.Hour) // TTL negativo = já expirou

	if _, err := svc.ValidarAccessToken(par.AccessToken); err == nil {
		t.Error("token expirado deveria falhar")
	}
}
