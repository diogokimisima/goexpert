package main

import (
	"fmt"

	"github.com/diogokimisima/goexpert/packaging/1/math"
)

func main() {
	m := math.NewMath(1, 2)

	fmt.Println(m.Add())
}
