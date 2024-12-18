package services

import (
	"app/database"
	"app/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckStatusLocked(id primitive.ObjectID) (bool, error) {
	database.UserCollection = database.GetCollection("testDB", "users")
	var user models.User
	err := database.UserCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return false, err
	}
	return user.Locked, nil
}
func LockUser(id primitive.ObjectID, StatusLocked bool) error {
	if id.IsZero() {
		return fmt.Errorf("user ID is empty")
	}
	database.UserCollection = database.GetCollection("testDB", "users")
	isLocked, err := CheckStatusLocked(id)
	if err != nil {
		return err
	}
	if isLocked {
		return fmt.Errorf("user with ID %v already locked", id)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"locked": true}}
	_, err = database.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
func UnlockUser(id primitive.ObjectID) error {
	if id.IsZero() {
		return fmt.Errorf("user ID is empty")
	}
	database.UserCollection = database.GetCollection("testDB", "users")
	Islocked, err := CheckStatusLocked(id)
	if err != nil {
		return err
	}
	if !Islocked {
		return fmt.Errorf("user with ID %v already unlocked", id)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"locked": false}}
	_, err = database.UserCollection.UpdateOne(context.Background(), filter, update)
	return err
}
func DeleteUser(id primitive.ObjectID) error {
	database.UserCollection = database.GetCollection("testDB", "users")
	_, err := database.UserCollection.DeleteOne(
		context.Background(),
		bson.M{"_id": id},
	)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
