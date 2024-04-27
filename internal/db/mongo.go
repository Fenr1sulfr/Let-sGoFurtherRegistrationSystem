package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnetWithMongo() *mongo.Client {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://sulfur:sulfur@cluster0.depm68j.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if err := client.Database("sample_user").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}


