package main

import "fmt"

// func main() {
// 	total := func() int {
// 		return sum(1, 3, 45, 6, 123, 6, 23, 3, 41) * 2
// 	}()

// 	fmt.Println(total)
// }

// func sum(numeros ...int) int {
// 	total := 0
// 	for _, numero := range numeros {
// 		total += numero
// 	}
// 	return total
// }

func main() {
	somador := sum()
	fmt.Println(somador(1))
	fmt.Println(somador(2))
	fmt.Println(somador(3))
}

func sum() func(int) int {
	total := 0
	return func(num int) int {
		total += num
		return total
	}
}
