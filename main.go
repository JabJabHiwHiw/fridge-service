package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/JabJabHiwHiw/fridge-service/internal/services"
	"github.com/JabJabHiwHiw/fridge-service/proto"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:example@mongodb:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connected to MongoDB")

	db := client.Database("fridge-service")
	fridgeItemCollection := db.Collection("fridge_item")
	ingredientCollection := db.Collection("ingredient")

	fridgeItemService := services.FridgeItemService{
		Collection: fridgeItemCollection,
	}

	ingredientService := services.IngredientService{
		Collection: ingredientCollection,
	}

	grpcServer := grpc.NewServer()
	proto.RegisterFridgeItemServiceServer(grpcServer, &fridgeItemService)
	proto.RegisterIngredientServiceServer(grpcServer, &ingredientService)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started on port :50052")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
