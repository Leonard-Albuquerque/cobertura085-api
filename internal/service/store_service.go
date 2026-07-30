package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var (
	ErrInvalidLineOfBusiness = errors.New("ramo de atuação inválido")
)

type StoreService interface {
	GetPublicStores(ctx context.Context) ([]dto.PublicStoreResponse, error)
	GetStoreBySlug(ctx context.Context, slug string) (*model.Store, []model.PickupPoint, error)
	GetStoreByQRToken(ctx context.Context, token string) (*model.Store, error)
	UpdateStoreSettings(ctx context.Context, storeID string, input dto.UpdateStoreSettingsInput) error
	GetDashboardStats(ctx context.Context, storeID string) (*dto.StoreDashboardStatsResponse, error)
}

type storeService struct {
	storeRepo        repository.StoreRepository
	neighborhoodRepo repository.NeighborhoodRepository
	lobRepo          repository.LineOfBusinessRepository
}

func NewStoreService(
	storeRepo repository.StoreRepository,
	neighborhoodRepo repository.NeighborhoodRepository,
	lobRepo repository.LineOfBusinessRepository,
) StoreService {
	return &storeService{
		storeRepo:        storeRepo,
		neighborhoodRepo: neighborhoodRepo,
		lobRepo:          lobRepo,
	}
}

func (s *storeService) GetPublicStores(ctx context.Context) ([]dto.PublicStoreResponse, error) {
	return s.storeRepo.FindAllPublic(ctx)
}

func (s *storeService) GetStoreBySlug(ctx context.Context, slug string) (*model.Store, []model.PickupPoint, error) {
	return s.storeRepo.FindBySlugWithPickupPoints(ctx, slug)
}

func (s *storeService) GetStoreByQRToken(ctx context.Context, token string) (*model.Store, error) {
	return s.storeRepo.FindByQRToken(ctx, token)
}

func (s *storeService) UpdateStoreSettings(ctx context.Context, storeID string, input dto.UpdateStoreSettingsInput) error {
	// 1. Sanitizar número de WhatsApp (apenas dígitos)
	reg := regexp.MustCompile(`\D`)
	cleanWhatsapp := reg.ReplaceAllString(input.Whatsapp, "")
	if len(cleanWhatsapp) == 0 {
		cleanWhatsapp = input.Whatsapp
	}

	// 2. Validação de Ramo de Atuação (LineOfBusiness)
	var finalLob *string
	var finalCustomLob *string

	if input.LineOfBusiness != nil && *input.LineOfBusiness != "" {
		lobCode := *input.LineOfBusiness
		if lobCode == "other" {
			val := "other"
			finalLob = &val
			finalCustomLob = input.CustomLineOfBusiness
		} else {
			_, err := s.lobRepo.FindByCode(ctx, lobCode)
			if err != nil {
				return ErrInvalidLineOfBusiness
			}
			finalLob = &lobCode
			finalCustomLob = nil
		}
	}

	// 3. Serializar operatingHoursJson se fornecido
	var jsonBytes []byte
	if input.OperatingHoursJSON != nil {
		b, err := json.Marshal(input.OperatingHoursJSON)
		if err == nil {
			jsonBytes = b
		}
	}

	// 4. Montar modelo da Loja para atualização
	storeUpdate := &model.Store{
		Name:                   input.Name,
		Whatsapp:               cleanWhatsapp,
		Address:                input.Address,
		PickupEnabled:          input.PickupEnabled,
		LogoURL:                input.LogoURL,
		BannerURL:              input.BannerURL,
		Description:            input.Description,
		Instagram:              input.Instagram,
		CatalogURL:             input.CatalogURL,
		WebsiteURL:             input.WebsiteURL,
		DeliveryTimeDefault:    input.DeliveryTimeDefault,
		DeliveryAvailableMsg:   input.DeliveryAvailableMsg,
		DeliveryUnavailableMsg: input.DeliveryUnavailableMsg,
		SameDayCutoff:          input.SameDayCutoff,
		CutoffMessage:          input.CutoffMessage,
		LineOfBusiness:         finalLob,
		CustomLineOfBusiness:   finalCustomLob,
		UpdatedAt:              time.Now(),
	}

	if len(jsonBytes) > 0 {
		storeUpdate.OperatingHoursJSON = datatypes.JSON(jsonBytes)
	}

	// 5. Montar novos PickupPoints
	pickupPoints := make([]model.PickupPoint, 0, len(input.PickupPoints))
	for _, p := range input.PickupPoints {
		pickupPoints = append(pickupPoints, model.PickupPoint{
			ID:           uuid.New().String(),
			StoreID:      storeID,
			Name:         p.Name,
			Address:      p.Address,
			Latitude:     p.Latitude,
			Longitude:    p.Longitude,
			Instructions: p.Instructions,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	// 6. Executar atualização atômica no repositório
	return s.storeRepo.UpdateSettingsWithTx(ctx, storeID, storeUpdate, pickupPoints)
}

func (s *storeService) GetDashboardStats(ctx context.Context, storeID string) (*dto.StoreDashboardStatsResponse, error) {
	return s.neighborhoodRepo.GetDashboardStats(ctx, storeID)
}
