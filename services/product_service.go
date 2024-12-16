package services

import (
	"app/database"
	"app/models"
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateProduct(product models.Product) []string {
	var errors []string
	validate := validator.New()
	err := validate.Struct(product)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field %s: %s", e.Field(), e.Tag()))
		}
	}
	return errors
}
func CreateProduct(product models.Product, category models.Category) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	_, err := database.CategoryCollection.UpdateMany(
		context.Background(),
		bson.M{"_id": bson.M{"$in": product.CategoryID}},
		bson.M{
			"$set": bson.M{
				"product_ids": bson.M{"$ifNull": []primitive.ObjectID{}},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error updating category with product ID: %v", err)
	}
	database.ProductCollection = database.GetCollection("testDB", "products")
	_, err = database.ProductCollection.InsertOne(context.Background(), product)
	if err != nil {
		return err
	}
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	_, err = database.CategoryCollection.UpdateMany(
		context.Background(),
		bson.M{"_id": product.CategoryID},
		bson.M{
			"$push": bson.M{"product_ids": product.ID},
		})
	if err != nil {
		return fmt.Errorf("error updating category with product ID: %v", err)
	}
	return nil
}
func GetProduct(id primitive.ObjectID) (models.Product, error) {
	database.ProductCollection = database.GetCollection("testDB", "products")
	var product models.Product
	err := database.ProductCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	return product, err
}
func UpdateProduct(product models.Product) error {
	product.UpdatedAt = time.Now()
	database.ProductCollection = database.GetCollection("testDB", "products")
	_, err := database.ProductCollection.UpdateOne(context.Background(), bson.M{"_id": product.ID}, bson.M{"$set": product})
	if err != nil {
		return err
	}
	return nil
}
