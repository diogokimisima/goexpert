package main

import "fmt"

func main() {
	defer fmt.Println("Executando o defer")
	fmt.Println("Primeira linha")
	fmt.Println("Segunda linha")
	fmt.Println("Terceira linha")
}
