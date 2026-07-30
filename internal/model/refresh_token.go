package model

import (
	"time"
)

// RefreshToken representa a tabela "RefreshToken" no banco de dados.
type RefreshToken struct {
	ID        string     `gorm:"primaryKey;column:id" json:"id"`
	UserID    string     `gorm:"column:userId;not null;index" json:"userId"`
	TokenHash string     `gorm:"column:tokenHash;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expiresAt;not null" json:"expiresAt"`
	RevokedAt *time.Time `gorm:"column:revokedAt" json:"revokedAt,omitempty"`
	UserAgent *string    `gorm:"column:userAgent" json:"userAgent,omitempty"`
	IPHash    *string    `gorm:"column:ipHash" json:"ipHash,omitempty"`
	CreatedAt time.Time  `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`

	// Relacionamento com o Usuário
	User *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (RefreshToken) TableName() string {
	return "RefreshToken"
}
