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

	database.ProductCollection = database.GetCollection("testDB", "products")
	err := database.ProductCollection.FindOne(
		context.Background(),
		bson.M{"title": product.Title},
	).Decode(&models.Product{})
	if err == nil {
		return fmt.Errorf("product with title %s already exists", product.Title)
	}
	_, err = database.ProductCollection.InsertOne(context.Background(), product)
	if err != nil {
		return fmt.Errorf("error inserting product: %v", err)
	}

	err = AddProductToCategory(product.ID, product.CategoryID)
	if err != nil {
		return fmt.Errorf("error adding product to category: %v", err)
	}

	return nil
}
func GetAllProductByTitle(title string) ([]models.Product, error) {
	database.ProductCollection = database.GetCollection("testDB", "products")
	var products []models.Product
	productCursor, err := database.ProductCollection.Find(context.Background(), bson.M{"title": title})
	if err != nil {
		return nil, err
	}
	for productCursor.Next(context.Background()) {
		var product models.Product
		err := productCursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
func GetAllProduct() ([]models.Product, error) {
	database.ProductCollection = database.GetCollection("testDB", "products")
	var products []models.Product
	productCursor, err := database.ProductCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	for productCursor.Next(context.Background()) {
		var product models.Product
		err := productCursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
func UpdatePriceProduct(productID primitive.ObjectID, price float64) error {
	database.ProductCollection = database.GetCollection("testDB", "products")
	_, err := database.ProductCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": productID},
		bson.M{"$set": bson.M{"price": price, "updated_at": time.Now()}},
	)
	if err != nil {
		return err
	}
	return nil
}
func DeleteProduct(productID primitive.ObjectID) error {
	database.ProductCollection = database.GetCollection("testDB", "products")
	_, err := database.ProductCollection.DeleteOne(context.Background(), bson.M{"_id": productID})
	if err != nil {
		return fmt.Errorf("error deleting product: %v", err)
	}
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	_, err = database.CategoryCollection.UpdateMany(
		context.Background(),
		bson.M{"product_ids": productID},
		bson.M{"$pull": bson.M{"product_ids": productID}},
	)
	if err != nil {
		return fmt.Errorf("error deleting product from categories: %v", err)
	}
	return nil
}
