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
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/middleware"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service/external"
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

	// 4. Clientes HTTP Externos (ViaCEP e OpenStreetMap Nominatim)
	viaCEPClient := external.NewViaCEPClient()
	nominatimClient := external.NewNominatimClient()

	// Repositórios
	userRepo := repository.NewUserRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	neighborhoodRepo := repository.NewNeighborhoodRepository(db)
	baseNeighborhoodRepo := repository.NewBaseNeighborhoodRepository(db)
	lobRepo := repository.NewLineOfBusinessRepository(db)
	searchEventRepo := repository.NewSearchEventRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Serviços
	jwtService := service.NewJWTService(cfg.JWTSecret)
	authService := service.NewAuthService(userRepo, storeRepo, refreshTokenRepo, jwtService, cfg)
	storeService := service.NewStoreService(storeRepo, neighborhoodRepo, lobRepo)
	neighborhoodService := service.NewNeighborhoodService(neighborhoodRepo)
	commonService := service.NewCommonService(lobRepo, baseNeighborhoodRepo)
	telemetryService := service.NewTelemetryService(searchEventRepo)
	shippingService := service.NewShippingService(
		viaCEPClient,
		nominatimClient,
		storeRepo,
		baseNeighborhoodRepo,
		neighborhoodRepo,
		telemetryService,
	)
	geoJSONService := service.NewGeoJSONService("")

	// Handlers HTTP
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authService)
	storeHandler := handler.NewStoreHandler(storeService)
	neighborhoodHandler := handler.NewNeighborhoodHandler(neighborhoodService)
	commonHandler := handler.NewCommonHandler(commonService)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService)
	shippingHandler := handler.NewShippingHandler(shippingService)
	geoJSONHandler := handler.NewGeoJSONHandler(geoJSONService)

	// 5. Configurar Roteador Gin
	router := gin.Default()

	// Registro de Rotas Base
	router.GET("/health", healthHandler.Check)

	// Grupo API v1
	apiV1 := router.Group("/api/v1")
	{
		// Rotas de Autenticação
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.POST("/logout", authHandler.Logout)

			// Rotas Protegidas de Autenticação
			authProtected := authGroup.Group("", middleware.AuthMiddleware(jwtService))
			{
				authProtected.GET("/me", authHandler.GetProfile)
			}
		}

		// Rotas do Módulo de Lojas
		storesGroup := apiV1.Group("/stores")
		{
			storesGroup.GET("", storeHandler.ListPublic)
			storesGroup.GET("/qr/:token", storeHandler.GetByQRToken)
			storesGroup.GET("/:id", storeHandler.GetBySlug)
			storesGroup.PUT("/:id/settings", storeHandler.UpdateSettings)
			storesGroup.GET("/:id/dashboard-stats", storeHandler.GetDashboardStats)
			storesGroup.GET("/:id/neighborhoods", neighborhoodHandler.ListByStore)
			storesGroup.GET("/:id/neighborhoods/check", neighborhoodHandler.CheckByStore)
		}

		// Rotas do Módulo de Bairros
		neighborhoodsGroup := apiV1.Group("/neighborhoods")
		{
			neighborhoodsGroup.PATCH("/bulk", neighborhoodHandler.UpdateBulk)
			neighborhoodsGroup.PATCH("/:id", neighborhoodHandler.UpdateSingle)
		}

		// Rotas de Frete e Consultas Externas
		shippingGroup := apiV1.Group("/shipping")
		{
			shippingGroup.POST("/lookup-cep", shippingHandler.LookupCEP)
			shippingGroup.POST("/lookup-address", shippingHandler.LookupAddress)
			shippingGroup.POST("/lookup-coords", shippingHandler.LookupCoords)
			shippingGroup.POST("/lookup-selected-address", shippingHandler.LookupSelectedAddress)
			shippingGroup.GET("/address-suggestions", shippingHandler.GetAddressSuggestions)
		}

		// Rotas do Módulo de Utilitários e Domínio
		commonGroup := apiV1.Group("")
		{
			commonGroup.GET("/lines-of-business", commonHandler.ListLinesOfBusiness)
			commonGroup.GET("/base-neighborhoods/by-name/:name", commonHandler.GetBaseNeighborhoodByName)
			commonGroup.GET("/geojson/bairros-fortaleza", geoJSONHandler.GetBairrosFortaleza)
		}

		// Rotas do Módulo de Telemetria e Analytics
		telemetryGroup := apiV1.Group("/telemetry")
		{
			telemetryGroup.POST("/search-events", telemetryHandler.CreateSearchEvent)
			telemetryGroup.GET("/summary", telemetryHandler.GetSummary)
			telemetryGroup.GET("/logs", telemetryHandler.GetRecentLogs)
			telemetryGroup.GET("/stores/:storeId", telemetryHandler.GetBusinessStats)
		}
	}

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
