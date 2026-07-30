package service

import (
	"context"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
)

type CommonService interface {
	GetLinesOfBusiness(ctx context.Context) ([]model.LineOfBusiness, error)
	GetBaseNeighborhoodByName(ctx context.Context, name string) (*model.BaseNeighborhood, error)
}

type commonService struct {
	lobRepo  repository.LineOfBusinessRepository
	baseRepo repository.BaseNeighborhoodRepository
}

func NewCommonService(
	lobRepo repository.LineOfBusinessRepository,
	baseRepo repository.BaseNeighborhoodRepository,
) CommonService {
	return &commonService{
		lobRepo:  lobRepo,
		baseRepo: baseRepo,
	}
}

func (s *commonService) GetLinesOfBusiness(ctx context.Context) ([]model.LineOfBusiness, error) {
	return s.lobRepo.FindActive(ctx)
}

func (s *commonService) GetBaseNeighborhoodByName(ctx context.Context, name string) (*model.BaseNeighborhood, error) {
	return s.baseRepo.FindByName(ctx, name)
}
