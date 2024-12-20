package main

import (
	"app/database"
	"app/routes"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Print("\\ \\        /         __ \\ _)      _) |         | \n" +
		" \\ \\  \\   / _ \\  _ \\ |   | |  _` | | __|  _` | | \n" +
		"  \\ \\  \\ /  __/  __/ |   | | (   | | |   (   | | \n" +
		"   \\_/\\_/ \\___|\\___|____/ _|\\__, |_|\\__|\\__,_|_| \n" +
		"                            |___/                \n")
	uri := "mongodb://localhost:27017"
	err := database.ConnectMongoDB(uri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	err = database.CreateCollection("testDB", "users", &options.CreateCollectionOptions{})
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	err = database.CreateCollection("testDB", "products", &options.CreateCollectionOptions{})
	if err != nil {
		log.Fatal("Failed to create collection:", err)
	}
	err = database.CreateCollection("testDB", "categories", &options.CreateCollectionOptions{})
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	err = database.CreateCollection("testDB", "carts", &options.CreateCollectionOptions{})
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	routes.AddRoutes()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
