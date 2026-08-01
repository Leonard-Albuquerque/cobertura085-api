package dto

// LookupCEPInput representa a payload para consulta de frete por CEP
type LookupCEPInput struct {
	StoreSlug string `json:"storeSlug" binding:"required"`
	CEP       string `json:"cep" binding:"required"`
}

// LookupAddressInput representa a payload para consulta de frete por texto de endereço
type LookupAddressInput struct {
	StoreSlug string `json:"storeSlug" binding:"required"`
	Address   string `json:"address" binding:"required"`
}

// LookupCoordsInput representa a payload para consulta por coordenadas GPS
type LookupCoordsInput struct {
	StoreSlug string  `json:"storeSlug" binding:"required"`
	Lat       float64 `json:"lat" binding:"required"`
	Lon       float64 `json:"lon" binding:"required"`
}

// LookupSelectedAddressInput representa a payload para consulta com endereço pré-selecionado do autocompletar
type LookupSelectedAddressInput struct {
	StoreSlug  string  `json:"storeSlug" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	Lat        float64 `json:"lat" binding:"required"`
	Lon        float64 `json:"lon" binding:"required"`
	BairroName string  `json:"bairroName,omitempty"`
}

// LookupResultResponse representa a resposta unificada esperada pelo frontend
type LookupResultResponse struct {
	Success               bool     `json:"success"`
	Error                 string   `json:"error,omitempty"`
	Bairro                string   `json:"bairro,omitempty"`
	Street                string   `json:"street,omitempty"`
	DeliveryEnabled       bool     `json:"deliveryEnabled"`
	Fee                   *float64 `json:"fee,omitempty"`
	DeliveryTime          string   `json:"deliveryTime,omitempty"`
	MinimumOrder          *float64 `json:"minimumOrder,omitempty"`
	FreeDeliveryThreshold *float64 `json:"freeDeliveryThreshold,omitempty"`
	Notes                 *string  `json:"notes,omitempty"`
	StoreAddress          string   `json:"storeAddress,omitempty"`
	StoreWhatsapp         string   `json:"storeWhatsapp,omitempty"`
	PickupEnabled         bool     `json:"pickupEnabled"`
}

// AddressSuggestionItem representa um item de sugestão de endereço para o autocompletar
type AddressSuggestionItem struct {
	DisplayName string  `json:"display_name"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Bairro      string  `json:"bairro,omitempty"`
	Road        string  `json:"road,omitempty"`
}
