package services

import (
	"app/database"
	"app/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCart(user models.User) error {
	database.CartCollection = database.GetCollection("testDB", "carts")
	var cart models.Cart
	cart.ID = primitive.NewObjectID()
	cart.UserID = user.ID
	cart.Products = []models.Product{}
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()
	_, err := database.CartCollection.InsertOne(context.Background(), cart)
	if err != nil {
		return err
	}
	return nil
}	
