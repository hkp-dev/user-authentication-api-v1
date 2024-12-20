package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	FullName  string             `bson:"full_name" json:"full_name"`
	Phone     string             `bson:"phone" json:"phone"`
	Street    string             `bson:"street" json:"street"`
	City      string             `bson:"city" json:"city"`
	State     string             `bson:"state" json:"state"`
	ZipCode   string             `bson:"zip_code" json:"zip_code"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
