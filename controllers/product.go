package controllers

import (
	"app/models"
	"app/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	createProductError := services.ValidateProduct(product)
	if len(createProductError) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": createProductError,
		})
		return
	}
	var categories models.Category
	if err := services.CreateProduct(product, categories); err != nil {
		http.Error(w, `{"error": "Failed to create product"}`, http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	response := map[string]string{
		"message": "Product created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetProducts(w http.ResponseWriter, r *http.Request)        {}
func DeleteProduct(w http.ResponseWriter, r *http.Request)      {}
func UpdatePriceProduct(w http.ResponseWriter, r *http.Request) {}

// func CreateProduct(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var product models.Product
// 	err := json.NewDecoder(r.Body).Decode(&product)
// 	if err != nil {
// 		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Validate the product data
// 	validationErrors := services.ValidateProduct(product)
// 	if len(validationErrors) > 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"error": validationErrors,
// 		})
// 		return
// 	}

// 	product.CreatedAt = time.Now()
// 	product.UpdatedAt = time.Now()

// 	// Ensure database collection is initialized
// 	database.ProductCollection = database.GetCollection("testDB", "products")

// 	// Insert the product into the database
// 	product.ID = primitive.NewObjectID() // generate new product ID
// 	result, err := database.ProductCollection.InsertOne(context.Background(), product)
// 	if err != nil {
// 		http.Error(w, `{"error": "Failed to create product"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	// Ensure CategoryID is not nil or empty before attempting to update
// 	if len(product.CategoryID) == 0 {
// 		http.Error(w, `{"error": "CategoryID cannot be empty"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Update categories with the new product's ID
// 	database.CategoryCollection = database.GetCollection("testDB", "categories")
// 	if database.CategoryCollection == nil {
// 		http.Error(w, `{"error": "Failed to get category collection"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	_, err = database.CategoryCollection.UpdateMany(
// 		context.Background(),
// 		bson.M{"_id": bson.M{"$in": product.CategoryID}},
// 		bson.M{
// 			"$set": bson.M{
// 				"product_ids": bson.M{"$ifNull": []primitive.ObjectID{}},
// 			},
// 			"$push": bson.M{
// 				"product_ids": product.ID,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Error updating categories: %v", err)
// 		return
// 	}

// 	// Prepare response
// 	product.ID = result.InsertedID.(primitive.ObjectID)
// 	response := map[string]interface{}{
// 		"message": "Product created successfully",
// 		"id":      product.ID.Hex(),
// 		"title":   product.Title,
// 		"desc":    product.Description,
// 		"price":   product.Price,
// 		"stock":   product.Stock,
// 	}

// 	// Send response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }
