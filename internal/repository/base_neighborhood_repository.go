package repository

import (
	"context"
	"errors"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"gorm.io/gorm"
)

var (
	ErrBaseNeighborhoodNotFound = errors.New("bairro base não encontrado")
)

type BaseNeighborhoodRepository interface {
	FindByName(ctx context.Context, name string) (*model.BaseNeighborhood, error)
	FindAll(ctx context.Context) ([]model.BaseNeighborhood, error)
}

type baseNeighborhoodRepository struct {
	db *gorm.DB
}

func NewBaseNeighborhoodRepository(db *gorm.DB) BaseNeighborhoodRepository {
	return &baseNeighborhoodRepository{db: db}
}

func (r *baseNeighborhoodRepository) FindByName(ctx context.Context, name string) (*model.BaseNeighborhood, error) {
	var bn model.BaseNeighborhood
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&bn).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBaseNeighborhoodNotFound
	}
	return &bn, err
}

func (r *baseNeighborhoodRepository) FindAll(ctx context.Context) ([]model.BaseNeighborhood, error) {
	var list []model.BaseNeighborhood
	err := r.db.WithContext(ctx).
		Order("\"officialName\" ASC").
		Find(&list).Error
	return list, err
}
