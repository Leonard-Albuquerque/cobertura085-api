package dto

import "time"

// BaseNeighborhoodResponse representa os dados do bairro base aninhado no bairro da loja
type BaseNeighborhoodResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OfficialName string `json:"officialName"`
}

// NeighborhoodResponse representa o retorno completo do bairro com baseNeighborhood incluso
type NeighborhoodResponse struct {
	ID                    string                    `json:"id"`
	StoreID               string                    `json:"storeId"`
	BaseNeighborhoodID    string                    `json:"baseNeighborhoodId"`
	DeliveryEnabled       bool                      `json:"deliveryEnabled"`
	Fee                   float64                   `json:"fee"`
	DeliveryTime          *string                   `json:"deliveryTime,omitempty"`
	MinimumOrder          *float64                  `json:"minimumOrder,omitempty"`
	FreeDeliveryThreshold *float64                  `json:"freeDeliveryThreshold,omitempty"`
	Notes                 *string                   `json:"notes,omitempty"`
	CreatedAt             time.Time                 `json:"createdAt"`
	UpdatedAt             time.Time                 `json:"updatedAt"`
	BaseNeighborhood      *BaseNeighborhoodResponse `json:"baseNeighborhood,omitempty"`
}

// UpdateNeighborhoodInput representa o payload de atualização de um único bairro
type UpdateNeighborhoodInput struct {
	DeliveryEnabled       *bool    `json:"deliveryEnabled,omitempty"`
	Fee                   *float64 `json:"fee,omitempty"`
	DeliveryTime          *string  `json:"deliveryTime,omitempty"`
	MinimumOrder          *float64 `json:"minimumOrder,omitempty"`
	FreeDeliveryThreshold *float64 `json:"freeDeliveryThreshold,omitempty"`
	Notes                 *string  `json:"notes,omitempty"`
}

// BulkUpdateNeighborhoodsInput representa o payload para atualização em lote de bairros
type BulkUpdateNeighborhoodsInput struct {
	IDs  []string                `json:"ids" binding:"required"`
	Data UpdateNeighborhoodInput `json:"data" binding:"required"`
}
