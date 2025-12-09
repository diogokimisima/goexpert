package main

import (
	"context"
	"log"

	"github.com/diogokimisima/fullcycle-auction/configuration/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	databaseClient, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return
	}

	log.Println("Successfully connected to MongoDB:", databaseClient.Name())
}
