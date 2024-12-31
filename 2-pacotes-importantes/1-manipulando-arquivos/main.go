package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Criação e escrita no arquivo
	f, err := os.Create("arquivo.txt")
	if err != nil {
		panic(err)
	}

	tamanho, err := f.Write([]byte("Escrevendo dados no arquivo"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)

	// Fecha o arquivo após a escrita
	err = f.Close()
	if err != nil {
		panic(err)
	}

	// Leitura completa do arquivo
	arquivo, err := os.ReadFile("arquivo.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(arquivo))

	// Leitura em partes com buffer
	arquivo2, err := os.Open("arquivo.txt")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(arquivo2)
	buffer := make([]byte, 10)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println(string(buffer[:n]))
	}

	// Fecha o arquivo após a leitura
	err = arquivo2.Close()
	if err != nil {
		panic(err)
	}

	// Remoção do arquivo
	err = os.Remove("arquivo.txt")
	if err != nil {
		panic(err)
	}
}
