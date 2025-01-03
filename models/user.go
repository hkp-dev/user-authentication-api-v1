package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"`
	Address   []Address          `bson:"address" json:"address"`
	Locked    bool               `bson:"locked" json:"locked"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	OTP       string             `bson:"otp" json:"otp"`
	OTPExpiry time.Time          `bson:"otp_expiry" json:"otp_expiry"`
	CartID    primitive.ObjectID `bson:"cart_id" json:"cart_id"`
}

// type LoginRequest struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=6"`
// }

// type LoginResponse struct {
// 	User  User   `json:"user"`
// 	Token string `json:"token"`
// }
