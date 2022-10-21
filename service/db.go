package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

const URI = ""
const DATABASE = "2019-nCov"

var (
	client *mongo.Client
	db     *mongo.Database
)

func InitDB() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	db = client.Database(DATABASE)
	err = db.CreateCollection(context.TODO(), "Overall")
	err = db.CreateCollection(context.TODO(), "Abroad")
	err = db.CreateCollection(context.TODO(), "Area")
	err = db.CreateCollection(context.TODO(), "Timeline")
	if err != nil {
		log.Println(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
}
