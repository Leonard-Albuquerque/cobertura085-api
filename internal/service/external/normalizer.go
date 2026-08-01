package external

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9\s-]`)
	multipleSpacesRegex  = regexp.MustCompile(`\s+`)
)

// NormalizeName remove acentos, caracteres especiais e converte a string para caixa baixa normalizada.
func NormalizeName(str string) string {
	// 1. Converter para caixa baixa
	lower := strings.ToLower(str)

	// 2. Remover acentos (NFD decomposition -> remove Mn mark category)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, lower)
	if err != nil {
		result = lower
	}

	// 3. Remover caracteres especiais exceto letras, números, espaços e hífen
	clean := nonAlphanumericRegex.ReplaceAllString(result, "")

	// 4. Normalizar múltiplos espaços para espaço único e dar Trim
	clean = multipleSpacesRegex.ReplaceAllString(clean, " ")
	return strings.TrimSpace(clean)
}

// NominatimAddress representa o objeto "address" retornado pela API do OpenStreetMap Nominatim.
type NominatimAddress struct {
	Road         string `json:"road,omitempty"`
	Suburb       string `json:"suburb,omitempty"`
	Neighbourhood string `json:"neighbourhood,omitempty"`
	Quarter      string `json:"quarter,omitempty"`
	CityDistrict string `json:"city_district,omitempty"`
	CityBlock    string `json:"city_block,omitempty"`
	Residential  string `json:"residential,omitempty"`
	City         string `json:"city,omitempty"`
	Town         string `json:"town,omitempty"`
	Municipality string `json:"municipality,omitempty"`
	Postcode     string `json:"postcode,omitempty"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country,omitempty"`
}

// NominatimSearchResult representa o item retornado pelas APIs de busca e reverse geocode do Nominatim.
type NominatimSearchResult struct {
	PlaceID     int64            `json:"place_id"`
	Lat         string           `json:"lat"`
	Lon         string           `json:"lon"`
	DisplayName string           `json:"display_name"`
	Type        string           `json:"type,omitempty"`
	Name        string           `json:"name,omitempty"`
	Address     NominatimAddress `json:"address"`
}

// ExtractBairro extrai o nome do bairro do resultado do Nominatim descartando secretarias regionais (SER).
func ExtractBairro(item NominatimSearchResult) string {
	addr := item.Address

	candidates := []string{
		addr.Suburb,
		addr.Neighbourhood,
		addr.Quarter,
		addr.CityDistrict,
		addr.CityBlock,
		addr.Residential,
	}

	// Filtrar secretarias regionais (SER I, SER II, etc.)
	filtered := make([]string, 0, len(candidates))
	for _, c := range candidates {
		trimmed := strings.TrimSpace(c)
		if trimmed == "" {
			continue
		}
		lc := strings.ToLower(trimmed)
		if strings.Contains(lc, "regional") || strings.Contains(lc, "ser ") || strings.HasPrefix(lc, "ser") || strings.Contains(lc, "secretaria executiva") {
			continue
		}
		filtered = append(filtered, trimmed)
	}

	if len(filtered) > 0 {
		return filtered[0]
	}

	// Fallback por tipo do item
	validTypes := map[string]bool{
		"suburb":        true,
		"neighbourhood": true,
		"quarter":       true,
		"city_district": true,
		"district":      true,
	}

	if validTypes[item.Type] {
		if item.Name != "" {
			return strings.TrimSpace(item.Name)
		}
		if item.DisplayName != "" {
			parts := strings.Split(item.DisplayName, ",")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}
	}

	return ""
}
