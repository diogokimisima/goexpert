package main

import (
	"net/http"

	"github.com/diogokimisima/goexpert/7-APIS/configs"
	_ "github.com/diogokimisima/goexpert/7-APIS/docs"
	"github.com/diogokimisima/goexpert/7-APIS/internal/entity"
	"github.com/diogokimisima/goexpert/7-APIS/internal/infra/database"
	"github.com/diogokimisima/goexpert/7-APIS/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwager "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title 			Go Expert API Example
// @version 	   	1.0
// @description 	Product API with authentication
// @termsOfService 	http://swagger.io/terms/

// @contact.name 	Diogo Kimisima
// @contact.url 	http://www.diogokimisima.com.br
// @contact.email 	dkimisima@gmail.com

// @license.name 	Full Cycle License
// @license.url 	http://www.fullcycle.com.br

// @host 			localhost:8000
// @BasePath 		/
// @securityDefinitions.apikey ApiKeyAuth
// @in 				header
// @name 			Authorization
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
	userHandler := handlers.NewUserHandler(userDb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Get("/", productHandler.ListProducts)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJwt)

	r.Get("/docs/*", httpSwager.Handler(httpSwager.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
