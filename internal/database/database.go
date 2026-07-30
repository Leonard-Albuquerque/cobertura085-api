package database

import (
	"fmt"
	"log"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa e retorna a conexão GORM com o banco de dados PostgreSQL.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados PostgreSQL: %w", err)
	}

	log.Println("Conexão com PostgreSQL estabelecida com sucesso.")
	return db, nil
}
