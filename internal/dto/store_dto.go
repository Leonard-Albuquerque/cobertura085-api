package dto

import "time"

// PublicStoreResponse representa a DTO leve de loja exibida na busca da página inicial
type PublicStoreResponse struct {
	ID                string  `json:"id"`
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	LogoURL           *string `json:"logoUrl,omitempty"`
	Address           string  `json:"address"`
	OperatingHours    string  `json:"operatingHours"`
	PickupEnabled     bool    `json:"pickupEnabled"`
	HasDelivery       bool    `json:"hasDelivery"`
	PickupPointsCount int64   `json:"pickupPointsCount"`
}

// PickupPointInput representa os dados de um ponto de retirada recebidos no formulário
type PickupPointInput struct {
	Name         *string  `json:"name,omitempty"`
	Address      string   `json:"address"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	Instructions *string  `json:"instructions,omitempty"`
}

// UpdateStoreSettingsInput representa o payload de atualização de configurações da loja
type UpdateStoreSettingsInput struct {
	Name                   string             `json:"name" binding:"required"`
	Whatsapp               string             `json:"whatsapp" binding:"required"`
	Address                string             `json:"address" binding:"required"`
	PickupEnabled          bool               `json:"pickupEnabled"`
	LogoURL                *string            `json:"logoUrl,omitempty"`
	BannerURL              *string            `json:"bannerUrl,omitempty"`
	Description            *string            `json:"description,omitempty"`
	Instagram              *string            `json:"instagram,omitempty"`
	CatalogURL             *string            `json:"catalogUrl,omitempty"`
	WebsiteURL             *string            `json:"websiteUrl,omitempty"`
	DeliveryTimeDefault    string             `json:"deliveryTimeDefault"`
	DeliveryAvailableMsg   *string            `json:"deliveryAvailableMsg,omitempty"`
	DeliveryUnavailableMsg *string            `json:"deliveryUnavailableMsg,omitempty"`
	SameDayCutoff          *string            `json:"sameDayCutoff,omitempty"`
	CutoffMessage          *string            `json:"cutoffMessage,omitempty"`
	LineOfBusiness         *string            `json:"lineOfBusiness,omitempty"`
	CustomLineOfBusiness   *string            `json:"customLineOfBusiness,omitempty"`
	OperatingHoursJSON     any                `json:"operatingHoursJson,omitempty"`
	PickupPoints           []PickupPointInput `json:"pickupPoints"`
}

// StoreDashboardStatsResponse representa as métricas exibidas no painel da loja
type StoreDashboardStatsResponse struct {
	ActiveNeighborhoods   int64      `json:"activeNeighborhoods"`
	InactiveNeighborhoods int64      `json:"inactiveNeighborhoods"`
	AverageFee            float64    `json:"averageFee"`
	LastUpdatedAt         *time.Time `json:"lastUpdatedAt,omitempty"`
}
