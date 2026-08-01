package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var (
	ErrViaCEPNotFound   = errors.New("CEP não encontrado")
	ErrViaCEPBadRequest = errors.New("formato de CEP inválido")
)

type ViaCEPResponse struct {
	CEP         string      `json:"cep"`
	Logradouro  string      `json:"logradouro"`
	Complemento string      `json:"complemento"`
	Bairro      string      `json:"bairro"`
	Localidade  string      `json:"localidade"`
	UF          string      `json:"uf"`
	IBGE        string      `json:"ibge"`
	DDD         string      `json:"ddd"`
	Erro        interface{} `json:"erro,omitempty"`
}

func (r *ViaCEPResponse) HasError() bool {
	if r.Erro == nil {
		return false
	}
	switch v := r.Erro.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1"
	default:
		return false
	}
}

type ViaCEPClient interface {
	Lookup(ctx context.Context, rawCEP string) (*ViaCEPResponse, error)
}

type viaCEPClient struct {
	httpClient *http.Client
}

func NewViaCEPClient() ViaCEPClient {
	return &viaCEPClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *viaCEPClient) Lookup(ctx context.Context, rawCEP string) (*ViaCEPResponse, error) {
	reg := regexp.MustCompile(`\D`)
	cleanCEP := reg.ReplaceAllString(rawCEP, "")

	if len(cleanCEP) != 8 {
		return nil, ErrViaCEPBadRequest
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cleanCEP)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP retornou status HTTP %d", resp.StatusCode)
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.HasError() {
		return nil, ErrViaCEPNotFound
	}

	return &data, nil
}
