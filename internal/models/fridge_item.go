package models

import (
	"time"
)

type FridgeItem struct {
	ID           string     `bson:"_id"`
	UserID       string     `bson:"user_id"`
	IngredientID string     `bson:"ingredient_id"`
	Quantity     string     `bson:"quantity"`
	AddedDate    *time.Time `bson:"added_date"`
	ExpiredDate  *time.Time `bson:"expired_date"`
}
