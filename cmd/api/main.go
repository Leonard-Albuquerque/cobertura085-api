package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/config"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/database"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 1. Carregar Configurações de Variáveis de Ambiente
	cfg := config.Load()

	// 2. Configurar modo do Gin
	gin.SetMode(cfg.GinMode)

	// 3. Inicializar Conexão com o Banco de Dados (GORM + PostgreSQL)
	var db *gorm.DB
	var err error
	db, err = database.InitDB(cfg)
	if err != nil {
		log.Printf("[AVISO] Não foi possível conectar ao PostgreSQL durante a inicialização: %v\n", err)
	}

	// 4. Injeção de Dependência Manual (Construtores)
	healthHandler := handler.NewHealthHandler(db)

	// 5. Configurar Roteador Gin
	router := gin.Default()

	// Registro de Rotas Base
	router.GET("/health", healthHandler.Check)

	// 6. Configurar Servidor HTTP
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Iniciar Servidor em uma goroutine separada
	go func() {
		log.Printf("Servidor API iniciado na porta %s (GIN_MODE=%s)...\n", cfg.Port, cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro crítico no servidor HTTP: %v\n", err)
		}
	}()

	// 7. Configuração do Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Aguarda sinal de interrupção (SIGINT / SIGTERM)
	sig := <-quit
	log.Printf("Sinal OS recebido: %v. Encerrando servidor graciosamente...\n", sig)

	// Contexto com limite de tempo de 5 segundos para encerrar requisições ativas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Servidor forçado a encerrar: %v\n", err)
	}

	log.Println("Servidor finalizado com sucesso.")
}
