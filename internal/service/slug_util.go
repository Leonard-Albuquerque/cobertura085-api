package service

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converte uma string em formato de slug URL amigável.
// Exemplo: "Lojão do menino" -> "lojao-do-menino"
func Slugify(s string) string {
	// 1. Remover acentos e diacríticos (normalização NFD + remoção de Mn runes)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		result = s
	}

	// 2. Converter para minúsculas
	result = strings.ToLower(result)

	// 3. Substituir caracteres não alfanuméricos por hífen
	result = nonAlphanumericRegex.ReplaceAllString(result, "-")

	// 4. Remover hífens do início e do fim
	result = strings.Trim(result, "-")

	return result
}
