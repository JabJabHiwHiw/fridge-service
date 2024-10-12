package services

import (
	"context"
	"fmt"
	"time"

	"github.com/JabJabHiwHiw/fridge-service/internal/models"
	"github.com/JabJabHiwHiw/fridge-service/proto"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ proto.FridgeItemServiceServer = (*FridgeItemService)(nil)

type FridgeItemService struct {
	proto.UnimplementedFridgeItemServiceServer
	Collection *mongo.Collection
}

func (r *FridgeItemService) GetFridge(ctx context.Context, req *proto.FridgeRequest) (*proto.FridgeItemsResponse, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": req.GetUserId()})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var items []*proto.FridgeItem

	for cursor.Next(ctx) {
		var item models.FridgeItem

		err := cursor.Decode(&item)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		items = append(items, &proto.FridgeItem{
			Id:           item.ID,
			UserId:       item.UserID,
			IngredientId: item.IngredientID,
			Quantity:     item.Quantity,
			AddedDate:    timestamppb.New(*item.AddedDate),
			ExpiredDate:  timestamppb.New(*item.ExpiredDate),
		})
	}

	return &proto.FridgeItemsResponse{
		Items: items,
	}, nil
}

func (r *FridgeItemService) GetFridgeItem(ctx context.Context, req *proto.FridgeItemRequest) (*proto.FridgeItemResponse, error) {
	var item models.FridgeItem
	err := r.Collection.FindOne(ctx, bson.M{"_id": req.GetId()}).Decode(&item)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.FridgeItemResponse{
		Item: &proto.FridgeItem{
			Id:           item.ID,
			UserId:       item.UserID,
			IngredientId: item.IngredientID,
			Quantity:     item.Quantity,
			AddedDate:    timestamppb.New(*item.AddedDate),
			ExpiredDate:  timestamppb.New(*item.ExpiredDate),
		},
	}, nil
}

func (r *FridgeItemService) GetExpiredItems(ctx context.Context, req *proto.FridgeRequest) (*proto.FridgeItemsResponse, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{
		"user_id":      req.GetUserId(),
		"expired_date": bson.M{"$lt": time.Now()},
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var items []*proto.FridgeItem

	for cursor.Next(ctx) {
		var item models.FridgeItem

		err := cursor.Decode(&item)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		items = append(items, &proto.FridgeItem{
			Id:           item.ID,
			UserId:       item.UserID,
			IngredientId: item.IngredientID,
			Quantity:     item.Quantity,
			AddedDate:    timestamppb.New(*item.AddedDate),
			ExpiredDate:  timestamppb.New(*item.ExpiredDate),
		})
	}

	return &proto.FridgeItemsResponse{
		Items: items,
	}, nil
}

func (r *FridgeItemService) AddItem(ctx context.Context, req *proto.FridgeItem) (*proto.FridgeItemResponse, error) {
	item := models.FridgeItem{
		ID:           uuid.New().String(),
		UserID:       req.GetUserId(),
		IngredientID: req.GetIngredientId(),
		Quantity:     req.GetQuantity(),
	}

	if req.GetAddedDate() != nil {
		addedDate := req.GetAddedDate().AsTime()
		item.AddedDate = &addedDate
	}

	if req.GetExpiredDate() != nil {
		expiredDate := req.GetExpiredDate().AsTime()
		item.ExpiredDate = &expiredDate
	}

	_, err := r.Collection.InsertOne(ctx, item)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.FridgeItemResponse{
		Item: &proto.FridgeItem{
			Id:           item.ID,
			UserId:       item.UserID,
			IngredientId: item.IngredientID,
			Quantity:     item.Quantity,
			AddedDate:    timestamppb.New(*item.AddedDate),
			ExpiredDate:  timestamppb.New(*item.ExpiredDate),
		},
	}, nil
}

// AddItem input example
// {
//     "user_id": "555",
//     "ingredient_id": "d85af0b8-0b74-49a4-81aa-b36068404ab9",
//     "quantity": "5 kg",
//     "added_date": {
//         "seconds": 1697125800,
//         "nanos": 0
//     },
//     "expired_date": {
//         "seconds": 1697212200,
//         "nanos": 0
//     }
// }

func (r *FridgeItemService) UpdateItem(ctx context.Context, req *proto.FridgeItem) (*proto.FridgeItemResponse, error) {
	updates := bson.M{}

	if req.GetIngredientId() != "" {
		updates["ingredient_id"] = req.GetIngredientId()
	}
	if req.GetUserId() != "" {
		updates["user_id"] = req.GetUserId()
	}
	if req.GetQuantity() != "" {
		updates["quantity"] = req.GetQuantity()
	}
	if req.GetAddedDate() != nil {
		addedDate := req.GetAddedDate().AsTime()
		updates["added_date"] = &addedDate
	}
	if req.GetExpiredDate() != nil {
		expiredTime := req.GetExpiredDate().AsTime()
		updates["expired_date"] = &expiredTime
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": req.GetId()}, bson.M{"$set": updates})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var updatedItem models.FridgeItem
	err = r.Collection.FindOne(ctx, bson.M{"_id": req.GetId()}).Decode(&updatedItem)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var addedDate, expiredDate *timestamppb.Timestamp
	if updatedItem.AddedDate != nil {
		addedDate = timestamppb.New(*updatedItem.AddedDate)
	}
	if updatedItem.ExpiredDate != nil {
		expiredDate = timestamppb.New(*updatedItem.ExpiredDate)
	}

	return &proto.FridgeItemResponse{
		Item: &proto.FridgeItem{
			Id:           updatedItem.ID,
			UserId:       updatedItem.UserID,
			IngredientId: updatedItem.IngredientID,
			Quantity:     updatedItem.Quantity,
			AddedDate:    addedDate,
			ExpiredDate:  expiredDate,
		},
	}, nil
}

func (r *FridgeItemService) RemoveItem(ctx context.Context, req *proto.FridgeItemRequest) (*proto.Empty, error) {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.Empty{}, nil
}
