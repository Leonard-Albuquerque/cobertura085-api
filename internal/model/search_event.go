package model

import "time"

// SearchEvent representa a tabela "SearchEvent" no banco de dados.
type SearchEvent struct {
	ID                    string            `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt             time.Time         `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP;index" json:"createdAt"`
	StoreID               string            `gorm:"column:storeId;not null;index" json:"storeId"`
	EventType             string            `gorm:"column:eventType;type:varchar(30);not null" json:"eventType"`
	SearchType            string            `gorm:"column:searchType;type:varchar(20);not null" json:"searchType"`
	SearchedValue         string            `gorm:"column:searchedValue;type:varchar(120);not null" json:"searchedValue"`
	SearchedNeighborhood  *string           `gorm:"column:searchedNeighborhood;type:varchar(80)" json:"searchedNeighborhood,omitempty"`
	MatchedNeighborhoodID *string           `gorm:"column:matchedNeighborhoodId" json:"matchedNeighborhoodId,omitempty"`
	DeliveryAvailable     bool              `gorm:"column:deliveryAvailable;not null" json:"deliveryAvailable"`
	DeliveryPrice         float64           `gorm:"column:deliveryPrice;not null" json:"deliveryPrice"`
	ResponseTimeMs        int               `gorm:"column:responseTimeMs;not null" json:"responseTimeMs"`
	SessionID             string            `gorm:"column:sessionId;type:uuid;not null" json:"sessionId"`
	IPHash                string            `gorm:"column:ipHash;type:varchar(64);not null" json:"ipHash"`
	UserAgent             string            `gorm:"column:userAgent;not null" json:"userAgent"`

	// Relacionamentos GORM
	Store               *Store            `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	MatchedNeighborhood *BaseNeighborhood `gorm:"foreignKey:MatchedNeighborhoodID;references:ID" json:"matchedNeighborhood,omitempty"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (SearchEvent) TableName() string {
	return "SearchEvent"
}
