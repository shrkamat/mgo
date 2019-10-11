package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"

	log "github.com/sirupsen/logrus"
)

type Mango struct {
	Id           string `json:"id" bson:"_id,omitempty"`
	Name         string
	ShelfLife    int
	PriceX       int
	Availability string
}

type DBCollection struct {
	coll *mongo.Collection
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

type Mangoes []Mango

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

	mango := Mango{
		Id:           "PEL~PED~122~langda",
		Name:         "lagnda",
		ShelfLife:    6,
		PriceX:       4,
		Availability: "North India",
	}

	err = InsertOne(coll, mango)

	logCollection(coll)

	tota := Mango{
		Id:           "PEL~PED~122~langda",
		Name:         "lagnda",
		ShelfLife:    10,
		PriceX:       8,
		Availability: "South India",
	}

	UpdateOne(coll, bson.D{{"_id", "PEL~PED~122~langda"}}, tota)

	logCollection(coll)

	count, _ := DeleteMany(context.TODO(), coll, bson.D{{}})

	fmt.Println("deleted ", count, "items")
}

func UpdateOne(coll *mongo.Collection, filter interface{}, data interface{}) error {

	res := coll.FindOneAndReplace(context.TODO(), filter, data)

	if err := res.Err(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Upsert(coll *mongo.Collection, filter interface{}, data interface{}) error {

	opts := options.FindOneAndReplace().SetUpsert(true)

	res := coll.FindOneAndReplace(context.TODO(), filter, data, opts)

	err := res.Err()
	if err != nil {
		fmt.Println(err)
		fmt.Println(res.DecodeBytes())
		return err
	}

	bytes, _ := res.DecodeBytes()

	fmt.Println("update", bytes)

	return nil
}

func DeleteMany(ctx context.Context, coll *mongo.Collection, filter interface{}) (int64, error) {
	dResult, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println("Failed to Del", err)
		return 0, err
	}
	return dResult.DeletedCount, nil
}

func InsertOne(coll *mongo.Collection, entity interface{}) error {
	res, err := coll.InsertOne(context.TODO(), entity)
	if err != nil {
		fmt.Println("Failed to insert", err)
	} else {
		fmt.Println("inserted:", res.InsertedID)
	}
	return err
}

func logCollection(coll *mongo.Collection) {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var mangoes Mangoes
	err = cursor.All(context.TODO(), &mangoes)

	if err != nil {
		fmt.Println("fetch err", err)
	}

	fmt.Println(mangoes)
}
