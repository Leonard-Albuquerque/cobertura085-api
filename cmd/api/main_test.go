package main

import (
	"testing"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/handler"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service/external"
	"github.com/gin-gonic/gin"
)

func TestRouterInitialization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	storeRepo := repository.NewStoreRepository(nil)
	neighborhoodRepo := repository.NewNeighborhoodRepository(nil)
	baseNeighborhoodRepo := repository.NewBaseNeighborhoodRepository(nil)
	lobRepo := repository.NewLineOfBusinessRepository(nil)
	searchEventRepo := repository.NewSearchEventRepository(nil)

	viaCEPClient := external.NewViaCEPClient()
	nominatimClient := external.NewNominatimClient()

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

	storeHandler := handler.NewStoreHandler(storeService)
	neighborhoodHandler := handler.NewNeighborhoodHandler(neighborhoodService)
	commonHandler := handler.NewCommonHandler(commonService)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService)
	shippingHandler := handler.NewShippingHandler(shippingService)
	geoJSONHandler := handler.NewGeoJSONHandler(geoJSONService)

	apiV1 := router.Group("/api/v1")
	{
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

		neighborhoodsGroup := apiV1.Group("/neighborhoods")
		{
			neighborhoodsGroup.PATCH("/bulk", neighborhoodHandler.UpdateBulk)
			neighborhoodsGroup.PATCH("/:id", neighborhoodHandler.UpdateSingle)
		}

		shippingGroup := apiV1.Group("/shipping")
		{
			shippingGroup.POST("/lookup-cep", shippingHandler.LookupCEP)
			shippingGroup.POST("/lookup-address", shippingHandler.LookupAddress)
			shippingGroup.POST("/lookup-coords", shippingHandler.LookupCoords)
			shippingGroup.POST("/lookup-selected-address", shippingHandler.LookupSelectedAddress)
			shippingGroup.GET("/address-suggestions", shippingHandler.GetAddressSuggestions)
		}

		commonGroup := apiV1.Group("")
		{
			commonGroup.GET("/lines-of-business", commonHandler.ListLinesOfBusiness)
			commonGroup.GET("/base-neighborhoods/by-name/:name", commonHandler.GetBaseNeighborhoodByName)
			commonGroup.GET("/geojson/bairros-fortaleza", geoJSONHandler.GetBairrosFortaleza)
		}

		telemetryGroup := apiV1.Group("/telemetry")
		{
			telemetryGroup.POST("/search-events", telemetryHandler.CreateSearchEvent)
			telemetryGroup.GET("/summary", telemetryHandler.GetSummary)
			telemetryGroup.GET("/logs", telemetryHandler.GetRecentLogs)
			telemetryGroup.GET("/stores/:storeId", telemetryHandler.GetBusinessStats)
		}
	}

	t.Log("Roteador Gin inicializado com sucesso com todas as rotas da Fase 1 e Fase 2!")
}
