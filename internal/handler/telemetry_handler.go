package handler

import (
	"net/http"
	"strconv"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type TelemetryHandler struct {
	service service.TelemetryService
}

func NewTelemetryHandler(service service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{service: service}
}

// CreateSearchEvent grava um novo log de telemetria de busca
func (h *TelemetryHandler) CreateSearchEvent(c *gin.Context) {
	var input dto.CreateSearchEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Payload de log inválido: "+err.Error())
		return
	}

	clientIP := c.ClientIP()
	err := h.service.LogSearchEvent(c.Request.Context(), input, clientIP)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao gravar log de telemetria")
		return
	}

	RespondSuccess(c, http.StatusCreated, map[string]string{"message": "Log registrado com sucesso"})
}

// GetSummary retorna o dashboard global de telemetria
func (h *TelemetryHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao carregar resumo de telemetria")
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetRecentLogs retorna o histórico recente de logs
func (h *TelemetryHandler) GetRecentLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	logs, err := h.service.GetRecentLogs(c.Request.Context(), limit)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao buscar histórico de logs")
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetBusinessStats retorna estatísticas de telemetria filtradas por empresa
func (h *TelemetryHandler) GetBusinessStats(c *gin.Context) {
	storeID := c.Param("storeId")
	if storeID == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "ID da loja é obrigatório")
		return
	}

	stats, err := h.service.GetBusinessStats(c.Request.Context(), storeID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao carregar telemetria da empresa")
		return
	}

	c.JSON(http.StatusOK, stats)
}
