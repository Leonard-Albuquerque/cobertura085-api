package repository

import (
	"context"
	"errors"

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
