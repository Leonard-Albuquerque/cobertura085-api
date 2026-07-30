package service

import (
	"context"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
)

type NeighborhoodService interface {
	GetStoreNeighborhoods(ctx context.Context, storeID string) ([]dto.NeighborhoodResponse, error)
	CheckStoreNeighborhood(ctx context.Context, storeID, baseNeighborhoodID string) (*dto.NeighborhoodResponse, error)
	UpdateNeighborhood(ctx context.Context, id string, input dto.UpdateNeighborhoodInput) error
	UpdateNeighborhoodsBulk(ctx context.Context, input dto.BulkUpdateNeighborhoodsInput) error
}

type neighborhoodService struct {
	repo repository.NeighborhoodRepository
}

func NewNeighborhoodService(repo repository.NeighborhoodRepository) NeighborhoodService {
	return &neighborhoodService{repo: repo}
}

func (s *neighborhoodService) GetStoreNeighborhoods(ctx context.Context, storeID string) ([]dto.NeighborhoodResponse, error) {
	return s.repo.FindByStoreID(ctx, storeID)
}

func (s *neighborhoodService) CheckStoreNeighborhood(ctx context.Context, storeID, baseNeighborhoodID string) (*dto.NeighborhoodResponse, error) {
	return s.repo.FindStoreNeighborhoodByBaseID(ctx, storeID, baseNeighborhoodID)
}

func (s *neighborhoodService) UpdateNeighborhood(ctx context.Context, id string, input dto.UpdateNeighborhoodInput) error {
	return s.repo.Update(ctx, id, input)
}

func (s *neighborhoodService) UpdateNeighborhoodsBulk(ctx context.Context, input dto.BulkUpdateNeighborhoodsInput) error {
	return s.repo.UpdateBulk(ctx, input.IDs, input.Data)
}
