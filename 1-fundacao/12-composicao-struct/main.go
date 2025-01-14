package main

import "fmt"

type Endereco struct {
	Logradouro string
	Numero     int
	Cidade     string
	Estado     string
}

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
	Endereco
}

// Alterar para usar ponteiro (c *Cliente)
func (c *Cliente) Desativar() {
	c.Ativo = false
	fmt.Printf("O cliente %s, foi desativado\n", c.Nome)
}

func main() {
	wesley := Cliente{
		Nome:  "Wesley",
		Idade: 25,
		Ativo: true,
	}

	// Passando o ponteiro de wesley (o & wesley)
	wesley.Desativar()

	// Verifique se o cliente foi desativado
	fmt.Printf("Status final do cliente %s: %v\n", wesley.Nome, wesley.Ativo)
}
