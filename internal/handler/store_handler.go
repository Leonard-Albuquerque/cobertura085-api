package handler

import (
	"errors"
	"net/http"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService service.StoreService
}

func NewStoreHandler(storeService service.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

// ListPublic retorna a lista de todas as lojas para a busca pública na home
func (h *StoreHandler) ListPublic(c *gin.Context) {
	stores, err := h.storeService.GetPublicStores(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Falha ao carregar lojas")
		return
	}
	c.JSON(http.StatusOK, stores)
}

// GetBySlug retorna os dados completos da loja pelo slug ou ID
func (h *StoreHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("id")
	if slug == "" {
		slug = c.Param("slug")
	}
	if slug == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Slug ou ID da loja é obrigatório")
		return
	}

	store, pickupPoints, err := h.storeService.GetStoreBySlug(c.Request.Context(), slug)
	if errors.Is(err, repository.ErrStoreNotFound) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Loja não encontrada")
		return
	}
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao buscar loja")
		return
	}

	// Monta o payload incluindo pickupPoints no JSON
	response := map[string]interface{}{
		"id":                     store.ID,
		"slug":                   store.Slug,
		"name":                   store.Name,
		"whatsapp":               store.Whatsapp,
		"address":                store.Address,
		"operatingHours":         store.OperatingHours,
		"operatingHoursJson":     store.OperatingHoursJSON,
		"pickupEnabled":          store.PickupEnabled,
		"logoUrl":                store.LogoURL,
		"bannerUrl":              store.BannerURL,
		"description":            store.Description,
		"instagram":              store.Instagram,
		"catalogUrl":             store.CatalogURL,
		"websiteUrl":             store.WebsiteURL,
		"deliveryTimeDefault":    store.DeliveryTimeDefault,
		"deliveryAvailableMsg":   store.DeliveryAvailableMsg,
		"deliveryUnavailableMsg": store.DeliveryUnavailableMsg,
		"sameDayCutoff":          store.SameDayCutoff,
		"cutoffMessage":          store.CutoffMessage,
		"lineOfBusiness":         store.LineOfBusiness,
		"customLineOfBusiness":   store.CustomLineOfBusiness,
		"qrToken":                store.QRToken,
		"pickupPoints":           pickupPoints,
	}

	c.JSON(http.StatusOK, response)
}

// GetByQRToken resolve uma loja pelo seu QR Token
func (h *StoreHandler) GetByQRToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Token de QR Code é obrigatório")
		return
	}

	store, err := h.storeService.GetStoreByQRToken(c.Request.Context(), token)
	if errors.Is(err, repository.ErrStoreNotFound) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "QR Token não associado a nenhuma loja")
		return
	}
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao resolver QR Token")
		return
	}

	c.JSON(http.StatusOK, store)
}

// UpdateSettings atualiza as configurações da loja e recria seus pontos de retirada
func (h *StoreHandler) UpdateSettings(c *gin.Context) {
	storeID := c.Param("id")
	if storeID == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "ID da loja é obrigatório")
		return
	}

	var input dto.UpdateStoreSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "INVALID_INPUT", "Dados de entrada inválidos: "+err.Error())
		return
	}

	err := h.storeService.UpdateStoreSettings(c.Request.Context(), storeID, input)
	if errors.Is(err, service.ErrInvalidLineOfBusiness) {
		RespondError(c, http.StatusBadRequest, "INVALID_LOB", "Ramo de atuação fornecido é inválido")
		return
	}
	if errors.Is(err, repository.ErrStoreNotFound) {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Loja não encontrada")
		return
	}
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao atualizar configurações da loja")
		return
	}

	RespondSuccess(c, http.StatusOK, map[string]string{"message": "Configurações da loja atualizadas com sucesso"})
}

// GetDashboardStats retorna estatísticas da loja para o dashboard admin
func (h *StoreHandler) GetDashboardStats(c *gin.Context) {
	storeID := c.Param("id")
	if storeID == "" {
		RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "ID da loja é obrigatório")
		return
	}

	stats, err := h.storeService.GetDashboardStats(c.Request.Context(), storeID)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Erro ao buscar métricas do dashboard")
		return
	}

	c.JSON(http.StatusOK, stats)
}
