package dto

import "time"

// CreateSearchEventInput representa o payload para registrar um log de busca
type CreateSearchEventInput struct {
	StoreID               string  `json:"storeId" binding:"required"`
	EventType             string  `json:"eventType" binding:"required"`
	SearchType            string  `json:"searchType" binding:"required"`
	SearchedValue         string  `json:"searchedValue" binding:"required"`
	SearchedNeighborhood  *string `json:"searchedNeighborhood,omitempty"`
	MatchedNeighborhoodID *string `json:"matchedNeighborhoodId,omitempty"`
	DeliveryAvailable     bool    `json:"deliveryAvailable"`
	DeliveryPrice         float64 `json:"deliveryPrice"`
	ResponseTimeMs        int     `json:"responseTimeMs"`
	SessionID             string  `json:"sessionId" binding:"required"`
	IPHash                string  `json:"ipHash"`
	UserAgent             string  `json:"userAgent"`
}

// TopNeighborhoodItem representa a contagem de buscas por bairro
type TopNeighborhoodItem struct {
	Neighborhood string `json:"neighborhood"`
	Count        int64  `json:"count"`
}

// StoreSearchShareItem representa a distribuição de buscas por loja
type StoreSearchShareItem struct {
	StoreID   string `json:"storeId"`
	StoreName string `json:"storeName"`
	Count     int64  `json:"count"`
}

// TelemetrySummaryResponse representa o resumo global de métricas e analytics
type TelemetrySummaryResponse struct {
	TotalSearches         int64                  `json:"totalSearches"`
	UniqueVisitors        int64                  `json:"uniqueVisitors"`
	AvailableRate         float64                `json:"availableRate"`
	AverageResponseTimeMs float64                `json:"averageResponseTimeMs"`
	TopNeighborhoods      []TopNeighborhoodItem  `json:"topNeighborhoods"`
	StoreSearchShares     []StoreSearchShareItem `json:"storeSearchShares"`
}

// SearchEventResponse representa um registro individual de log no painel admin
type SearchEventResponse struct {
	ID                   string     `json:"id"`
	CreatedAt            time.Time  `json:"createdAt"`
	StoreID              string     `json:"storeId"`
	StoreName            *string    `json:"storeName,omitempty"`
	EventType            string     `json:"eventType"`
	SearchType           string     `json:"searchType"`
	SearchedValue        string     `json:"searchedValue"`
	SearchedNeighborhood *string    `json:"searchedNeighborhood,omitempty"`
	DeliveryAvailable    bool       `json:"deliveryAvailable"`
	DeliveryPrice        float64    `json:"deliveryPrice"`
	ResponseTimeMs       int        `json:"responseTimeMs"`
	SessionID            string     `json:"sessionId"`
	UserAgent            string     `json:"userAgent"`
}
