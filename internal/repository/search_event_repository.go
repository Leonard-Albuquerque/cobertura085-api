package repository

import (
	"context"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"gorm.io/gorm"
)

type SearchEventRepository interface {
	Create(ctx context.Context, event *model.SearchEvent) error
	GetSummary(ctx context.Context) (*dto.TelemetrySummaryResponse, error)
	GetRecentLogs(ctx context.Context, limit int) ([]dto.SearchEventResponse, error)
	GetBusinessStats(ctx context.Context, storeID string) (*dto.TelemetrySummaryResponse, error)
}

type searchEventRepository struct {
	db *gorm.DB
}

func NewSearchEventRepository(db *gorm.DB) SearchEventRepository {
	return &searchEventRepository{db: db}
}

func (r *searchEventRepository) Create(ctx context.Context, event *model.SearchEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *searchEventRepository) GetSummary(ctx context.Context) (*dto.TelemetrySummaryResponse, error) {
	var totalSearches int64
	r.db.WithContext(ctx).Model(&model.SearchEvent{}).Count(&totalSearches)

	var uniqueVisitors int64
	r.db.WithContext(ctx).Model(&model.SearchEvent{}).Select("COUNT(DISTINCT \"sessionId\")").Count(&uniqueVisitors)

	var availableCount int64
	r.db.WithContext(ctx).Model(&model.SearchEvent{}).Where("\"deliveryAvailable\" = ?", true).Count(&availableCount)

	availableRate := 0.0
	if totalSearches > 0 {
		availableRate = (float64(availableCount) / float64(totalSearches)) * 100
	}

	var avgMs float64
	row := r.db.WithContext(ctx).
		Model(&model.SearchEvent{}).
		Select("COALESCE(AVG(\"responseTimeMs\"), 0)").
		Row()
	_ = row.Scan(&avgMs)

	// Top 10 Bairros Buscados
	type topBairroRow struct {
		SearchedNeighborhood string
		Count                int64
	}
	var topRows []topBairroRow
	r.db.WithContext(ctx).
		Table("\"SearchEvent\"").
		Select("\"searchedNeighborhood\", COUNT(*) as count").
		Where("\"searchedNeighborhood\" IS NOT NULL AND \"searchedNeighborhood\" != ''").
		Group("\"searchedNeighborhood\"").
		Order("count DESC").
		Limit(10).
		Scan(&topRows)

	topNeighborhoods := make([]dto.TopNeighborhoodItem, 0, len(topRows))
	for _, row := range topRows {
		topNeighborhoods = append(topNeighborhoods, dto.TopNeighborhoodItem{
			Neighborhood: row.SearchedNeighborhood,
			Count:        row.Count,
		})
	}

	// Share de Buscas por Loja
	type storeShareRow struct {
		StoreID   string
		StoreName string
		Count     int64
	}
	var shareRows []storeShareRow
	r.db.WithContext(ctx).
		Table("\"SearchEvent\"").
		Select("\"SearchEvent\".\"storeId\", \"Store\".\"name\" as store_name, COUNT(*) as count").
		Joins("LEFT JOIN \"Store\" ON \"SearchEvent\".\"storeId\" = \"Store\".\"id\"").
		Group("\"SearchEvent\".\"storeId\", \"Store\".\"name\"").
		Order("count DESC").
		Limit(10).
		Scan(&shareRows)

	storeShares := make([]dto.StoreSearchShareItem, 0, len(shareRows))
	for _, row := range shareRows {
		storeShares = append(storeShares, dto.StoreSearchShareItem{
			StoreID:   row.StoreID,
			StoreName: row.StoreName,
			Count:     row.Count,
		})
	}

	return &dto.TelemetrySummaryResponse{
		TotalSearches:         totalSearches,
		UniqueVisitors:        uniqueVisitors,
		AvailableRate:         availableRate,
		AverageResponseTimeMs: avgMs,
		TopNeighborhoods:      topNeighborhoods,
		StoreSearchShares:     storeShares,
	}, nil
}

func (r *searchEventRepository) GetRecentLogs(ctx context.Context, limit int) ([]dto.SearchEventResponse, error) {
	if limit <= 0 {
		limit = 50
	}

	var events []model.SearchEvent
	err := r.db.WithContext(ctx).
		Preload("Store").
		Order("\"createdAt\" DESC").
		Limit(limit).
		Find(&events).Error

	if err != nil {
		return nil, err
	}

	result := make([]dto.SearchEventResponse, 0, len(events))
	for _, e := range events {
		var storeName *string
		if e.Store != nil {
			storeName = &e.Store.Name
		}

		result = append(result, dto.SearchEventResponse{
			ID:                   e.ID,
			CreatedAt:            e.CreatedAt,
			StoreID:              e.StoreID,
			StoreName:            storeName,
			EventType:            e.EventType,
			SearchType:           e.SearchType,
			SearchedValue:        e.SearchedValue,
			SearchedNeighborhood: e.SearchedNeighborhood,
			DeliveryAvailable:    e.DeliveryAvailable,
			DeliveryPrice:        e.DeliveryPrice,
			ResponseTimeMs:       e.ResponseTimeMs,
			SessionID:            e.SessionID,
			UserAgent:            e.UserAgent,
		})
	}

	return result, nil
}

func (r *searchEventRepository) GetBusinessStats(ctx context.Context, storeID string) (*dto.TelemetrySummaryResponse, error) {
	var totalSearches int64
	r.db.WithContext(ctx).Model(&model.SearchEvent{}).Where("\"storeId\" = ?", storeID).Count(&totalSearches)

	var uniqueVisitors int64
	r.db.WithContext(ctx).Model(&model.SearchEvent{}).Where("\"storeId\" = ?", storeID).Select("COUNT(DISTINCT \"sessionId\")").Count(&uniqueVisitors)

	var availableCount int64
	r.db.WithContext(ctx).Model(&model.SearchEvent{}).Where("\"storeId\" = ? AND \"deliveryAvailable\" = ?", storeID, true).Count(&availableCount)

	availableRate := 0.0
	if totalSearches > 0 {
		availableRate = (float64(availableCount) / float64(totalSearches)) * 100
	}

	var avgMs float64
	row := r.db.WithContext(ctx).
		Model(&model.SearchEvent{}).
		Select("COALESCE(AVG(\"responseTimeMs\"), 0)").
		Where("\"storeId\" = ?", storeID).
		Row()
	_ = row.Scan(&avgMs)

	// Top 10 Bairros da Loja
	type topBairroRow struct {
		SearchedNeighborhood string
		Count                int64
	}
	var topRows []topBairroRow
	r.db.WithContext(ctx).
		Table("\"SearchEvent\"").
		Select("\"searchedNeighborhood\", COUNT(*) as count").
		Where("\"storeId\" = ? AND \"searchedNeighborhood\" IS NOT NULL AND \"searchedNeighborhood\" != ''", storeID).
		Group("\"searchedNeighborhood\"").
		Order("count DESC").
		Limit(10).
		Scan(&topRows)

	topNeighborhoods := make([]dto.TopNeighborhoodItem, 0, len(topRows))
	for _, row := range topRows {
		topNeighborhoods = append(topNeighborhoods, dto.TopNeighborhoodItem{
			Neighborhood: row.SearchedNeighborhood,
			Count:        row.Count,
		})
	}

	return &dto.TelemetrySummaryResponse{
		TotalSearches:         totalSearches,
		UniqueVisitors:        uniqueVisitors,
		AvailableRate:         availableRate,
		AverageResponseTimeMs: avgMs,
		TopNeighborhoods:      topNeighborhoods,
		StoreSearchShares:     []dto.StoreSearchShareItem{},
	}, nil
}
