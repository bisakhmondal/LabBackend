package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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


//getting data from db.
func (p* MongoClient)GetData() (Plist, error){
	collection := p.client.Database("users").Collection("info")

	var personsData Plist
	ctx,cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	curr, err := collection.Find(ctx, bson.M{},options.Find().SetProjection(bson.M{"_id":0,"password":0,"username":0}))

	if err!=nil{
		return nil, err
	}
	if err = curr.All(ctx,&personsData);err!=nil{
		return nil, err
	}

	return personsData,nil
}

func (p* MongoClient)FindUser(route *string) (*Person,error){

	coll := p.client.Database("users").Collection("info")
	ctx,cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	var personData Person

	err := coll.FindOne(
		ctx,
		bson.M{"route":route},
		options.FindOne().SetProjection(bson.M{"_id":0}),
	).Decode(&personData)

	if err!=nil{
		return nil,err
	}
	return &personData, nil
}
