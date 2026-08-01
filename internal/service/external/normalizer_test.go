package external

import (
	"testing"
)

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Aldeota", "aldeota"},
		{"Água Fria!", "agua fria"},
		{"São João do Tauape", "sao joao do tauape"},
		{"  Mesa - Redonda  ", "mesa - redonda"},
		{"Meireles @ Fortaleza", "meireles fortaleza"},
	}

	for _, tt := range tests {
		result := NormalizeName(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeName(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractBairro(t *testing.T) {
	item := NominatimSearchResult{
		Address: NominatimAddress{
			Suburb: "Aldeota",
		},
	}

	bairro := ExtractBairro(item)
	if bairro != "Aldeota" {
		t.Errorf("ExtractBairro() = %q; want %q", bairro, "Aldeota")
	}

	itemRegional := NominatimSearchResult{
		Address: NominatimAddress{
			Suburb:        "SER II",
			Neighbourhood: "Meireles",
		},
	}

	bairroReg := ExtractBairro(itemRegional)
	if bairroReg != "Meireles" {
		t.Errorf("ExtractBairro() = %q; want %q", bairroReg, "Meireles")
	}
}
