package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient    *mongo.Client
	UserCollection *mongo.Collection
	once           sync.Once
)

func ContextTimeOut(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t)
}
func ConnectMongoDB(uri string) error {
	var ConnectError error
	once.Do(func() {
		ctx, cancel := ContextTimeOut(5 * time.Second)
		defer cancel()
		MongoClient, ConnectError = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if ConnectError != nil {
			ConnectError = fmt.Errorf("Failed to connect to MongoDB: %w\n", ConnectError)
			return
		}
		if ConnectError = MongoClient.Ping(ctx, nil); ConnectError != nil {
			ConnectError = fmt.Errorf("Failed to ping MongoDB: %w\n", ConnectError)
			return
		}
		fmt.Printf("Connected to database successfully!!!\n")
	})
	return ConnectError
}
func GetCollection(database, name string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoDB client is not initialized. Call Connect() first.")
	}
	return MongoClient.Database(database).Collection(name)
}
func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := ContextTimeOut(10 * time.Second)
		defer cancel()

		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect MongoDB: %v", err)
		} else {
			fmt.Println("Disconnected from MongoDB successfully!")
		}
	}
}
func CreateCollection(database, name string, opts *options.CreateCollectionOptions) error {
	ctx, cancel := ContextTimeOut(5 * time.Second)
	defer cancel()

	collections, err := MongoClient.Database(database).ListCollectionNames(ctx, bson.M{"name": name})
	if err != nil {
		return fmt.Errorf("Failed to list collections: %w", err)
	}

	if len(collections) > 0 {
		fmt.Printf("Collection '%s' already exists in database '%s'.\n", name, database)
		return nil
	}

	err = MongoClient.Database(database).CreateCollection(ctx, name, opts)
	if err != nil {
		return fmt.Errorf("Failed to create collection '%s': %w", name, err)
	}

	fmt.Printf("Collection '%s' created successfully in database '%s'.\n", name, database)
	return nil
}
