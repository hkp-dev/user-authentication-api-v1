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
func GetProductByID(id primitive.ObjectID) (models.Product, error) {
	database.ProductCollection = database.GetCollection("testDB", "products")
	var product models.Product
	err := database.ProductCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
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
	defer productCursor.Close(context.Background())
	for productCursor.Next(context.Background()) {
		var product models.Product
		err := productCursor.Decode(&product)
		if err != nil {
			return nil, fmt.Errorf("Cursor decode error: %w", err)
		}
		products = append(products, product)
	}
	if err = productCursor.Err(); err != nil {
		return nil, fmt.Errorf("Cursor iteration error: %w", err)
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

// func AddProductToCart(userID, productID primitive.ObjectID, quantity int) error {
// 	cartCollection := database.GetCollection("testDB", "carts")
// 	productCollection := database.GetCollection("testDB", "products")
// 	userCollection := database.GetCollection("testDB", "users")
// 	var product models.Product

// 	// Find product user wants to add
// 	err := productCollection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
// 	if err != nil {
// 		return fmt.Errorf("product not found: %v", err)
// 	}

// 	// Check quantity
// 	if product.Stock < quantity {
// 		return fmt.Errorf("not enough stock for product: %s", product.Title)
// 	}

// 	// Find the user's cart
// var cart models.Cart
// err = cartCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&cart)
// if err != nil {
// 	// If cart does not exist, create a new cart for the user
// 	if err == mongo.ErrNoDocuments {
// 		newCart := models.Cart{
// 			ID:        primitive.NewObjectID(),
// 			UserID:    userID,
// 			Products:  []models.Product{{ID: product.ID, Title: product.Title, Description: product.Description, Price: product.Price, Quantity: quantity}},
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}
// 		_, insertErr := cartCollection.InsertOne(context.Background(), newCart)
// 		if insertErr != nil {
// 			return fmt.Errorf("failed to create new cart: %v", insertErr)
// 		}
// 		return nil
// 	}
// 	return fmt.Errorf("failed to query cart: %v", err)
// }
// 	// Check if product already exists in cart
// 	productExists := false
// 	for i, p := range cart.Products {
// 		if p.ID == productID {
// 			cart.Products[i].Quantity += quantity
// 			productExists = true
// 			break
// 		}
// 	}

// 	// If product doesn't exist, add it to the cart
// 	if !productExists {
// 		cart.Products = append(cart.Products, models.Product{
// 			ID:    product.ID,
// 			Title: product.Title,
// 			// Category:    GetCategoriesByTitle(product.Title),
// 			Description: product.Description,
// 			Price:       product.Price,
// 			Quantity:    quantity,
// 		})
// 	}

// 	// Update cart with the new product details
// 	cart.UpdatedAt = time.Now()
// 	_, updateErr := cartCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": cart.ID},
// 		bson.M{"$set": bson.M{
// 			"products":   cart.Products,
// 			"updated_at": cart.UpdatedAt,
// 		}},
// 	)
// 	if updateErr != nil {
// 		return fmt.Errorf("failed to update cart: %v", updateErr)
// 	}

// 	// Update product stock
// 	_, stockUpdateErr := productCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": productID},
// 		bson.M{"$inc": bson.M{"stock": -quantity}},
// 	)
// 	if stockUpdateErr != nil {
// 		return fmt.Errorf("failed to update product stock: %v", stockUpdateErr)
// 	}
// 	_, cartUserUpdatedErr := userCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": userID},
// 		bson.M{"$set": bson.M{
// 			"cart_id": cart.ID,
// 		}},
// 	)
// 	if cartUserUpdatedErr != nil {
// 		return fmt.Errorf("failed to update cart id for user %v", cartUserUpdatedErr)
// 	}
// 	return nil
// }
