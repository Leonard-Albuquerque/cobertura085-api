package repository

import (
	"context"
	"errors"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"gorm.io/gorm"
)

var (
	ErrStoreNotFound = errors.New("loja não encontrada")
)

type StoreRepository interface {
	Create(ctx context.Context, store *model.Store) error
	CreateTx(ctx context.Context, tx *gorm.DB, store *model.Store) error
	FindByNameOrSlug(ctx context.Context, name, slug string) (*model.Store, error)
	FindByID(ctx context.Context, id string) (*model.Store, error)
	FindBySlug(ctx context.Context, slug string) (*model.Store, error)
	FindBySlugWithPickupPoints(ctx context.Context, slug string) (*model.Store, []model.PickupPoint, error)
	FindByQRToken(ctx context.Context, token string) (*model.Store, error)
	FindAllPublic(ctx context.Context) ([]dto.PublicStoreResponse, error)
	UpdateSettingsWithTx(ctx context.Context, storeID string, updateStore *model.Store, pickupPoints []model.PickupPoint) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) CreateTx(ctx context.Context, tx *gorm.DB, store *model.Store) error {
	return tx.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) FindByNameOrSlug(ctx context.Context, name, slug string) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) OR slug = ?", name, slug).
		First(&store).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStoreNotFound
	}
	return &store, err
}

func (r *storeRepository) FindByID(ctx context.Context, id string) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&store).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStoreNotFound
	}
	return &store, err
}

func (r *storeRepository) FindBySlug(ctx context.Context, slug string) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&store).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStoreNotFound
	}
	return &store, err
}

func (r *storeRepository) FindBySlugWithPickupPoints(ctx context.Context, slug string) (*model.Store, []model.PickupPoint, error) {
	var store model.Store
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&store).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, ErrStoreNotFound
	}
	if err != nil {
		return nil, nil, err
	}

	var pickupPoints []model.PickupPoint
	err = r.db.WithContext(ctx).
		Where("\"storeId\" = ?", store.ID).
		Find(&pickupPoints).Error

	return &store, pickupPoints, err
}

func (r *storeRepository) FindByQRToken(ctx context.Context, token string) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).
		Where("\"qrToken\" = ?", token).
		First(&store).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStoreNotFound
	}
	return &store, err
}

func (r *storeRepository) FindAllPublic(ctx context.Context) ([]dto.PublicStoreResponse, error) {
	var stores []model.Store
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&stores).Error

	if err != nil {
		return nil, err
	}

	result := make([]dto.PublicStoreResponse, 0, len(stores))
	for _, s := range stores {
		// Checa se a loja tem ao menos um bairro com entrega ativada
		var deliveryCount int64
		r.db.WithContext(ctx).
			Model(&model.Neighborhood{}).
			Where("\"storeId\" = ? AND \"deliveryEnabled\" = ?", s.ID, true).
			Count(&deliveryCount)

		// Contagem de pontos de retirada
		var pickupCount int64
		r.db.WithContext(ctx).
			Model(&model.PickupPoint{}).
			Where("\"storeId\" = ?", s.ID).
			Count(&pickupCount)

		result = append(result, dto.PublicStoreResponse{
			ID:                s.ID,
			Slug:              s.Slug,
			Name:              s.Name,
			LogoURL:           s.LogoURL,
			Address:           s.Address,
			OperatingHours:    s.OperatingHours,
			PickupEnabled:     s.PickupEnabled,
			HasDelivery:       deliveryCount > 0,
			PickupPointsCount: pickupCount,
		})
	}

	return result, nil
}

func (r *storeRepository) UpdateSettingsWithTx(ctx context.Context, storeID string, updateStore *model.Store, pickupPoints []model.PickupPoint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Atualizar campos da loja
		if err := tx.Model(&model.Store{}).Where("id = ?", storeID).Updates(updateStore).Error; err != nil {
			return err
		}

		// 2. Apagar pontos de retirada antigos
		if err := tx.Where("\"storeId\" = ?", storeID).Delete(&model.PickupPoint{}).Error; err != nil {
			return err
		}

		// 3. Inserir novos pontos de retirada se existirem
		if len(pickupPoints) > 0 {
			if err := tx.Create(&pickupPoints).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
