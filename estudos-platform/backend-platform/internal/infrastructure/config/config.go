package config

import (
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config centraliza todas as variáveis de ambiente da aplicação.
type Config struct {
	AppEnv  string
	AppPort string

	// DatabaseURL (DATABASE_URL) tem prioridade sobre os DB_* — formato dos Postgres gerenciados.
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	// TrustProxy usa o 1º IP de X-Forwarded-For como cliente (API atrás de proxy/CDN).
	TrustProxy bool

	JWTSecret          string
	JWTAccessTTLMin    int
	JWTRefreshTTLHours int

	CORSAllowedOrigins []string

	// EditorEmails lista (em minúsculas) quem pode criar/editar roadmaps. Promove no primeiro acesso.
	EditorEmails []string
}

func Load() *Config {
	// .env é opcional (em produção as variáveis vêm do ambiente do host)
	_ = godotenv.Load()

	cfg := LoadDB()
	cfg.AppEnv = getEnv("APP_ENV", "development")
	// PaaS (Render, Railway...) injetam PORT; APP_PORT continua valendo localmente.
	cfg.AppPort = getEnv("APP_PORT", getEnv("PORT", "8080"))
	cfg.TrustProxy = getEnv("TRUST_PROXY", "") == "true"
	cfg.JWTSecret = getEnv("JWT_SECRET", "")
	cfg.JWTAccessTTLMin = getEnvInt("JWT_ACCESS_TTL_MIN", 15)
	cfg.JWTRefreshTTLHours = getEnvInt("JWT_REFRESH_TTL_HOURS", 168)
	cfg.CORSAllowedOrigins = splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"))
	for _, e := range splitCSV(getEnv("EDITOR_EMAILS", "")) {
		cfg.EditorEmails = append(cfg.EditorEmails, strings.ToLower(e))
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET é obrigatório")
	}
	if cfg.AppEnv == "production" && len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET precisa de pelo menos 32 caracteres em produção")
	}
	return cfg
}

// LoadDB lê só o Postgres. CLIs (migrate, content-job) não exigem JWT_SECRET.
func LoadDB() *Config {
	_ = godotenv.Load()
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5433"),
		DBUser:      getEnv("DB_USER", "estudos"),
		DBPassword:  getEnv("DB_PASSWORD", "estudos_dev"),
		DBName:      getEnv("DB_NAME", "estudos_platform"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),
	}
}

// PostgresURL devolve a URL de conexão (DATABASE_URL ou montada dos DB_*).
// `channel_binding` (padrão das strings do Neon) é removido: o driver do golang-migrate não o aceita.
func (c *Config) PostgresURL() string {
	if c.DatabaseURL != "" {
		u, err := url.Parse(c.DatabaseURL)
		if err != nil {
			return c.DatabaseURL
		}
		q := u.Query()
		q.Del("channel_binding")
		u.RawQuery = q.Encode()
		return u.String()
	}
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DBUser, c.DBPassword),
		Host:     net.JoinHostPort(c.DBHost, c.DBPort),
		Path:     "/" + c.DBName,
		RawQuery: url.Values{"sslmode": {c.DBSSLMode}}.Encode(),
	}
	return u.String()
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
