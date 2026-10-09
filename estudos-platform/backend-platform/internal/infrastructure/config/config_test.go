package config

import "testing"

func TestPostgresURL_MontaDosCamposComEscape(t *testing.T) {
	c := &Config{DBHost: "localhost", DBPort: "5433", DBUser: "estudos", DBPassword: "p@ss:w/rd", DBName: "estudos_platform", DBSSLMode: "disable"}
	want := "postgres://estudos:p%40ss%3Aw%2Frd@localhost:5433/estudos_platform?sslmode=disable"
	if got := c.PostgresURL(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPostgresURL_DatabaseURLTemPrioridadeERemoveChannelBinding(t *testing.T) {
	c := &Config{
		DatabaseURL: "postgresql://u:p@ep-x.sa-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require",
		DBHost:      "ignorado",
	}
	want := "postgresql://u:p@ep-x.sa-east-1.aws.neon.tech/neondb?sslmode=require"
	if got := c.PostgresURL(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLoad_PortaDoPaaSEProxy(t *testing.T) {
	t.Setenv("JWT_SECRET", "segredo-de-teste")
	t.Setenv("PORT", "10000")
	t.Setenv("TRUST_PROXY", "true")
	cfg := Load()
	if cfg.AppPort != "10000" {
		t.Fatalf("sem APP_PORT deve usar PORT, got %q", cfg.AppPort)
	}
	if !cfg.TrustProxy {
		t.Fatal("TRUST_PROXY=true deve ligar TrustProxy")
	}

	t.Setenv("APP_PORT", "8080")
	if got := Load().AppPort; got != "8080" {
		t.Fatalf("APP_PORT tem prioridade, got %q", got)
	}
}
