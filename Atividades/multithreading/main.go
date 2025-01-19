package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro,omitempty"`
	Bairro      string `json:"bairro,omitempty"`
	Cidade      string `json:"cidade,omitempty"`
	Estado      string `json:"estado,omitempty"`
	Complemento string `json:"complemento,omitempty"`
	API         string `json:"api"`
}

type result struct {
	Address Address
	Err     error
}

func fetch(url, api string, parseFunc func(map[string]interface{}) Address, ch chan<- result) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- result{Err: fmt.Errorf("%s: %v", api, err)}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- result{Err: fmt.Errorf("%s: %v", api, err)}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ch <- result{Err: fmt.Errorf("%s: %v", api, err)}
		return
	}

	ch <- result{Address: parseFunc(data)}
}

func main() {
	cep := "16078095"
	ch := make(chan result, 2)

	// Thread 1 Requisição para BrasilAPI
	go fetch(
		fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep),
		"BrasilAPI",
		func(data map[string]interface{}) Address {
			return Address{
				Cep:        data["cep"].(string),
				Logradouro: data["street"].(string),
				Bairro:     data["neighborhood"].(string),
				Cidade:     data["city"].(string),
				Estado:     data["state"].(string),
				API:        "BrasilAPI",
			}
		},
		ch,
	)

	// Thread 2 Requisição para ViaCEP
	go fetch(
		fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep),
		"ViaCEP",
		func(data map[string]interface{}) Address {
			return Address{
				Cep:         data["cep"].(string),
				Logradouro:  data["logradouro"].(string),
				Bairro:      data["bairro"].(string),
				Cidade:      data["localidade"].(string),
				Estado:      data["uf"].(string),
				Complemento: data["complemento"].(string),
				API:         "ViaCEP",
			}
		},
		ch,
	)

	select {
	case res := <-ch:
		if res.Err != nil {
			fmt.Println("Erro:", res.Err)
		} else {
			fmt.Printf("Resultado mais rápido:\nAPI: %s\nEndereço: %+v\n", res.Address.API, res.Address)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: Timeout nas requisições.")
	}
}
