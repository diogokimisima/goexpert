package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID   int `gorm:"column:id;primaryKey"`
	Name string
}

type Product struct {
	ID           int `gorm:"column:id;primaryKey"`
	Name         string
	Price        float64
	CategoryID   int
	Category     Category
	SerialNumber SerialNumber
	gorm.Model
}

type SerialNumber struct {
	ID        int `gorm:"column:id;primaryKey"`
	Number    string
	ProductID int
}

func main() {
	// Adicionando parseTime=true ao DSN
	dsn := "root:root@tcp(localhost:3306)/goexpert?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	category := Category{Name: "Eletrônicos"}
	db.Create(&category)

	db.Create(&Product{
		Name:       "Mouse",
		Price:      1000.00,
		CategoryID: 1, // Ajustando para usar o ID do objeto criado
	})

	db.Create(&SerialNumber{
		Number:    "123456",
		ProductID: 1,
	})

	var products []Product
	db.Preload("Category").Preload("SerialNumber").Find(&products)
	for _, product := range products {
		println(product.Name, product.Category.Name, product.SerialNumber.Number)
	}
}
