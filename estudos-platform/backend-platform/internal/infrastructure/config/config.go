package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/tools/go/cfg"
)

//Config centraliza todas as variáveis de ambiente da aplicação.
type Config struct {
	AppEnv 				string
	AppPort 			string

	DBHost 				string
	DBPort 				string
	DBUser 				string
	DBPassword 			string
	DBName 				string
	DBSSLMode 			string

	JWTSecret 			string
	JWTAcessTTLMin 		int
	JWTRefreshTTLHours 	int
}

func Load() *Config {
	// .env é opcional (em produção as variáveis vêm do ambiente do host)
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv: 				getEnv("APP_ENV", "development"),
		AppPort: 				getEnv("APP_PORT", "8080"),

		DBHost: 				getEnv("DB_HOST", "localhost"),
		DBPort: 				getEnv("DB_PORT", "5432"),
		DBUser: 				getEnv("DB_USER", "estudos"),
		DBPassword: 			getEnv("DB_PASSWORD", "estudos_dev"),
		DBName: 				getEnv("DB_NAME", "estudos_platform"),
		DBSSLMode: 				getEnv("DB_SSL_MODE", "disable"),

		JWTSecret: 				getEnv("JWT_SECRET", ""),
		JWTAcessTTLMin: 		getEnv("JWT_ACESS_TTL_MIN", 15),
		JWTRefreshTTLHours: 	getEnv("JWT_REFRESH_TTL_HOURS", 168),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET é obrigatório.")
	}
	return cfg
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