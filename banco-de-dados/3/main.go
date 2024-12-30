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
	ID         int `gorm:"column:id;primaryKey"`
	Name       string
	Price      float64
	CategoryID int
	Category   Category
	gorm.Model
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{})

	category := Category{Name: "Eletrônicos"}
	db.Create(&category)

	// db.Create(&Product{
	// 	Name:       "Mouse",
	// 	Price:      1000.00,
	// 	CategoryID: category.ID,
	// })

	var products []Product
	db.Preload("Category").Find(&products)
	for _, product := range products {
		println(product.Name, product.Category.Name)
	}
}
