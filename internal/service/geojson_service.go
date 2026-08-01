package service

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

type GeoJSONService interface {
	GetBairrosFortaleza(ctx context.Context) ([]byte, error)
}

type geoJSONService struct {
	filePath string
	cache    []byte
	mu       sync.RWMutex
}

func NewGeoJSONService(customPath string) GeoJSONService {
	if customPath == "" {
		customPath = "public/bairros-fortaleza.geojson"
	}
	return &geoJSONService{filePath: customPath}
}

func (s *geoJSONService) GetBairrosFortaleza(ctx context.Context) ([]byte, error) {
	s.mu.RLock()
	if len(s.cache) > 0 {
		defer s.mu.RUnlock()
		return s.cache, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Checar cache novamente
	if len(s.cache) > 0 {
		return s.cache, nil
	}

	pathsToTry := []string{
		s.filePath,
		filepath.Join(".", "public", "bairros-fortaleza.geojson"),
		filepath.Join("..", "project-archives", "public", "bairros-fortaleza.geojson"),
		"/Users/jamesmachome/Desktop/fretefortal/project-archives/public/bairros-fortaleza.geojson",
	}

	var data []byte
	var lastErr error

	for _, path := range pathsToTry {
		b, err := os.ReadFile(path)
		if err == nil {
			data = b
			break
		}
		lastErr = err
	}

	if len(data) == 0 {
		return nil, lastErr
	}

	s.cache = data
	return data, nil
}
