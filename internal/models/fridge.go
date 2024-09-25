package models

import (
	"time"
)

type FridgeItem struct {
	ID               string     `bson:"_id"`
	Owner            string     `bson:"owner"`
	Name             string     `bson:"name"`
	Type             string     `bson:"type"`
	Amount           string     `bson:"amount"`
	ManufacturedDate *time.Time `bson:"manufactured_date"`
	ExpiredDate      *time.Time `bson:"expired_date"`
}
