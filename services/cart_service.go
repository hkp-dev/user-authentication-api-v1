package services

import (
	"app/database"
	"app/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCart(user models.User) error {
	database.CartCollection = database.GetCollection("testDB", "carts")
	database.UserCollection = database.GetCollection("testDB", "users")
	var cart models.Cart
	cart.ID = primitive.NewObjectID()
	cart.UserID = user.ID
	// cart.Products = []models.Product{}
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()
	_, err := database.CartCollection.InsertOne(context.Background(), cart)
	if err != nil {
		return err
	}
	_, err = database.UserCollection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"cart_id": cart.ID}})
	if err != nil {
		return err
	}
	return nil
}
func GetCarts() ([]models.Cart, error) {
	database.CartCollection = database.GetCollection("testDB", "carts")
	var carts []models.Cart
	cursor, err := database.CartCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var cart models.Cart
		err := cursor.Decode(&cart)
		if err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func DeleteCart(idCart primitive.ObjectID) error {
	database.CartCollection = database.GetCollection("testDB", "carts")
	_, err := database.CartCollection.DeleteOne(context.Background(), bson.M{"user_id": idCart})
	if err != nil {
		return err
	}
	return nil
}
func GetCartByUserID(userID primitive.ObjectID) (models.Cart, error) {
	database.CartCollection = database.GetCollection("testDB", "carts")
	var cart models.Cart
	err := database.CartCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return cart, err
	}
	return cart, nil
}
func AddProductToCart(userID, productID primitive.ObjectID, quantity int) error {
	database.ProductCollection = database.GetCollection("testDB", "products")
	database.CartCollection = database.GetCollection("testDB", "carts")
	database.UserCollection = database.GetCollection("testDB", "users")
	//find product
	product, err := GetProductByID(productID)
	if err != nil {
		return err
	}

	//Check quantity user want to know does product has enough to add
	if product.Stock < quantity {
		return fmt.Errorf("not enough stock for product: %s", product.Title)
	}
	//find cart of user
	cart, err := GetCartByUserID(userID)
	if err != nil {
		return err
	}

	//Check if product already exists in cart
	//if product exists, update quantity
	productFound := false
	for i, p := range cart.Products {
		if p.ProductID == productID {
			cart.Products[i].Quantity += quantity
			productFound = true
			break
		}
	}
	//if product doesn't exist, add it to the cart
	if !productFound {
		cart.Products = append(cart.Products, models.CartItem{
			ProductID: productID,
			Quantity:  quantity,
		})
	}
	// Update cart with the new product details
	cart.UpdatedAt = time.Now()

	_, err = database.CartCollection.UpdateOne(context.Background(), bson.M{"_id": cart.ID}, bson.M{"$set": bson.M{"products": cart.Products, "updated_at": cart.UpdatedAt}})
	if err != nil {
		return err
	}

	// Update stock of product after user add product to cart
	_, err = database.ProductCollection.UpdateOne(context.Background(), bson.M{"_id": productID}, bson.M{"$inc": bson.M{"stock": -quantity}})
	if err != nil {
		return err
	}
	return nil
}
func RemoveProductFromCart(userID, productID primitive.ObjectID) error {
	database.CartCollection = database.GetCollection("testDB", "carts")
	database.ProductCollection = database.GetCollection("testDB", "products")
	//find cart by userID
	cart, err := GetCartByUserID(userID)
	if err != nil {
		return err
	}
	//remove product from cart
	var productQuantity int
	for i, p := range cart.Products {
		if p.ProductID == productID {
			productQuantity = p.Quantity
			cart.Products[i] = cart.Products[len(cart.Products)-1]
			cart.Products = cart.Products[:len(cart.Products)-1]
			break
		}
	}
	//update cart
	cart.UpdatedAt = time.Now()
	_, err = database.CartCollection.UpdateOne(context.Background(), bson.M{"_id": cart.ID}, bson.M{"$set": bson.M{"products": cart.Products, "updated_at": cart.UpdatedAt}})
	if err != nil {
		return err
	}
	//update stock of product by increment productQuantity
	_, err = database.ProductCollection.UpdateOne(context.Background(), bson.M{"_id": productID}, bson.M{"$inc": bson.M{"stock": productQuantity}})
	if err != nil {
		return err
	}
	return nil
}
