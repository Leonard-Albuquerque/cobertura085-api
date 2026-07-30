package handler

import (
	"errors"
	"net/http"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	service service.CommonService
}

func NewCommonHandler(service service.CommonService) *CommonHandler {
	return &CommonHandler{service: service}
}

// ListLinesOfBusiness lista todos os ramos de atuação ativos
func (h *CommonHandler) ListLinesOfBusiness(c *gin.Context) {
	lobs, err := h.service.GetLinesOfBusiness(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao buscar ramos de atuação")
		return
	}
	c.JSON(http.StatusOK, lobs)
}

// GetBaseNeighborhoodByName realiza busca normalizada do bairro base pelo nome
func (h *CommonHandler) GetBaseNeighborhoodByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Nome do bairro é obrigatório")
		return
	}

	bn, err := h.service.GetBaseNeighborhoodByName(c.Request.Context(), name)
	if errors.Is(err, repository.ErrBaseNeighborhoodNotFound) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Bairro base não encontrado")
		return
	}
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao buscar bairro base")
		return
	}

	c.JSON(http.StatusOK, bn)
}
