package controllers_test

import (
	"app/controllers"
	"app/database"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"`
	Locked    bool               `bson:"locked" json:"locked"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	OTP       string             `bson:"otp" json:"otp"`
	OTPExpiry time.Time          `bson:"otp_expiry" json:"otp_expiry"`
	CartID    primitive.ObjectID `bson:"cart_id" json:"card_id"`
}

//	func AddUserAfterTest(user models.User) error {
//		user = models.User{
//			ID: primitive.NewObjectID(),
//		}
//	}
func DeleteUserAfterTest(userID primitive.ObjectID) error {
	userCollection := database.GetCollection("testDB", "users")
	_, err := userCollection.DeleteOne(context.Background(), bson.M{"_id": userID})
	if err != nil {
		return err
	}
	return nil
}
func TestMain(m *testing.M) {
	err := database.ConnectMongoDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("Couldn't connect to mongoDB: %v", err)
	}
	code := m.Run()
	if database.MongoClient != nil {
		_ = database.MongoClient.Disconnect(context.Background())
	}
	os.Exit(code)
}
func TestLockUser(t *testing.T) {
	userCollection := database.GetCollection("testDB", "users")
	userID := primitive.NewObjectID()
	_, err := userCollection.InsertOne(context.Background(), bson.M{
		"_id":     userID,
		"locked":  false,
		"created": time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}
	reqBody := map[string]primitive.ObjectID{"id": userID}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/unlock-user", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	controllers.LockUser(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response["message"] != "User locked successfully" {
		t.Errorf("Expected message 'User locked successfully', got %v", response["message"])
	}
	err = DeleteUserAfterTest(userID)
	if err != nil {
		t.Errorf("Failed to delete user after test: %v", err)
	}
}
func TestUnlockUser(t *testing.T) {
	userCollection := database.GetCollection("testDB", "users")
	userID := primitive.NewObjectID()
	_, err := userCollection.InsertOne(context.Background(), bson.M{
		"_id":     userID,
		"locked":  true,
		"created": time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}
	reqBody := map[string]primitive.ObjectID{"id": userID}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/unlock-user", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	controllers.UnlockUser(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if response["message"] != "User unlocked successfully" {
		t.Errorf("Expected message 'User unlocked successfully', got %v", response["message"])
	}
	err = DeleteUserAfterTest(userID)
	if err != nil {
		t.Errorf("Failed to delete user after test: %v", err)
	}
}
func TestDeleteUser(t *testing.T) {
	userCollection := database.GetCollection("testDB", "users")
	userID := primitive.NewObjectID()
	_, err := userCollection.InsertOne(context.Background(), bson.M{
		"_id":     userID,
		"locked":  false,
		"created": time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}
	reqBody := map[string]primitive.ObjectID{"id": userID}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/delete-user", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	controllers.DeleteUser(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "User deleted successfully" {
		t.Errorf("Expected message 'User deleted successfully', got %v", response["message"])
	}

	var deletedUser bson.M
	err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&deletedUser)
	if err == nil {
		t.Errorf("Expected user to be deleted, but found user in database")
	}
}