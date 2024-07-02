package models

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	log.Println("Item package initialized")
}

type Item struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Price     float64            `bson:"price" json:"price"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
