package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/dto"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service/external"
	"github.com/google/uuid"
)

type ShippingService interface {
	LookupCEP(ctx context.Context, input dto.LookupCEPInput, clientIP, userAgent string) (*dto.LookupResultResponse, error)
	LookupAddress(ctx context.Context, input dto.LookupAddressInput, clientIP, userAgent string) (*dto.LookupResultResponse, error)
	LookupCoords(ctx context.Context, input dto.LookupCoordsInput, clientIP, userAgent string) (*dto.LookupResultResponse, error)
	LookupSelectedAddress(ctx context.Context, input dto.LookupSelectedAddressInput, clientIP, userAgent string) (*dto.LookupResultResponse, error)
	GetAddressSuggestions(ctx context.Context, query string) ([]dto.AddressSuggestionItem, error)
}

type shippingService struct {
	viaCEPClient         external.ViaCEPClient
	nominatimClient      external.NominatimClient
	storeRepo            repository.StoreRepository
	baseNeighborhoodRepo repository.BaseNeighborhoodRepository
	neighborhoodRepo     repository.NeighborhoodRepository
	telemetryService     TelemetryService
}

func NewShippingService(
	viaCEPClient external.ViaCEPClient,
	nominatimClient external.NominatimClient,
	storeRepo repository.StoreRepository,
	baseNeighborhoodRepo repository.BaseNeighborhoodRepository,
	neighborhoodRepo repository.NeighborhoodRepository,
	telemetryService TelemetryService,
) ShippingService {
	return &shippingService{
		viaCEPClient:         viaCEPClient,
		nominatimClient:      nominatimClient,
		storeRepo:            storeRepo,
		baseNeighborhoodRepo: baseNeighborhoodRepo,
		neighborhoodRepo:     neighborhoodRepo,
		telemetryService:     telemetryService,
	}
}

func (s *shippingService) LookupCEP(ctx context.Context, input dto.LookupCEPInput, clientIP, userAgent string) (*dto.LookupResultResponse, error) {
	startTime := time.Now()

	store, err := s.storeRepo.FindBySlug(ctx, input.StoreSlug)
	if err != nil {
		return &dto.LookupResultResponse{Success: false, Error: "Loja não encontrada."}, nil
	}

	sessionID := uuid.New().String()

	viaCepData, err := s.viaCEPClient.Lookup(ctx, input.CEP)
	if err != nil {
		ms := int(time.Since(startTime).Milliseconds())
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:           store.ID,
			EventType:         "SEARCH",
			SearchType:        "CEP",
			SearchedValue:     input.CEP,
			DeliveryAvailable: false,
			DeliveryPrice:     0,
			ResponseTimeMs:    ms,
			SessionID:         sessionID,
			UserAgent:         userAgent,
		}, clientIP)
		return &dto.LookupResultResponse{Success: false, Error: "CEP não encontrado ou erro na consulta."}, nil
	}

	city := strings.ToLower(viaCepData.Localidade)
	if city != "fortaleza" {
		ms := int(time.Since(startTime).Milliseconds())
		rawBairro := viaCepData.Bairro
		if rawBairro == "" {
			rawBairro = "Fora de Fortaleza"
		}
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "CEP",
			SearchedValue:        input.CEP,
			SearchedNeighborhood: &rawBairro,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			DeliveryEnabled: false,
			Bairro:          rawBairro,
			Error:           "Infelizmente entregamos apenas em Fortaleza (CE).",
		}, nil
	}

	rawBairro := viaCepData.Bairro
	if rawBairro == "" {
		return &dto.LookupResultResponse{Success: false, Error: "Não foi possível identificar o bairro para este CEP."}, nil
	}

	normalizedBairro := external.NormalizeName(rawBairro)
	baseBairro, err := s.baseNeighborhoodRepo.FindByName(ctx, normalizedBairro)
	if err != nil || baseBairro == nil {
		ms := int(time.Since(startTime).Milliseconds())
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "CEP",
			SearchedValue:        input.CEP,
			SearchedNeighborhood: &rawBairro,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          rawBairro,
			Street:          viaCepData.Logradouro,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	neighborhood, err := s.neighborhoodRepo.FindStoreNeighborhoodByBaseID(ctx, store.ID, baseBairro.ID)
	ms := int(time.Since(startTime).Milliseconds())

	if err != nil || neighborhood == nil {
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:               store.ID,
			EventType:             "SEARCH",
			SearchType:            "CEP",
			SearchedValue:         input.CEP,
			SearchedNeighborhood:  &baseBairro.OfficialName,
			MatchedNeighborhoodID: &baseBairro.ID,
			DeliveryAvailable:     false,
			DeliveryPrice:         0,
			ResponseTimeMs:        ms,
			SessionID:             sessionID,
			UserAgent:             userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          baseBairro.OfficialName,
			Street:          viaCepData.Logradouro,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	deliveryTime := "24h"
	if neighborhood.DeliveryTime != nil && *neighborhood.DeliveryTime != "" {
		deliveryTime = *neighborhood.DeliveryTime
	}

	_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
		StoreID:               store.ID,
		EventType:             "SEARCH",
		SearchType:            "CEP",
		SearchedValue:         input.CEP,
		SearchedNeighborhood:  &baseBairro.OfficialName,
		MatchedNeighborhoodID: &baseBairro.ID,
		DeliveryAvailable:     neighborhood.DeliveryEnabled,
		DeliveryPrice:         neighborhood.Fee,
		ResponseTimeMs:        ms,
		SessionID:             sessionID,
		UserAgent:             userAgent,
	}, clientIP)

	return &dto.LookupResultResponse{
		Success:               true,
		Bairro:                baseBairro.OfficialName,
		Street:                viaCepData.Logradouro,
		DeliveryEnabled:       neighborhood.DeliveryEnabled,
		Fee:                   &neighborhood.Fee,
		DeliveryTime:          deliveryTime,
		MinimumOrder:          neighborhood.MinimumOrder,
		FreeDeliveryThreshold: neighborhood.FreeDeliveryThreshold,
		Notes:                 neighborhood.Notes,
		StoreAddress:          store.Address,
		StoreWhatsapp:         store.Whatsapp,
		PickupEnabled:         store.PickupEnabled,
	}, nil
}

func (s *shippingService) LookupAddress(ctx context.Context, input dto.LookupAddressInput, clientIP, userAgent string) (*dto.LookupResultResponse, error) {
	startTime := time.Now()

	if len(strings.TrimSpace(input.Address)) < 5 {
		return &dto.LookupResultResponse{Success: false, Error: "Endereço muito curto. Digite o nome da rua e número."}, nil
	}

	store, err := s.storeRepo.FindBySlug(ctx, input.StoreSlug)
	if err != nil {
		return &dto.LookupResultResponse{Success: false, Error: "Loja não encontrada."}, nil
	}

	sessionID := uuid.New().String()

	results, err := s.nominatimClient.SearchAddress(ctx, input.Address)
	if err != nil || len(results) == 0 {
		ms := int(time.Since(startTime).Milliseconds())
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:           store.ID,
			EventType:         "SEARCH",
			SearchType:        "ADDRESS",
			SearchedValue:     input.Address,
			DeliveryAvailable: false,
			DeliveryPrice:     0,
			ResponseTimeMs:    ms,
			SessionID:         sessionID,
			UserAgent:         userAgent,
		}, clientIP)
		return &dto.LookupResultResponse{Success: false, Error: "Endereço não localizado em Fortaleza. Tente incluir o número ou ajustar a grafia."}, nil
	}

	result := results[0]
	city := strings.ToLower(result.Address.City)
	if city == "" {
		city = strings.ToLower(result.Address.Town)
	}
	if city == "" {
		city = strings.ToLower(result.Address.Municipality)
	}

	if city != "fortaleza" {
		ms := int(time.Since(startTime).Milliseconds())
		fora := "Fora de Fortaleza"
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "ADDRESS",
			SearchedValue:        input.Address,
			SearchedNeighborhood: &fora,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)
		return &dto.LookupResultResponse{
			Success:         true,
			DeliveryEnabled: false,
			Error:           "Infelizmente entregamos apenas em Fortaleza (CE).",
		}, nil
	}

	rawBairro := external.ExtractBairro(result)
	if rawBairro == "" {
		return &dto.LookupResultResponse{Success: false, Error: "Bairro não identificado. Tente pesquisar informando o CEP."}, nil
	}

	normalizedBairro := external.NormalizeName(rawBairro)
	baseBairro, err := s.baseNeighborhoodRepo.FindByName(ctx, normalizedBairro)
	street := result.Address.Road

	if err != nil || baseBairro == nil {
		ms := int(time.Since(startTime).Milliseconds())
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "ADDRESS",
			SearchedValue:        input.Address,
			SearchedNeighborhood: &rawBairro,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          rawBairro,
			Street:          street,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	neighborhood, err := s.neighborhoodRepo.FindStoreNeighborhoodByBaseID(ctx, store.ID, baseBairro.ID)
	ms := int(time.Since(startTime).Milliseconds())

	if err != nil || neighborhood == nil {
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:               store.ID,
			EventType:             "SEARCH",
			SearchType:            "ADDRESS",
			SearchedValue:         input.Address,
			SearchedNeighborhood:  &baseBairro.OfficialName,
			MatchedNeighborhoodID: &baseBairro.ID,
			DeliveryAvailable:     false,
			DeliveryPrice:         0,
			ResponseTimeMs:        ms,
			SessionID:             sessionID,
			UserAgent:             userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          baseBairro.OfficialName,
			Street:          street,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	deliveryTime := "24h"
	if neighborhood.DeliveryTime != nil && *neighborhood.DeliveryTime != "" {
		deliveryTime = *neighborhood.DeliveryTime
	}

	_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
		StoreID:               store.ID,
		EventType:             "SEARCH",
		SearchType:            "ADDRESS",
		SearchedValue:         input.Address,
		SearchedNeighborhood:  &baseBairro.OfficialName,
		MatchedNeighborhoodID: &baseBairro.ID,
		DeliveryAvailable:     neighborhood.DeliveryEnabled,
		DeliveryPrice:         neighborhood.Fee,
		ResponseTimeMs:        ms,
		SessionID:             sessionID,
		UserAgent:             userAgent,
	}, clientIP)

	return &dto.LookupResultResponse{
		Success:               true,
		Bairro:                baseBairro.OfficialName,
		Street:                street,
		DeliveryEnabled:       neighborhood.DeliveryEnabled,
		Fee:                   &neighborhood.Fee,
		DeliveryTime:          deliveryTime,
		MinimumOrder:          neighborhood.MinimumOrder,
		FreeDeliveryThreshold: neighborhood.FreeDeliveryThreshold,
		Notes:                 neighborhood.Notes,
		StoreAddress:          store.Address,
		StoreWhatsapp:         store.Whatsapp,
		PickupEnabled:         store.PickupEnabled,
	}, nil
}

func (s *shippingService) LookupCoords(ctx context.Context, input dto.LookupCoordsInput, clientIP, userAgent string) (*dto.LookupResultResponse, error) {
	startTime := time.Now()

	store, err := s.storeRepo.FindBySlug(ctx, input.StoreSlug)
	if err != nil {
		return &dto.LookupResultResponse{Success: false, Error: "Loja não encontrada."}, nil
	}

	sessionID := uuid.New().String()
	coordStr := fmt.Sprintf("%.6f, %.6f", input.Lat, input.Lon)

	result, err := s.nominatimClient.ReverseGeocode(ctx, input.Lat, input.Lon)
	if err != nil || result == nil {
		ms := int(time.Since(startTime).Milliseconds())
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:           store.ID,
			EventType:         "SEARCH",
			SearchType:        "LOCATION",
			SearchedValue:     coordStr,
			DeliveryAvailable: false,
			DeliveryPrice:     0,
			ResponseTimeMs:    ms,
			SessionID:         sessionID,
			UserAgent:         userAgent,
		}, clientIP)
		return &dto.LookupResultResponse{Success: false, Error: "Localização não identificada."}, nil
	}

	city := strings.ToLower(result.Address.City)
	if city == "" {
		city = strings.ToLower(result.Address.Town)
	}
	if city == "" {
		city = strings.ToLower(result.Address.Municipality)
	}

	if city != "fortaleza" {
		ms := int(time.Since(startTime).Milliseconds())
		fora := "Fora de Fortaleza"
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "LOCATION",
			SearchedValue:        coordStr,
			SearchedNeighborhood: &fora,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)
		return &dto.LookupResultResponse{
			Success:         true,
			DeliveryEnabled: false,
			Error:           "Infelizmente entregamos apenas em Fortaleza (CE).",
		}, nil
	}

	rawBairro := external.ExtractBairro(*result)
	if rawBairro == "" {
		return &dto.LookupResultResponse{Success: false, Error: "Bairro não identificado na sua localização. Tente pesquisar por CEP."}, nil
	}

	normalizedBairro := external.NormalizeName(rawBairro)
	baseBairro, err := s.baseNeighborhoodRepo.FindByName(ctx, normalizedBairro)
	street := result.Address.Road

	if err != nil || baseBairro == nil {
		ms := int(time.Since(startTime).Milliseconds())
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "LOCATION",
			SearchedValue:        coordStr,
			SearchedNeighborhood: &rawBairro,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          rawBairro,
			Street:          street,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	neighborhood, err := s.neighborhoodRepo.FindStoreNeighborhoodByBaseID(ctx, store.ID, baseBairro.ID)
	ms := int(time.Since(startTime).Milliseconds())

	if err != nil || neighborhood == nil {
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:               store.ID,
			EventType:             "SEARCH",
			SearchType:            "LOCATION",
			SearchedValue:         coordStr,
			SearchedNeighborhood:  &baseBairro.OfficialName,
			MatchedNeighborhoodID: &baseBairro.ID,
			DeliveryAvailable:     false,
			DeliveryPrice:         0,
			ResponseTimeMs:        ms,
			SessionID:             sessionID,
			UserAgent:             userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          baseBairro.OfficialName,
			Street:          street,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	deliveryTime := "24h"
	if neighborhood.DeliveryTime != nil && *neighborhood.DeliveryTime != "" {
		deliveryTime = *neighborhood.DeliveryTime
	}

	_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
		StoreID:               store.ID,
		EventType:             "SEARCH",
		SearchType:            "LOCATION",
		SearchedValue:         coordStr,
		SearchedNeighborhood:  &baseBairro.OfficialName,
		MatchedNeighborhoodID: &baseBairro.ID,
		DeliveryAvailable:     neighborhood.DeliveryEnabled,
		DeliveryPrice:         neighborhood.Fee,
		ResponseTimeMs:        ms,
		SessionID:             sessionID,
		UserAgent:             userAgent,
	}, clientIP)

	return &dto.LookupResultResponse{
		Success:               true,
		Bairro:                baseBairro.OfficialName,
		Street:                street,
		DeliveryEnabled:       neighborhood.DeliveryEnabled,
		Fee:                   &neighborhood.Fee,
		DeliveryTime:          deliveryTime,
		MinimumOrder:          neighborhood.MinimumOrder,
		FreeDeliveryThreshold: neighborhood.FreeDeliveryThreshold,
		Notes:                 neighborhood.Notes,
		StoreAddress:          store.Address,
		StoreWhatsapp:         store.Whatsapp,
		PickupEnabled:         store.PickupEnabled,
	}, nil
}

func (s *shippingService) LookupSelectedAddress(ctx context.Context, input dto.LookupSelectedAddressInput, clientIP, userAgent string) (*dto.LookupResultResponse, error) {
	startTime := time.Now()

	store, err := s.storeRepo.FindBySlug(ctx, input.StoreSlug)
	if err != nil {
		return &dto.LookupResultResponse{Success: false, Error: "Loja não encontrada."}, nil
	}

	sessionID := uuid.New().String()
	matchedBairroName := input.BairroName

	if matchedBairroName == "" {
		rev, err := s.nominatimClient.ReverseGeocode(ctx, input.Lat, input.Lon)
		if err == nil && rev != nil {
			matchedBairroName = external.ExtractBairro(*rev)
		}
	}

	if matchedBairroName == "" {
		return &dto.LookupResultResponse{
			Success:         true,
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
			Error:           "Bairro não identificado para esta localização.",
		}, nil
	}

	normalizedBairro := external.NormalizeName(matchedBairroName)
	baseBairro, err := s.baseNeighborhoodRepo.FindByName(ctx, normalizedBairro)
	ms := int(time.Since(startTime).Milliseconds())

	if err != nil || baseBairro == nil {
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:              store.ID,
			EventType:            "SEARCH",
			SearchType:           "ADDRESS",
			SearchedValue:        input.Address,
			SearchedNeighborhood: &matchedBairroName,
			DeliveryAvailable:    false,
			DeliveryPrice:        0,
			ResponseTimeMs:       ms,
			SessionID:            sessionID,
			UserAgent:            userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          matchedBairroName,
			Street:          "",
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	neighborhood, err := s.neighborhoodRepo.FindStoreNeighborhoodByBaseID(ctx, store.ID, baseBairro.ID)
	if err != nil || neighborhood == nil {
		_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
			StoreID:               store.ID,
			EventType:             "SEARCH",
			SearchType:            "ADDRESS",
			SearchedValue:         input.Address,
			SearchedNeighborhood:  &baseBairro.OfficialName,
			MatchedNeighborhoodID: &baseBairro.ID,
			DeliveryAvailable:     false,
			DeliveryPrice:         0,
			ResponseTimeMs:        ms,
			SessionID:             sessionID,
			UserAgent:             userAgent,
		}, clientIP)

		return &dto.LookupResultResponse{
			Success:         true,
			Bairro:          baseBairro.OfficialName,
			Street:          "",
			DeliveryEnabled: false,
			StoreAddress:    store.Address,
			StoreWhatsapp:   store.Whatsapp,
			PickupEnabled:   store.PickupEnabled,
		}, nil
	}

	deliveryTime := "2h"
	if neighborhood.DeliveryTime != nil && *neighborhood.DeliveryTime != "" {
		deliveryTime = *neighborhood.DeliveryTime
	} else if store.DeliveryTimeDefault != "" {
		deliveryTime = store.DeliveryTimeDefault
	}

	_ = s.telemetryService.LogSearchEvent(ctx, dto.CreateSearchEventInput{
		StoreID:               store.ID,
		EventType:             "SEARCH",
		SearchType:            "ADDRESS",
		SearchedValue:         input.Address,
		SearchedNeighborhood:  &baseBairro.OfficialName,
		MatchedNeighborhoodID: &baseBairro.ID,
		DeliveryAvailable:     neighborhood.DeliveryEnabled,
		DeliveryPrice:         neighborhood.Fee,
		ResponseTimeMs:        ms,
		SessionID:             sessionID,
		UserAgent:             userAgent,
	}, clientIP)

	return &dto.LookupResultResponse{
		Success:               true,
		Bairro:                baseBairro.OfficialName,
		Street:                "",
		DeliveryEnabled:       neighborhood.DeliveryEnabled,
		Fee:                   &neighborhood.Fee,
		DeliveryTime:          deliveryTime,
		MinimumOrder:          neighborhood.MinimumOrder,
		FreeDeliveryThreshold: neighborhood.FreeDeliveryThreshold,
		Notes:                 neighborhood.Notes,
		StoreAddress:          store.Address,
		StoreWhatsapp:         store.Whatsapp,
		PickupEnabled:         store.PickupEnabled,
	}, nil
}

func (s *shippingService) GetAddressSuggestions(ctx context.Context, query string) ([]dto.AddressSuggestionItem, error) {
	if len(strings.TrimSpace(query)) < 3 {
		return []dto.AddressSuggestionItem{}, nil
	}

	results, err := s.nominatimClient.SearchAddressSuggestions(ctx, query)
	if err != nil {
		return []dto.AddressSuggestionItem{}, nil
	}

	items := make([]dto.AddressSuggestionItem, 0, len(results))
	for _, item := range results {
		lat, _ := strconv.ParseFloat(item.Lat, 64)
		lon, _ := strconv.ParseFloat(item.Lon, 64)

		bairro := external.ExtractBairro(item)
		road := item.Address.Road

		items = append(items, dto.AddressSuggestionItem{
			DisplayName: item.DisplayName,
			Lat:         lat,
			Lon:         lon,
			Bairro:      bairro,
			Road:        road,
		})
	}

	return items, nil
}
