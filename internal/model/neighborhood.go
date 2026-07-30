package model

import "time"

// Neighborhood representa a tabela "Neighborhood" no banco de dados.
type Neighborhood struct {
	ID                    string            `gorm:"primaryKey;column:id" json:"id"`
	StoreID               string            `gorm:"column:storeId;not null;uniqueIndex:idx_store_base_neighborhood" json:"storeId"`
	BaseNeighborhoodID    string            `gorm:"column:baseNeighborhoodId;not null;uniqueIndex:idx_store_base_neighborhood" json:"baseNeighborhoodId"`
	DeliveryEnabled       bool              `gorm:"column:deliveryEnabled;not null;default:false" json:"deliveryEnabled"`
	Fee                   float64           `gorm:"column:fee;not null;default:0.00" json:"fee"`
	DeliveryTime          *string           `gorm:"column:deliveryTime;default:'24h'" json:"deliveryTime,omitempty"`
	MinimumOrder          *float64          `gorm:"column:minimumOrder" json:"minimumOrder,omitempty"`
	FreeDeliveryThreshold *float64          `gorm:"column:freeDeliveryThreshold" json:"freeDeliveryThreshold,omitempty"`
	Notes                 *string           `gorm:"column:notes" json:"notes,omitempty"`
	CreatedAt             time.Time         `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt             time.Time         `gorm:"column:updatedAt;not null" json:"updatedAt"`

	// Relacionamentos GORM
	Store            *Store            `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	BaseNeighborhood *BaseNeighborhood `gorm:"foreignKey:BaseNeighborhoodID;references:ID" json:"baseNeighborhood,omitempty"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (Neighborhood) TableName() string {
	return "Neighborhood"
}
