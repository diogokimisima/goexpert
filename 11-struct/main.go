package main

import "fmt"

type Cliente struct {
	Nome  string
	Idade int
	Ativo bool
}

func main() {
	wesley := Cliente{
		Nome:  "Wesley",
		Idade: 25,
		Ativo: true,
	}

	wesley.Ativo = false

	fmt.Printf("Nome: %s, idade: %d, ativo: %t", wesley.Nome, wesley.Idade, wesley.Ativo)
}
