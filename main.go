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
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

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
	collection := db.Collection("fridge")

	fridgeService := services.FridgeService{
		Collection: collection,
	}

	grpcServer := grpc.NewServer()
	proto.RegisterFridgeServiceServer(grpcServer, &fridgeService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started on port :8080")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
