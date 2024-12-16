package controllers

import (
	"app/models"
	"app/services"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	createCategoryError := services.ValidateCategory(category)
	if len(createCategoryError) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": createCategoryError,
		})
		return
	}
	result, err := services.CreateCategory(category)
	if err != nil {
		http.Error(w, `{"error": "Failed to create category"}`, http.StatusInternalServerError)
		return
	}
	category.ID = result.InsertedID.(primitive.ObjectID)
	response := map[string]interface{}{
		"message": "Category created successfully",
		"id":      category.ID.Hex(),
		"title":   category.Title,
		"desc":    category.Description,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func GetCategories(w http.ResponseWriter, r *http.Request) {}

func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {}
