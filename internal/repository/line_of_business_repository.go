package repository

import (
	"context"
	"errors"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"gorm.io/gorm"
)

var (
	ErrLineOfBusinessNotFound = errors.New("ramo de atuação não encontrado")
)

type LineOfBusinessRepository interface {
	FindActive(ctx context.Context) ([]model.LineOfBusiness, error)
	FindByCode(ctx context.Context, code string) (*model.LineOfBusiness, error)
}

type lineOfBusinessRepository struct {
	db *gorm.DB
}

func NewLineOfBusinessRepository(db *gorm.DB) LineOfBusinessRepository {
	return &lineOfBusinessRepository{db: db}
}

func (r *lineOfBusinessRepository) FindActive(ctx context.Context) ([]model.LineOfBusiness, error) {
	var list []model.LineOfBusiness
	err := r.db.WithContext(ctx).
		Where("\"isActive\" = ?", true).
		Order("\"sortOrder\" ASC, name ASC").
		Find(&list).Error
	return list, err
}

func (r *lineOfBusinessRepository) FindByCode(ctx context.Context, code string) (*model.LineOfBusiness, error) {
	var lob model.LineOfBusiness
	err := r.db.WithContext(ctx).
		Where("code = ? AND \"isActive\" = ?", code, true).
		First(&lob).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLineOfBusinessNotFound
	}
	return &lob, err
}
