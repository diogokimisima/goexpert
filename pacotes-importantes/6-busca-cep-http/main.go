package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	http.HandleFunc("/", BuscaCepHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	// Processa requisições GET.
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Verifica se o parâmetro cep foi fornecido na query string.
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Chama a função BuscaCEP para buscar o CEP
	cep, error := BuscaCEP(cepParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retorna o resultado da busca em formato JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// result, err := json.Marshal(cep)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(result)

	json.NewEncoder(w).Encode(cep)
	return
}

func BuscaCEP(cep string) (*ViaCEP, error) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json")
	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
