package handler

import (
	"net/http"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type ShippingHandler struct {
	shippingService service.ShippingService
}

func NewShippingHandler(shippingService service.ShippingService) *ShippingHandler {
	return &ShippingHandler{shippingService: shippingService}
}

// LookupCEP realiza a busca de frete a partir de um CEP
func (h *ShippingHandler) LookupCEP(c *gin.Context) {
	var input dto.LookupCEPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Payload inválido: "+err.Error())
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	res, err := h.shippingService.LookupCEP(c.Request.Context(), input, clientIP, userAgent)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao consultar CEP")
		return
	}

	c.JSON(http.StatusOK, res)
}

// LookupAddress realiza a busca de frete a partir do texto de endereço
func (h *ShippingHandler) LookupAddress(c *gin.Context) {
	var input dto.LookupAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Payload inválido: "+err.Error())
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	res, err := h.shippingService.LookupAddress(c.Request.Context(), input, clientIP, userAgent)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao consultar endereço")
		return
	}

	c.JSON(http.StatusOK, res)
}

// LookupCoords realiza a busca de frete a partir das coordenadas GPS (lat, lon)
func (h *ShippingHandler) LookupCoords(c *gin.Context) {
	var input dto.LookupCoordsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Payload inválido: "+err.Error())
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	res, err := h.shippingService.LookupCoords(c.Request.Context(), input, clientIP, userAgent)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao consultar coordenadas")
		return
	}

	c.JSON(http.StatusOK, res)
}

// LookupSelectedAddress realiza a consulta a partir de um item selecionado do autocompletar
func (h *ShippingHandler) LookupSelectedAddress(c *gin.Context) {
	var input dto.LookupSelectedAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Payload inválido: "+err.Error())
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	res, err := h.shippingService.LookupSelectedAddress(c.Request.Context(), input, clientIP, userAgent)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao processar endereço selecionado")
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetAddressSuggestions retorna até 5 sugestões de endereço em Fortaleza
func (h *ShippingHandler) GetAddressSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, []dto.AddressSuggestionItem{})
		return
	}

	suggestions, err := h.shippingService.GetAddressSuggestions(c.Request.Context(), query)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao buscar sugestões")
		return
	}

	c.JSON(http.StatusOK, suggestions)
}
