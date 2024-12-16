package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title" validate:"required,min=8,max=50"`
	Description string               `bson:"description" json:"description" validate:"required,min=8,max=200"`
	Price       float64              `bson:"price" json:"price" validate:"required"`
	CategoryID  []primitive.ObjectID `bson:"category_id" json:"category_id" validate:"required"`
	Quantity    int                  `bson:"quantity" json:"quantity" validate:"required"`
	Stock       int                  `bson:"stock" json:"stock" validate:"required"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}
