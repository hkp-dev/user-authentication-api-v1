package services

import (
	"app/database"
	"app/models"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ValidateCategory(category models.Category) []string {
	var errors []string
	validate := validator.New()
	err := validate.Struct(category)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field %s: %s", e.Field(), e.Tag()))
		}
	}
	return errors
}
func GetCategoryByID(categoryID primitive.ObjectID) (models.Category, error) {
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	var category models.Category
	err := database.CategoryCollection.FindOne(context.Background(), bson.M{"_id": categoryID}).Decode(&category)
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}
func CheckCategoryExists(title string) (bool, error) {
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	var category models.Category
	err := database.CategoryCollection.FindOne(context.Background(), bson.M{"title": title}).Decode(&category)
	if err != nil {
		return false, err
	}
	return true, nil
}
func CreateCategory(category models.Category) (*mongo.InsertOneResult, error) {
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	category.ID = primitive.NewObjectID()
	result, err := database.CategoryCollection.InsertOne(context.Background(), category)
	if err != nil {
		return result, err
	}
	return result, nil
}
func UpdateFieldToArray(productID primitive.ObjectID, fieldName string, value interface{}) error {
	if productID.IsZero() {
		return fmt.Errorf("product ID is empty")
	}
	if fieldName == "" {
		return fmt.Errorf("field name is empty")
	}
	database.CategoryCollection = database.GetCollection("testDB", "categories")

	_, err := database.CategoryCollection.UpdateMany(
		context.Background(),
		bson.M{"product_ids": bson.M{"$not": bson.M{"$type": "array"}}},
		bson.M{
			"$set": bson.M{
				"product_ids": []primitive.ObjectID{}}},
	)
	if err != nil {
		return err
	}
	return nil
}
func DeleteProductInCategory(productID primitive.ObjectID, categoryID primitive.ObjectID) error {
	if productID.IsZero() {
		return fmt.Errorf("product ID is empty")
	}
	if categoryID.IsZero() {
		return fmt.Errorf("category ID is empty")
	}
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	_, err := database.CategoryCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": categoryID},
		bson.M{"$pull": bson.M{"product_ids": productID}},
	)
	if err != nil {
		return fmt.Errorf("error pulling product ID from category: %v", err)
	}
	return nil
}
func AddProductToCategory(productID primitive.ObjectID, categoryIDs []primitive.ObjectID) error {
	if categoryIDs == nil {
		return fmt.Errorf("category IDs are empty")
	}
	if productID.IsZero() {
		return fmt.Errorf("product ID is empty")
	}
	err := UpdateFieldToArray(productID, "product_ids", categoryIDs)
	if err != nil {
		return err
	}
	database.CategoryCollection = database.GetCollection("testDB", "categories")

	_, err = database.CategoryCollection.UpdateMany(
		context.Background(),
		bson.M{"_id": bson.M{"$in": categoryIDs}},
		bson.M{
			"$push": bson.M{"product_ids": productID},
		},
	)
	if err != nil {
		return fmt.Errorf("error pushing product ID to categories: %v", err)
	}
	return nil
}
func GetCategoriesByTitle(title string) (models.Category, error) {
	if title == "" {
		return models.Category{}, fmt.Errorf("category title is empty")
	}
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	var category models.Category
	err := database.CategoryCollection.FindOne(context.Background(), bson.M{"title": title}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return category, fmt.Errorf("category with title %s not found", title)
		}
		return category, err
	}
	return category, nil
}
func GetProductsByCategoryTitle(title string) ([]models.Product, error) {
	if title == "" {
		return nil, fmt.Errorf("category title is empty")
	}
	database.ProductCollection = database.GetCollection("testDB", "products")
	var products []models.Product
	cursor, err := database.ProductCollection.Find(context.Background(), bson.M{"category": title})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var product models.Product
		err := cursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
func GetAllCategory() ([]models.Category, error) {
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	var categories []models.Category
	categoriesCursor, err := database.CategoryCollection.Find(
		context.Background(),
		bson.M{})
	if err != nil {
		return nil, err
	}
	for categoriesCursor.Next(context.Background()) {
		var category models.Category
		err := categoriesCursor.Decode(category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
func DeleteCategory(categoryID primitive.ObjectID) error {
	if categoryID.IsZero() {
		return fmt.Errorf("category ID is empty")
	}

	database.CategoryCollection = database.GetCollection("testDB", "categories")
	_, err := database.CategoryCollection.Find(context.Background(), bson.M{"_id": categoryID})
	if err == nil {
		return fmt.Errorf("category with ID %s not found", categoryID)
	}
	_, err = database.CategoryCollection.DeleteOne(context.Background(), bson.M{"_id": categoryID})
	if err != nil {
		return fmt.Errorf("error deleting category: %v", err)
	}
	database.ProductCollection = database.GetCollection("testDB", "products")
	_, err = database.ProductCollection.DeleteMany(context.Background(), bson.M{"category_id": categoryID})
	if err != nil {
		return fmt.Errorf("error deleting products in category: %v", err)
	}
	return nil
}
