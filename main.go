package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mango struct {
	Name         string
	ShelfLife    int
	PriceX       int
	Availability string
}

func main() {
	fmt.Println("Hello, World")

	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		fmt.Println("Mongo Connect", err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println("Mongo init", err)
	}

	coll := client.Database("mango-test").Collection("inventory")

	res, err := coll.InsertOne(context.TODO(), Mango{
		Name:         "Totapuri",
		ShelfLife:    3,
		PriceX:       2,
		Availability: "South India",
	})

	fmt.Println("inserted:", res.InsertedID)
}
