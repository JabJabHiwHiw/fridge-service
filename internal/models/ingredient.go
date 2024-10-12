package models

type Ingredient struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Category string `bson:"category"`
}
