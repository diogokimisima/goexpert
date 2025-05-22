package location

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ViaCEPService struct{}

type Address struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

// GetAddress fetches address information from a zip code
func (s *ViaCEPService) GetAddress(zipCode string) (*Address, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to ViaCEP API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP API returned status code %d", resp.StatusCode)
	}

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, fmt.Errorf("error decoding ViaCEP API response: %w", err)
	}

	// Check if the response contains an error
	if address.CEP == "" {
		return nil, fmt.Errorf("zip code not found")
	}

	return &address, nil
}

// For backwards compatibility
func GetAddress(zipCode string) (*Address, error) {
	service := &ViaCEPService{}
	return service.GetAddress(zipCode)
}
