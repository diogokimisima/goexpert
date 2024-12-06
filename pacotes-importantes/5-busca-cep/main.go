package main

import (
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		println(url)
	}
}
