package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/google/uuid"
)

type TelemetryService interface {
	LogSearchEvent(ctx context.Context, input dto.CreateSearchEventInput, clientIP string) error
	GetSummary(ctx context.Context) (*dto.TelemetrySummaryResponse, error)
	GetRecentLogs(ctx context.Context, limit int) ([]dto.SearchEventResponse, error)
	GetBusinessStats(ctx context.Context, storeID string) (*dto.TelemetrySummaryResponse, error)
}

type telemetryService struct {
	repo repository.SearchEventRepository
}

func NewTelemetryService(repo repository.SearchEventRepository) TelemetryService {
	return &telemetryService{repo: repo}
}

func (s *telemetryService) LogSearchEvent(ctx context.Context, input dto.CreateSearchEventInput, clientIP string) error {
	ipHash := input.IPHash
	if ipHash == "" && clientIP != "" {
		hash := sha256.Sum256([]byte(clientIP))
		ipHash = hex.EncodeToString(hash[:])
	}

	event := &model.SearchEvent{
		ID:                    uuid.New().String(),
		CreatedAt:             time.Now(),
		StoreID:               input.StoreID,
		EventType:             input.EventType,
		SearchType:            input.SearchType,
		SearchedValue:         input.SearchedValue,
		SearchedNeighborhood:  input.SearchedNeighborhood,
		MatchedNeighborhoodID: input.MatchedNeighborhoodID,
		DeliveryAvailable:     input.DeliveryAvailable,
		DeliveryPrice:         input.DeliveryPrice,
		ResponseTimeMs:        input.ResponseTimeMs,
		SessionID:             input.SessionID,
		IPHash:                ipHash,
		UserAgent:             input.UserAgent,
	}

	return s.repo.Create(ctx, event)
}

func (s *telemetryService) GetSummary(ctx context.Context) (*dto.TelemetrySummaryResponse, error) {
	return s.repo.GetSummary(ctx)
}

func (s *telemetryService) GetRecentLogs(ctx context.Context, limit int) ([]dto.SearchEventResponse, error) {
	return s.repo.GetRecentLogs(ctx, limit)
}

func (s *telemetryService) GetBusinessStats(ctx context.Context, storeID string) (*dto.TelemetrySummaryResponse, error) {
	return s.repo.GetBusinessStats(ctx, storeID)
}
