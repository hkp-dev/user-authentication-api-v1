package utils

import (
	"app/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	// claims["authorized"] = true
	// claims["created_at"] = user.CreatedAt.Unix()
	// claims["updated_at"] = user.UpdatedAt.Unix()
	// claims["otp"] = user.OTP
	// claims["otp_expiry"] = user.OTPExpiry.Unix()
	claims["_id"] = user.ID.Hex()
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["locked"] = user.Locked
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func ValidateJWT(r *http.Request) (models.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return models.User{}, fmt.Errorf("Authorization header is missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return models.User{}, fmt.Errorf("Invalid authorization format")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})

	if err != nil {
		return models.User{}, fmt.Errorf("Invalid token: %s", err.Error())
	}

	if !token.Valid {
		return models.User{}, fmt.Errorf("Token is not valid")
	}

	userID, ok := claims["_id"].(string)
	if !ok {
		return models.User{}, fmt.Errorf("User ID not found in token")
	}

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("Invalid user ID in token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return models.User{}, fmt.Errorf("Username not found in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return models.User{}, fmt.Errorf("Email not found in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return models.User{}, fmt.Errorf("Role not found in token")
	}

	locked, ok := claims["locked"].(bool)
	if !ok {
		return models.User{}, fmt.Errorf("Locked status not found in token")
	}

	user := models.User{
		ID:       userIDObj,
		Username: username,
		Email:    email,
		Role:     role,
		Locked:   locked,
	}

	return user, nil
}

// createdAtUnix, ok := claims["created_at"].(float64)
// if !ok {
// 	return models.User{}, fmt.Errorf("Created_at not found in token")
// }
// createdAt := time.Unix(int64(createdAtUnix), 0)

// updatedAtUnix, ok := claims["updated_at"].(float64)
// if !ok {
// 	return models.User{}, fmt.Errorf("Updated_at not found in token")
// }
// updatedAt := time.Unix(int64(updatedAtUnix), 0)

// otp, ok := claims["otp"].(string)
// if !ok {
// 	return models.User{}, fmt.Errorf("OTP not found in token")
// }

// otpExpiryUnix, ok := claims["otp_expiry"].(float64)
// if !ok {
// 	return models.User{}, fmt.Errorf("OTP expiry not found in token")
// }
// otpExpiry := time.Unix(int64(otpExpiryUnix), 0)
