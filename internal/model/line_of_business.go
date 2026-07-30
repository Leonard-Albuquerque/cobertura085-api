package model

import "time"

// LineOfBusiness representa a tabela "LineOfBusiness" no banco de dados.
type LineOfBusiness struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	Code      string    `gorm:"column:code;not null;unique" json:"code"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	IsActive  bool      `gorm:"column:isActive;not null;default:true" json:"isActive"`
	SortOrder int       `gorm:"column:sortOrder;not null;default:0" json:"sortOrder"`
	CreatedAt time.Time `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;not null" json:"updatedAt"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (LineOfBusiness) TableName() string {
	return "LineOfBusiness"
}
