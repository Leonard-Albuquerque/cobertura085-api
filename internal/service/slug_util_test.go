package service_test

import (
	"testing"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/service"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Lojão do menino", "lojao-do-menino"},
		{"  Super   Mercado 085!  ", "super-mercado-085"},
		{"Café & Cia.", "cafe-cia"},
		{"Açougue São José - Fortaleza", "acougue-sao-jose-fortaleza"},
		{"123 TESTE!!!", "123-teste"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := service.Slugify(tt.input)
			if got != tt.expected {
				t.Errorf("Slugify(%q) = %q; esperado %q", tt.input, got, tt.expected)
			}
		})
	}
}
