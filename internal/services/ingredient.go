package services

import (
	"context"
	"fmt"

	"github.com/JabJabHiwHiw/fridge-service/internal/models"
	"github.com/JabJabHiwHiw/fridge-service/proto"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ proto.IngredientServiceServer = (*IngredientService)(nil)

type IngredientService struct {
	proto.UnimplementedIngredientServiceServer
	Collection *mongo.Collection
}

func (r *IngredientService) GetIngredients(ctx context.Context, req *proto.Empty) (*proto.IngredientsResponse, error) {

	cursor, err := r.Collection.Find(ctx, bson.D{})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var ingredients []*proto.IngredientItem

	for cursor.Next(ctx) {
		var ingredient models.Ingredient

		err := cursor.Decode(&ingredient)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		ingredients = append(ingredients, &proto.IngredientItem{
			Id:       ingredient.ID,
			Name:     ingredient.Name,
			Category: ingredient.Category,
		})
	}

	return &proto.IngredientsResponse{
		Ingredients: ingredients,
	}, nil
}

func (r *IngredientService) GetIngredientItem(ctx context.Context, req *proto.IngredientItemRequest) (*proto.IngredientItemResponse, error) {
	var ingredient models.Ingredient
	err := r.Collection.FindOne(ctx, bson.M{"_id": req.GetId()}).Decode(&ingredient)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.IngredientItemResponse{
		Ingredient: &proto.IngredientItem{
			Id:       ingredient.ID,
			Name:     ingredient.Name,
			Category: ingredient.Category,
		},
	}, nil
}

func (r *IngredientService) AddIngredient(ctx context.Context, req *proto.IngredientItem) (*proto.IngredientItemResponse, error) {

	ingredient := models.Ingredient{
		ID:       uuid.New().String(),
		Name:     req.GetName(),
		Category: req.GetCategory(),
	}

	_, err := r.Collection.InsertOne(ctx, ingredient)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.IngredientItemResponse{
		Ingredient: &proto.IngredientItem{
			Id:       ingredient.ID,
			Name:     ingredient.Name,
			Category: ingredient.Category,
		},
	}, nil
}
