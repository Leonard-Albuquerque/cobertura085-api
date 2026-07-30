package main

import (
	"testing"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/handler"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
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

	storeService := service.NewStoreService(storeRepo, neighborhoodRepo, lobRepo)
	neighborhoodService := service.NewNeighborhoodService(neighborhoodRepo)
	commonService := service.NewCommonService(lobRepo, baseNeighborhoodRepo)
	telemetryService := service.NewTelemetryService(searchEventRepo)

	storeHandler := handler.NewStoreHandler(storeService)
	neighborhoodHandler := handler.NewNeighborhoodHandler(neighborhoodService)
	commonHandler := handler.NewCommonHandler(commonService)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService)

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

		commonGroup := apiV1.Group("")
		{
			commonGroup.GET("/lines-of-business", commonHandler.ListLinesOfBusiness)
			commonGroup.GET("/base-neighborhoods/by-name/:name", commonHandler.GetBaseNeighborhoodByName)
		}

		telemetryGroup := apiV1.Group("/telemetry")
		{
			telemetryGroup.POST("/search-events", telemetryHandler.CreateSearchEvent)
			telemetryGroup.GET("/summary", telemetryHandler.GetSummary)
			telemetryGroup.GET("/logs", telemetryHandler.GetRecentLogs)
			telemetryGroup.GET("/stores/:storeId", telemetryHandler.GetBusinessStats)
		}
	}

	t.Log("Gin radix tree router initialized successfully without panic!")
}
