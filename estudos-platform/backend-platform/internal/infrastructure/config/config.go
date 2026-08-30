package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config centraliza todas as variáveis de ambiente da aplicação.
type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret          string
	JWTAccessTTLMin    int
	JWTRefreshTTLHours int

	CORSAllowedOrigins []string
}

func Load() *Config {
	// .env é opcional (em produção as variáveis vêm do ambiente do host)
	_ = godotenv.Load()

	cfg := LoadDB()
	cfg.AppEnv = getEnv("APP_ENV", "development")
	cfg.AppPort = getEnv("APP_PORT", "8080")
	cfg.JWTSecret = getEnv("JWT_SECRET", "")
	cfg.JWTAccessTTLMin = getEnvInt("JWT_ACCESS_TTL_MIN", 15)
	cfg.JWTRefreshTTLHours = getEnvInt("JWT_REFRESH_TTL_HOURS", 168)
	cfg.CORSAllowedOrigins = splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"))

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET é obrigatório")
	}
	return cfg
}

// LoadDB lê só o Postgres. CLIs (migrate, content-job) não exigem JWT_SECRET.
func LoadDB() *Config {
	_ = godotenv.Load()
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "estudos"),
		DBPassword: getEnv("DB_PASSWORD", "estudos_dev"),
		DBName:     getEnv("DB_NAME", "estudos_platform"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
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
