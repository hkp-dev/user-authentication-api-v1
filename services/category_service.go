package services

import (
	"app/database"
	"app/models"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
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
func CreateCategory(category models.Category) (*mongo.InsertOneResult, error) {
	database.CategoryCollection = database.GetCollection("testDB", "categories")
	category.ID = primitive.NewObjectID()
	result, err := database.CategoryCollection.InsertOne(context.Background(), category)
	if err != nil {
		return result, err
	}
	return result, nil
}
