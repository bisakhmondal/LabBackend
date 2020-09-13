package data

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	
)

//MongoClient for Database handling.
type MongoClient struct{
	ctx *context.Context
	client *mongo.Client
}

// New Client initialization.
func NewMongoClient(ct *context.Context, clt *mongo.Client) *MongoClient{
	return &MongoClient{
		ctx: ct,
		client: clt,
	}
}




