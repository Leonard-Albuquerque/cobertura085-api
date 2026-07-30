package model

import "time"

// PickupPoint representa a tabela "PickupPoint" no banco de dados.
type PickupPoint struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	StoreID      string    `gorm:"column:storeId;not null" json:"storeId"`
	Name         *string   `gorm:"column:name" json:"name,omitempty"`
	Address      string    `gorm:"column:address;not null" json:"address"`
	Latitude     *float64  `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude    *float64  `gorm:"column:longitude" json:"longitude,omitempty"`
	Instructions *string   `gorm:"column:instructions" json:"instructions,omitempty"`
	CreatedAt    time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;not null" json:"updatedAt"`

	// Relacionamento GORM
	Store *Store `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (PickupPoint) TableName() string {
	return "PickupPoint"
}
