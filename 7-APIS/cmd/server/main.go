package main

import (
	"net/http"

	"github.com/diogokimisima/goexpert/7-APIS/configs"
	"github.com/diogokimisima/goexpert/7-APIS/internal/entity"
	"github.com/diogokimisima/goexpert/7-APIS/internal/infra/database"
	"github.com/diogokimisima/goexpert/7-APIS/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDb := database.NewProduct(db)
	userDb := database.NewUser(db)

	productHandler := handlers.NewProductHandler(productDb)
	userHandler := handlers.NewUserHandler(userDb, configs.TokenAuth, configs.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)
	r.Get("/products", productHandler.ListProducts)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJwt)

	http.ListenAndServe(":8000", r)
}
