package main

import (
	"app/database"
	"app/routes"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://localhost:27017"
	err := database.ConnectMongoDB(uri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	err = database.CreateCollection("testDB", "users", &options.CreateCollectionOptions{})
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	routes.AddRoutes()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
