package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type db struct {
	client      *mongo.Client
	DXYOverall  *mongo.Collection
	DXYNews     *mongo.Collection
	DXYArea     *mongo.Collection
	DXYProvince *mongo.Collection
}

func NewDB() *db {
	client := newClient()
	DXYOverall := client.Database("2019-nCov").Collection("DXYOverall")
	DXYNews := client.Database("2019-nCov").Collection("DXYNews")
	DXYArea := client.Database("2019-nCov").Collection("DXYArea")
	DXYProvince := client.Database("2019-nCov").Collection("DXYProvince")
	return &db{
		client:      client,
		DXYOverall:  DXYOverall,
		DXYNews:     DXYNews,
		DXYArea:     DXYArea,
		DXYProvince: DXYProvince,
	}
}

func newClient() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return client
}
