package model

import (
	"time"

	"gorm.io/datatypes"
)

// Store representa a tabela "Store" no banco de dados.
type Store struct {
	ID                     string         `gorm:"primaryKey;column:id" json:"id"`
	Slug                   string         `gorm:"column:slug;not null;unique" json:"slug"`
	Name                   string         `gorm:"column:name;not null;default:'Minha Loja'" json:"name"`
	LogoURL                *string        `gorm:"column:logoUrl" json:"logoUrl,omitempty"`
	Whatsapp               string         `gorm:"column:whatsapp;not null;default:'5585999999999'" json:"whatsapp"`
	Address                string         `gorm:"column:address;not null;default:'Endereço da Loja, Fortaleza - CE'" json:"address"`
	OperatingHours         string         `gorm:"column:operatingHours;not null;default:'Segunda a Sexta: 08:00 às 18:00'" json:"operatingHours"`
	PickupEnabled          bool           `gorm:"column:pickupEnabled;not null;default:true" json:"pickupEnabled"`
	CreatedAt              time.Time      `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updatedAt;not null" json:"updatedAt"`
	BannerURL              *string        `gorm:"column:bannerUrl" json:"bannerUrl,omitempty"`
	CatalogURL             *string        `gorm:"column:catalogUrl" json:"catalogUrl,omitempty"`
	CutoffMessage          *string        `gorm:"column:cutoffMessage" json:"cutoffMessage,omitempty"`
	DeliveryAvailableMsg   *string        `gorm:"column:deliveryAvailableMsg" json:"deliveryAvailableMsg,omitempty"`
	DeliveryTimeDefault    string         `gorm:"column:deliveryTimeDefault;not null;default:'2 horas'" json:"deliveryTimeDefault"`
	DeliveryUnavailableMsg *string        `gorm:"column:deliveryUnavailableMsg" json:"deliveryUnavailableMsg,omitempty"`
	Description            *string        `gorm:"column:description" json:"description,omitempty"`
	Instagram              *string        `gorm:"column:instagram" json:"instagram,omitempty"`
	OperatingHoursJSON     datatypes.JSON `gorm:"column:operatingHoursJson;type:jsonb" json:"operatingHoursJson,omitempty"`
	SameDayCutoff          *string        `gorm:"column:sameDayCutoff" json:"sameDayCutoff,omitempty"`
	WebsiteURL             *string        `gorm:"column:websiteUrl" json:"websiteUrl,omitempty"`
	CustomLineOfBusiness   *string        `gorm:"column:customLineOfBusiness" json:"customLineOfBusiness,omitempty"`
	LineOfBusiness         *string        `gorm:"column:lineOfBusiness" json:"lineOfBusiness,omitempty"`
	QRToken                string         `gorm:"column:qrToken;not null;unique" json:"qrToken"`
}

// TableName define o nome exato da tabela no PostgreSQL.
func (Store) TableName() string {
	return "Store"
}
