package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"gorm.io/gorm"
)

var (
	ErrNeighborhoodNotFound = errors.New("bairro não encontrado")
)

type NeighborhoodRepository interface {
	FindByStoreID(ctx context.Context, storeID string) ([]dto.NeighborhoodResponse, error)
	FindStoreNeighborhoodByBaseID(ctx context.Context, storeID, baseNeighborhoodID string) (*dto.NeighborhoodResponse, error)
	Update(ctx context.Context, id string, data dto.UpdateNeighborhoodInput) error
	UpdateBulk(ctx context.Context, ids []string, data dto.UpdateNeighborhoodInput) error
	GetDashboardStats(ctx context.Context, storeID string) (*dto.StoreDashboardStatsResponse, error)
}

type neighborhoodRepository struct {
	db *gorm.DB
}

func NewNeighborhoodRepository(db *gorm.DB) NeighborhoodRepository {
	return &neighborhoodRepository{db: db}
}

func (r *neighborhoodRepository) FindByStoreID(ctx context.Context, storeID string) ([]dto.NeighborhoodResponse, error) {
	var neighborhoods []model.Neighborhood
	err := r.db.WithContext(ctx).
		Preload("BaseNeighborhood").
		Joins("LEFT JOIN \"BaseNeighborhood\" ON \"Neighborhood\".\"baseNeighborhoodId\" = \"BaseNeighborhood\".\"id\"").
		Where("\"Neighborhood\".\"storeId\" = ?", storeID).
		Order("\"BaseNeighborhood\".\"officialName\" ASC").
		Find(&neighborhoods).Error

	if err != nil {
		return nil, err
	}

	result := make([]dto.NeighborhoodResponse, 0, len(neighborhoods))
	for _, n := range neighborhoods {
		var baseResp *dto.BaseNeighborhoodResponse
		if n.BaseNeighborhood != nil {
			baseResp = &dto.BaseNeighborhoodResponse{
				ID:           n.BaseNeighborhood.ID,
				Name:         n.BaseNeighborhood.Name,
				OfficialName: n.BaseNeighborhood.OfficialName,
			}
		}

		result = append(result, dto.NeighborhoodResponse{
			ID:                    n.ID,
			StoreID:               n.StoreID,
			BaseNeighborhoodID:    n.BaseNeighborhoodID,
			DeliveryEnabled:       n.DeliveryEnabled,
			Fee:                   n.Fee,
			DeliveryTime:          n.DeliveryTime,
			MinimumOrder:          n.MinimumOrder,
			FreeDeliveryThreshold: n.FreeDeliveryThreshold,
			Notes:                 n.Notes,
			CreatedAt:             n.CreatedAt,
			UpdatedAt:             n.UpdatedAt,
			BaseNeighborhood:      baseResp,
		})
	}

	return result, nil
}

func (r *neighborhoodRepository) FindStoreNeighborhoodByBaseID(ctx context.Context, storeID, baseNeighborhoodID string) (*dto.NeighborhoodResponse, error) {
	var n model.Neighborhood
	err := r.db.WithContext(ctx).
		Preload("BaseNeighborhood").
		Where("\"storeId\" = ? AND \"baseNeighborhoodId\" = ?", storeID, baseNeighborhoodID).
		First(&n).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNeighborhoodNotFound
	}
	if err != nil {
		return nil, err
	}

	var baseResp *dto.BaseNeighborhoodResponse
	if n.BaseNeighborhood != nil {
		baseResp = &dto.BaseNeighborhoodResponse{
			ID:           n.BaseNeighborhood.ID,
			Name:         n.BaseNeighborhood.Name,
			OfficialName: n.BaseNeighborhood.OfficialName,
		}
	}

	return &dto.NeighborhoodResponse{
		ID:                    n.ID,
		StoreID:               n.StoreID,
		BaseNeighborhoodID:    n.BaseNeighborhoodID,
		DeliveryEnabled:       n.DeliveryEnabled,
		Fee:                   n.Fee,
		DeliveryTime:          n.DeliveryTime,
		MinimumOrder:          n.MinimumOrder,
		FreeDeliveryThreshold: n.FreeDeliveryThreshold,
		Notes:                 n.Notes,
		CreatedAt:             n.CreatedAt,
		UpdatedAt:             n.UpdatedAt,
		BaseNeighborhood:      baseResp,
	}, nil
}

func (r *neighborhoodRepository) Update(ctx context.Context, id string, data dto.UpdateNeighborhoodInput) error {
	updates := make(map[string]interface{})
	if data.DeliveryEnabled != nil {
		updates["deliveryEnabled"] = *data.DeliveryEnabled
	}
	if data.Fee != nil {
		updates["fee"] = *data.Fee
	}
	if data.DeliveryTime != nil {
		updates["deliveryTime"] = *data.DeliveryTime
	}
	if data.MinimumOrder != nil {
		updates["minimumOrder"] = *data.MinimumOrder
	}
	if data.FreeDeliveryThreshold != nil {
		updates["freeDeliveryThreshold"] = *data.FreeDeliveryThreshold
	}
	if data.Notes != nil {
		updates["notes"] = *data.Notes
	}
	updates["updatedAt"] = time.Now()

	res := r.db.WithContext(ctx).Model(&model.Neighborhood{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNeighborhoodNotFound
	}
	return nil
}

func (r *neighborhoodRepository) UpdateBulk(ctx context.Context, ids []string, data dto.UpdateNeighborhoodInput) error {
	if len(ids) == 0 {
		return nil
	}

	updates := make(map[string]interface{})
	if data.DeliveryEnabled != nil {
		updates["deliveryEnabled"] = *data.DeliveryEnabled
	}
	if data.Fee != nil {
		updates["fee"] = *data.Fee
	}
	if data.DeliveryTime != nil {
		updates["deliveryTime"] = *data.DeliveryTime
	}
	if data.MinimumOrder != nil {
		updates["minimumOrder"] = *data.MinimumOrder
	}
	if data.FreeDeliveryThreshold != nil {
		updates["freeDeliveryThreshold"] = *data.FreeDeliveryThreshold
	}
	if data.Notes != nil {
		updates["notes"] = *data.Notes
	}
	updates["updatedAt"] = time.Now()

	return r.db.WithContext(ctx).Model(&model.Neighborhood{}).Where("id IN ?", ids).Updates(updates).Error
}

func (r *neighborhoodRepository) GetDashboardStats(ctx context.Context, storeID string) (*dto.StoreDashboardStatsResponse, error) {
	var activeCount int64
	err := r.db.WithContext(ctx).
		Model(&model.Neighborhood{}).
		Where("\"storeId\" = ? AND \"deliveryEnabled\" = ?", storeID, true).
		Count(&activeCount).Error
	if err != nil {
		return nil, err
	}

	var inactiveCount int64
	err = r.db.WithContext(ctx).
		Model(&model.Neighborhood{}).
		Where("\"storeId\" = ? AND \"deliveryEnabled\" = ?", storeID, false).
		Count(&inactiveCount).Error
	if err != nil {
		return nil, err
	}

	var avgFee float64
	row := r.db.WithContext(ctx).
		Model(&model.Neighborhood{}).
		Select("COALESCE(AVG(fee), 0)").
		Where("\"storeId\" = ? AND \"deliveryEnabled\" = ?", storeID, true).
		Row()
	_ = row.Scan(&avgFee)

	var lastUpdated model.Neighborhood
	err = r.db.WithContext(ctx).
		Where("\"storeId\" = ?", storeID).
		Order("\"updatedAt\" DESC").
		First(&lastUpdated).Error

	var lastUpdatedAt *time.Time
	if err == nil {
		lastUpdatedAt = &lastUpdated.UpdatedAt
	}

	return &dto.StoreDashboardStatsResponse{
		ActiveNeighborhoods:   activeCount,
		InactiveNeighborhoods: inactiveCount,
		AverageFee:            avgFee,
		LastUpdatedAt:         lastUpdatedAt,
	}, nil
}
