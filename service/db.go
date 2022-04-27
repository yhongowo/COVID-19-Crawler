package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/*
	MongoDB configuration
*/
const URI = "mongodb://localhost:27017"
const DATABASE = "2019-nCov"

var client *mongo.Client
var db *mongo.Database

func InitDB() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	db = client.Database(DATABASE)
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
}

func CloseDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
