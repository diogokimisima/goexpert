package main

import (
	"fmt"
)

func main() {
	fmt.Println(sum(1, 3, 45, 6, 123, 6, 23, 3, 41))
	fmt.Println(mult("a", "b", "c", "d", "e"))
}

func sum(numeros ...int) int {
	total := 0
	for _, numero := range numeros {
		total += numero
	}
	return total
}

func mult(letras ...string) string {
	string := ""
	for _, letra := range letras {
		string += letra
	}
	return string
}
