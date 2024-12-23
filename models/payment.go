package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Method    string             `bson:"method" json:"method"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type PaymentMethod struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required"` // "cash", "card", "netbanking"
	Description string             `bson:"description" json:"description,omitempty"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	Fee         float64            `bson:"fee" json:"fee,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
