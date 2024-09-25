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

var _ proto.FridgeServiceServer = (*FridgeService)(nil)

type FridgeService struct {
	proto.UnimplementedFridgeServiceServer
	Collection *mongo.Collection
}

func (r *FridgeService) GetFridge(ctx context.Context, req *proto.FridgeRequest) (*proto.FridgeItemsResponse, error) {

	cursor, err := r.Collection.Find(ctx, bson.M{"owner": req.GetOwner()})

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
			Id:               item.ID,
			Owner:            item.Owner,
			Name:             item.Name,
			Type:             item.Type,
			Amount:           item.Amount,
			ManufacturedDate: timestamppb.New(*item.ManufacturedDate),
			ExpiredDate:      timestamppb.New(*item.ExpiredDate),
		})
	}

	return &proto.FridgeItemsResponse{
		Items: items,
	}, nil
}

func (r *FridgeService) GetFridgeItem(ctx context.Context, req *proto.FridgeItemRequest) (*proto.FridgeItemResponse, error) {
	var item models.FridgeItem
	err := r.Collection.FindOne(ctx, bson.M{"_id": req.GetId()}).Decode(&item)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.FridgeItemResponse{
		Item: &proto.FridgeItem{
			Id:               item.ID,
			Owner:            item.Owner,
			Name:             item.Name,
			Type:             item.Type,
			Amount:           item.Amount,
			ManufacturedDate: timestamppb.New(*item.ManufacturedDate),
			ExpiredDate:      timestamppb.New(*item.ExpiredDate),
		},
	}, nil
}

func (r *FridgeService) GetExpiredItems(ctx context.Context, req *proto.FridgeRequest) (*proto.FridgeItemsResponse, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{
		"owner":        req.GetOwner(),
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
			Id:               item.ID,
			Owner:            item.Owner,
			Name:             item.Name,
			Type:             item.Type,
			Amount:           item.Amount,
			ManufacturedDate: timestamppb.New(*item.ManufacturedDate),
			ExpiredDate:      timestamppb.New(*item.ExpiredDate),
		})
	}

	return &proto.FridgeItemsResponse{
		Items: items,
	}, nil
}

func (r *FridgeService) AddItem(ctx context.Context, req *proto.FridgeItem) (*proto.FridgeItemResponse, error) {

	item := models.FridgeItem{
		ID:     uuid.New().String(),
		Owner:  req.GetOwner(),
		Name:   req.GetName(),
		Type:   req.GetType(),
		Amount: req.GetAmount(),
	}

	if req.GetManufacturedDate() != nil {
		manufacturedDate := req.GetManufacturedDate().AsTime()
		item.ManufacturedDate = &manufacturedDate
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
			Id:               item.ID,
			Owner:            item.Owner,
			Name:             item.Name,
			Type:             item.Type,
			Amount:           item.Amount,
			ManufacturedDate: timestamppb.New(*item.ManufacturedDate),
			ExpiredDate:      timestamppb.New(*item.ExpiredDate),
		},
	}, nil
}

func (r *FridgeService) UpdateItem(ctx context.Context, req *proto.FridgeItem) (*proto.FridgeItemResponse, error) {
	updates := bson.M{}

	if req.GetName() != "" {
		updates["name"] = req.GetName()
	}
	if req.GetType() != "" {
		updates["type"] = req.GetType()
	}
	if req.GetAmount() != "" {
		updates["amount"] = req.GetAmount()
	}
	if req.GetManufacturedDate() != nil {
		manufacturedTime := req.GetManufacturedDate().AsTime()
		updates["manufactured_date"] = &manufacturedTime
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

	var manufacturedDate, expiredDate *timestamppb.Timestamp
	if updatedItem.ManufacturedDate != nil {
		manufacturedDate = timestamppb.New(*updatedItem.ManufacturedDate)
	}
	if updatedItem.ExpiredDate != nil {
		expiredDate = timestamppb.New(*updatedItem.ExpiredDate)
	}

	return &proto.FridgeItemResponse{
		Item: &proto.FridgeItem{
			Id:               updatedItem.ID,
			Owner:            updatedItem.Owner,
			Name:             updatedItem.Name,
			Type:             updatedItem.Type,
			Amount:           updatedItem.Amount,
			ManufacturedDate: manufacturedDate,
			ExpiredDate:      expiredDate,
		},
	}, nil
}

func (r *FridgeService) RemoveItem(ctx context.Context, req *proto.FridgeItemRequest) (*proto.Empty, error) {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.Empty{}, nil
}
