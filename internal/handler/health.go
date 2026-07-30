package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler lida com as requisições do endpoint de health check.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler construtor manual para HealthHandler (Injeção de Dependência).
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check verifica a saúde da API e a conectividade com o banco de dados.
func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "up"

	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}
	} else {
		dbStatus = "not_connected"
	}

	status := "healthy"
	if dbStatus == "down" {
		status = "degraded"
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"status":    status,
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
