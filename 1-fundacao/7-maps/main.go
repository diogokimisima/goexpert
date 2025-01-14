package main

import "fmt"

func main() {

	salarios := map[string]int{"Wesley": 1000, "João": 2000, "Maria": 3000}
	// fmt.Println(salarios["Wesley"])
	delete(salarios, "Wesley")
	// fmt.Println(salarios)
	salarios["Wes"] = 5000
	// fmt.Println(salarios)

	// sal := make(map[string]int)
	// sal1 := map[string]int{}
	// sal["Wesley"] = 100
	// sal1["Wesley"] = 101
	// fmt.Println(sal)
	// fmt.Println(sal1)

	for nome, salario := range salarios {
		fmt.Printf("O salário de %s é R$%d.\n", nome, salario)
	}

	for _, salario := range salarios {
		fmt.Printf("O salário é R$%d.\n", salario)
	}
}
