package main

import "github.com/diogokimisima/goexpert/7-APIS/configs"

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	println(config.DBDriver)
}
