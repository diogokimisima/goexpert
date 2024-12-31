package main

func main() {
	// for i := 0; i < 10; i++ {
	// 	println(i)
	// }

	numero := []string{"um", "dois", "tres"}
	for k, v := range numero {
		println(k, v)
	}

	numero2 := []string{"um", "dois", "tres"}
	for _, v := range numero2 {
		println(v)
	}
}
