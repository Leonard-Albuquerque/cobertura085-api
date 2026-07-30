package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config armazena as configurações da aplicação carregadas de variáveis de ambiente.
type Config struct {
	Port                 string
	GinMode              string
	DatabaseURL          string
	JWTSecret            string
	JWTAccessExpiration  string
	JWTRefreshExpiration string
}

// Load carrega as variáveis de ambiente com valores padrão de fallback.
func Load() *Config {
	// 🔴 CORREÇÃO: Tenta carregar o .env na raiz ou subpastas
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("[INFO] Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
		}
	}

	cfg := &Config{
		Port:                 getEnv("PORT", "8080"),
		GinMode:              getEnv("GIN_MODE", "debug"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cobertura085?sslmode=disable"),
		JWTSecret:            getEnv("JWT_SECRET", "cobertura085_super_secret_jwt_key_change_in_prod"),
		JWTAccessExpiration:  getEnv("JWT_ACCESS_EXPIRATION", "15m"),
		JWTRefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", "168h"),
	}

	log.Printf("[DEBUG] DSN Carregada: %s\n", cfg.DSN())
	return cfg
}

// DSN retorna a string de conexão (Data Source Name) para o PostgreSQL.
func (c *Config) DSN() string {
	fmt.Println("DSN de Conexão:", c.DatabaseURL)
	return c.DatabaseURL
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
