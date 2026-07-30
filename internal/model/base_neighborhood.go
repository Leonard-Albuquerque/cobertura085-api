package model

import "time"

// BaseNeighborhood representa a tabela "BaseNeighborhood" no banco de dados.
type BaseNeighborhood struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	Name         string    `gorm:"column:name;not null;index" json:"name"`
	OfficialName string    `gorm:"column:officialName;not null" json:"officialName"`
	CreatedAt    time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (BaseNeighborhood) TableName() string {
	return "BaseNeighborhood"
}
