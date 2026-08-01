package handler

import (
	"net/http"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type GeoJSONHandler struct {
	geoJSONService service.GeoJSONService
}

func NewGeoJSONHandler(geoJSONService service.GeoJSONService) *GeoJSONHandler {
	return &GeoJSONHandler{geoJSONService: geoJSONService}
}

// GetBairrosFortaleza serve o arquivo GeoJSON com os limites territoriais dos bairros de Fortaleza
func (h *GeoJSONHandler) GetBairrosFortaleza(c *gin.Context) {
	data, err := h.geoJSONService.GetBairrosFortaleza(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao carregar dados geográficos dos bairros")
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}
