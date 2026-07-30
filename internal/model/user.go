package model

import (
	"time"
)

// User representa a tabela "User" no banco de dados.
type User struct {
	ID           string     `gorm:"primaryKey;column:id" json:"id"`
	StoreID      string     `gorm:"column:storeId;not null;uniqueIndex" json:"storeId"`
	Email        string     `gorm:"column:email;not null;uniqueIndex" json:"email"`
	PasswordHash string     `gorm:"column:passwordHash;not null" json:"-"`
	Name         *string    `gorm:"column:name" json:"name,omitempty"`
	IsActive     bool       `gorm:"column:isActive;not null;default:true" json:"isActive"`
	LastLoginAt  *time.Time `gorm:"column:lastLoginAt" json:"lastLoginAt,omitempty"`
	CreatedAt    time.Time  `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updatedAt;not null" json:"updatedAt"`

	// Relacionamento com a loja
	Store *Store `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (User) TableName() string {
	return "User"
}
