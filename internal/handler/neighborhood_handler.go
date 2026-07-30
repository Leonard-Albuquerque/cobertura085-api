package handler

import (
	"errors"
	"net/http"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type NeighborhoodHandler struct {
	service service.NeighborhoodService
}

func NewNeighborhoodHandler(service service.NeighborhoodService) *NeighborhoodHandler {
	return &NeighborhoodHandler{service: service}
}

// ListByStore lista a matriz de bairros e taxas de uma determinada loja
func (h *NeighborhoodHandler) ListByStore(c *gin.Context) {
	storeID := c.Param("id")
	if storeID == "" {
		storeID = c.Param("storeId")
	}
	if storeID == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "ID da loja é obrigatório")
		return
	}

	list, err := h.service.GetStoreNeighborhoods(c.Request.Context(), storeID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao carregar bairros da loja")
		return
	}

	c.JSON(http.StatusOK, list)
}

// CheckByStore realiza o lookup de taxa para um bairro base em uma loja
func (h *NeighborhoodHandler) CheckByStore(c *gin.Context) {
	storeID := c.Param("id")
	if storeID == "" {
		storeID = c.Param("storeId")
	}
	baseID := c.Query("baseNeighborhoodId")

	if storeID == "" || baseID == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "storeId e baseNeighborhoodId são obrigatórios")
		return
	}

	neighborhood, err := h.service.CheckStoreNeighborhood(c.Request.Context(), storeID, baseID)
	if errors.Is(err, repository.ErrNeighborhoodNotFound) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Regra de entrega não encontrada para este bairro")
		return
	}
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao verificar taxa de bairro")
		return
	}

	c.JSON(http.StatusOK, neighborhood)
}

// UpdateSingle atualiza as regras de entrega de um único bairro
func (h *NeighborhoodHandler) UpdateSingle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "ID do bairro é obrigatório")
		return
	}

	var input dto.UpdateNeighborhoodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Corpo da requisição inválido")
		return
	}

	err := h.service.UpdateNeighborhood(c.Request.Context(), id, input)
	if errors.Is(err, repository.ErrNeighborhoodNotFound) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Bairro não encontrado")
		return
	}
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao atualizar bairro")
		return
	}

	RespondSuccess(c, http.StatusOK, map[string]string{"message": "Bairro atualizado com sucesso"})
}

// UpdateBulk atualiza as regras de entrega em massa para múltiplos bairros
func (h *NeighborhoodHandler) UpdateBulk(c *gin.Context) {
	var input dto.BulkUpdateNeighborhoodsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Payload de lote inválido: "+err.Error())
		return
	}

	err := h.service.UpdateNeighborhoodsBulk(c.Request.Context(), input)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao atualizar bairros em lote")
		return
	}

	RespondSuccess(c, http.StatusOK, map[string]string{"message": "Bairros atualizados em lote com sucesso"})
}
