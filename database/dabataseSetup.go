package database

import (
	"context"
	"fmt"

	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongo://localhost:27017"))

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// defer cancel()

	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.Ping(context.TODO(), nil)
	// if err != nil {
	// 	log.Println("failed to connect mongodb")
	// 	return nil
	// }
	// fmt.Println("successfully connected to mongoDB")
	// return client
	clientOptions := options.Client().ApplyURI("mongodb://development:testpassword@localhost:27017/?authSource=admin")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("failed to connect mongodb")
		return nil
	}
	fmt.Println("successfully connected to mongoDB")
	return client

}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var userCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return userCollection

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productCollection
}
