package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID       int `gorm:"column:id;primaryKey"`
	Name     string
	Products []Product
}

type Product struct {
	ID         int `gorm:"column:id;primaryKey"`
	Name       string
	Price      float64
	CategoryID int
	Category   Category
	gorm.Model
}

func main() {
	// Adicionando parseTime=true ao DSN
	dsn := "root:root@tcp(localhost:3306)/goexpert?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{})

	// category := Category{Name: "Eletrônicos"}
	// db.Create(&category)

	db.Create(&Product{
		Name:       "Mouse 2",
		Price:      1000.00,
		CategoryID: 1, // Ajustando para usar o ID do objeto criado
	})

	var products []Product
	db.Preload("Category").Find(&products)
	for _, product := range products {
		println(product.Name, product.Category.Name)
	}

	var categories []Category
	err = db.Model(&Category{}).Preload("Products").Find(&categories).Error
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		println(category.Name)
		for _, product := range category.Products {
			println("-", product.Name)
		}
	}

}
