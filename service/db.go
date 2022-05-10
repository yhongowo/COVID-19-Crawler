package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
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
	//Create Database and Time-series Collections
	db = client.Database(DATABASE)
	tso := options.TimeSeries().SetTimeField("updateTime").SetGranularity("hours")
	opts := options.CreateCollection().SetTimeSeriesOptions(tso)
	err = db.CreateCollection(context.TODO(), "Overall", opts)
	err = db.CreateCollection(context.TODO(), "Abroad", opts)
	err = db.CreateCollection(context.TODO(), "Province", opts)
	err = db.CreateCollection(context.TODO(), "Timeline", opts)
	if err != nil {
		log.Println("Collections created")
	}

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
