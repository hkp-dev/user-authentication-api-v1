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

// func GetAllProducts(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet { // Use GET instead of POST
// 		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	products, err := services.GetAllProduct()
// 	if err != nil {
// 		log.Printf("Error retrieving products: %v", err)
// 		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	if len(products) == 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"message": "No products found",
// 		})
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"message":  "All products retrieved successfully",
// 		"products": products,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	cateID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = services.DeleteCategory(cateID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete category"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Category with ID " + input.ID + " deleted successfully",
	})
}
func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var cates []models.Category
	cates, err := services.GetAllCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		`message`: `Get Product Successfully`,
		`cates`:   cates,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
