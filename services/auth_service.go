package services

import (
	"app/database"
	"app/models"
	"app/utils"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(user models.User) []string {
	var errors []string
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field %s: %s", e.Field(), e.Tag()))
		}
	}
	return errors
}
func GetUserByID(userID primitive.ObjectID) (models.User, error) {
	database.UserCollection = database.GetCollection("testDB", "users")
	var user models.User
	err := database.UserCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, err
}
func CheckUserExists(user models.User) (bool, error) {
	database.UserCollection = database.GetCollection("testDB", "users")
	filter := bson.M{"$or": []interface{}{
		bson.M{"username": user.Username},
		bson.M{"email": user.Email},
	}}
	var existingUser models.User
	err := database.UserCollection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return true, err
}
func CreateUser(user *models.User, cart *models.Cart) error {
	database.UserCollection = database.GetCollection("testDB", "users")
	user.ID = primitive.NewObjectID()
	if user.Email == "admin@gmail.com" {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.CartID = cart.ID
	_, err := database.UserCollection.InsertOne(context.Background(), user)
	return err
}
func CheckStatusUser(user models.User) bool {
	database.UserCollection = database.GetCollection("testDB", "users")
	filter := bson.M{"_id": user.ID}
	err := database.UserCollection.FindOne(context.Background(), filter).Decode(&user)
	if user.Locked {
		return false
	}
	return err == nil
}
func FindUserByUsername(username string) (models.User, error) {
	var user models.User
	database.UserCollection = database.GetCollection("testDB", "users")
	err := database.UserCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	return user, err
}
func ComparePassword(password string, user models.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
func GetUserByToken(r *http.Request) (models.User, error) {
	var user models.User
	user, err := utils.ValidateJWT(r)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	database.UserCollection = database.GetCollection("testDB", "users")
	err := database.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return user, err
}
func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
func UpdatePassword(username, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = database.UserCollection.UpdateOne(
		context.Background(),
		bson.M{"username": username},
		bson.M{"$set": bson.M{"password": string(hashedPassword), "updated_at": time.Now()}},
	)
	return err
}
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06v", rand.Intn(999999))
}
func GenerateExpiry() time.Time {
	return time.Now().Add(5 * time.Minute)
}
func SaveOTP(otp string, expiry time.Time, user models.User) error {
	database.UserCollection = database.GetCollection("testDB", "users")
	_, err := database.UserCollection.UpdateOne(context.Background(), bson.M{"email": user.Email}, bson.M{
		"$set": bson.M{
			"otp":        otp,
			"otp_expiry": expiry,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
func GetOTP(user models.User) (string, time.Time, error) {
	database.UserCollection = database.GetCollection("testDB", "users")
	var u models.User
	err := database.UserCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&u)
	if err != nil {
		return "", time.Time{}, err
	}
	return u.OTP, u.OTPExpiry, nil
}
