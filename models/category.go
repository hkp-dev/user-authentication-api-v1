package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title" validate:"required,max=30"`
	Description string               `bson:"description" json:"description" validate:"required,min=8,max=200"`
	ProductIDs  []primitive.ObjectID `bson:"product_ids" json:"product_ids"`
}