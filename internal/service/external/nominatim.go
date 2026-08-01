package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NominatimClient interface {
	SearchAddress(ctx context.Context, query string) ([]NominatimSearchResult, error)
	ReverseGeocode(ctx context.Context, lat, lon float64) (*NominatimSearchResult, error)
	SearchAddressSuggestions(ctx context.Context, query string) ([]NominatimSearchResult, error)
}

type nominatimClient struct {
	httpClient *http.Client
	userAgent  string
}

func NewNominatimClient() NominatimClient {
	return &nominatimClient{
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
		userAgent: "FreteFortalApp/1.0 (contact: admin@fretefortal.com)",
	}
}

func (c *nominatimClient) SearchAddress(ctx context.Context, rawQuery string) ([]NominatimSearchResult, error) {
	query := strings.TrimSpace(rawQuery)
	if !strings.Contains(strings.ToLower(query), "fortaleza") {
		query = fmt.Sprintf("%s, Fortaleza, Ceará, Brasil", query)
	}

	reqURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&addressdetails=1&limit=1&countrycodes=br",
		url.QueryEscape(query),
	)

	return c.fetchSearch(ctx, reqURL)
}

func (c *nominatimClient) SearchAddressSuggestions(ctx context.Context, rawQuery string) ([]NominatimSearchResult, error) {
	query := strings.TrimSpace(rawQuery)
	if len(query) < 3 {
		return nil, nil
	}

	if !strings.Contains(strings.ToLower(query), "fortaleza") {
		query = fmt.Sprintf("%s, Fortaleza, Ceará, Brasil", query)
	}

	reqURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&addressdetails=1&limit=5&countrycodes=br",
		url.QueryEscape(query),
	)

	return c.fetchSearch(ctx, reqURL)
}

func (c *nominatimClient) ReverseGeocode(ctx context.Context, lat, lon float64) (*NominatimSearchResult, error) {
	reqURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?lat=%.6f&lon=%.6f&format=json&addressdetails=1",
		lat, lon,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Nominatim reverse retornou status HTTP %d", resp.StatusCode)
	}

	var result NominatimSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *nominatimClient) fetchSearch(ctx context.Context, reqURL string) ([]NominatimSearchResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Nominatim search retornou status HTTP %d", resp.StatusCode)
	}

	var results []NominatimSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}
