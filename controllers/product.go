package controllers

import (
	"app/models"
	"app/services"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	category, err := services.GetCategoriesByTitle(product.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	product.CategoryID = append(product.CategoryID, category.ID)
	err = services.CreateProduct(product, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"message": "Product created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	productId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = services.DeleteProduct(productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product with ID " + input.ID + " deleted successfully",
	})
}
func UpdatePriceProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID    string  `json:"id"`
		Price float64 `json:"price"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	productId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		http.Error(w, `{"error": "Invalid product ID format"}`, http.StatusBadRequest)
		return
	}

	err = services.UpdatePriceProduct(productId, input.Price)
	if err != nil {
		http.Error(w, `{"error": "Failed to update product price"}`, http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message": "Product price updated successfully",
		"product": map[string]interface{}{
			"id":    input.ID,
			"price": input.Price,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func GetAllProductsByTitle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Title string `json:"title"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var products []models.Product
	if len(input.Title) > 0 {
		products, err = services.GetAllProductByTitle(input.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if len(input.Title) == 0 {
		products, err = services.GetAllProduct()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	response := map[string]interface{}{
		"message":  "All products by title retrieved successfully",
		"products": products,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var products []models.Product
	products, err := services.GetAllProduct()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":  "All products retrieved successfully",
		"products": products,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Title string `json:"title"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	products, err := services.GetProductsByCategoryTitle(input.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":  "Products by category title retrieved successfully",
		"products": products,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// func AddProductToCart(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
// 		return
// 	}
// 	var input struct {
// 		UserID    string `json:"user_id"`
// 		ProductID string `json:"product_id"`
// 		Quantity  int    `json:"quantity"`
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	UserID, err := primitive.ObjectIDFromHex(input.UserID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	ProductID, err := primitive.ObjectIDFromHex(input.ProductID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	//check ProductID
// 	//check does product has exist
// 	var product models.Product
// 	database.ProductCollection = database.GetCollection("testDB", "products")
// 	err = database.ProductCollection.FindOne(
// 		context.Background(),
// 		bson.M{"_id": ProductID},
// 	).Decode(&product)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	//check does quantity of product less than quantity user add to cart
// 	if product.Quantity < input.Quantity {
// 		http.Error(w, `{"Error": "Quantity not enought"}`, http.StatusBadRequest)
// 		return
// 	}
// 	//check does user exists
// 	var user models.User
// 	database.UserCollection = database.GetCollection("testDB", "users")
// 	err = database.UserCollection.FindOne(
// 		context.Background(),
// 		bson.M{"_id": UserID},
// 	).Decode(&user)
// 	if err != nil {
// 		http.Error(w, `{"Error":"Not found user"}`, http.StatusUnauthorized)
// 		return
// 	}
// 	var cart models.Cart
// 	database.CartCollection = database.GetCollection("testDB", "carts")
// 	//find cart of user
// 	err = database.CartCollection.FindOne(
// 		context.Background(),
// 		bson.M{"user_id": user.ID},
// 	).Decode(&cart)
// 	if err != nil {
// 		//case didn't have cart before
// 		if err == mongo.ErrNoDocuments {
// 			newCart := models.Cart{
// 				ID:     primitive.NewObjectID(),
// 				UserID: user.ID,
// 				Products: []models.Product{{
// 					ID:       product.ID,
// 					Title:    product.Title,
// 					Price:    product.Price,
// 					Quantity: product.Quantity,
// 				}},
// 				CreatedAt: time.Now(),
// 				UpdatedAt: time.Now(),
// 			}
// 			_, err = database.CartCollection.InsertOne(
// 				context.Background(),
// 				newCart)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 			}
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 	}
// 	//check does product Existed on cart
// 	productExists := false
// 	for i, v := range cart.Products {
// 		if v.ID == ProductID {
// 			cart.Products[i].Quantity += input.Quantity
// 			productExists = true
// 			break
// 		}
// 	}
// 	if !productExists {
// 		cart.Products = append(cart.Products, models.Product{
// 			ID:         product.ID,
// 			Title:      product.Title,
// 			Price:      product.Price,
// 			CategoryID: product.CategoryID,
// 			Quantity:   input.Quantity,
// 			Stock:      product.Stock,
// 		})
// 	}
// 	//update field products and updated_at after add product to cart
// 	cart.UpdatedAt = time.Now()
// 	_, err = database.CartCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": cart.ID},
// 		bson.M{"$set": bson.M{
// 			"products":   cart.Products,
// 			"updated_at": cart.UpdatedAt,
// 		}},
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	// update stock of product after user add product to cart
// 	_, err = database.ProductCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": product.ID},
// 		bson.M{"$inc": bson.M{
// 			"stock": -input.Quantity,
// 		}},
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	} else {

// 	}
// 	_, err = database.UserCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": UserID},
// 		bson.M{"$set": bson.M{
// 			"cart_id": cart.ID,
// 		}},
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	response := map[string]string{
// 		"message": "Product have add to cart successfully",
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }
