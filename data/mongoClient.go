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


////getting data from db.
//func (p* MongoClient)GetData() (Plist, error){
//	collection := p.client.Database("users").Collection("info")
//
//	var personsData Plist
//	ctx,cancel := context.WithTimeout(context.TODO(),10*time.Second)
//	defer cancel()
//	curr, err := collection.Find(ctx, bson.M{},options.Find().SetProjection(bson.M{"_id":0,"password":0,"username":0}))
//
//	if err!=nil{
//		return nil, err
//	}
//	if err = curr.All(ctx,&personsData);err!=nil{
//		return nil, err
//	}
//
//	return personsData,nil
//}

//Authentication route.
func (p* MongoClient) CheckAuth(personData *Person, filter *bson.M) error{
	collection := p.client.Database("users").Collection("info")

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	option := options.FindOne().SetProjection(bson.M{
		"username":1,
		"password":1,
	})

	err := collection.FindOne(ctx , filter, option).Decode(personData)
	if  err != nil {
		return err
	}

	return nil
}

//Update Database.
func (p* MongoClient)UpdateDB(personData *Person)error{
	coll := p.client.Database("users").Collection("info")

	ctx,cancel := context.WithTimeout(context.TODO(),5*time.Second)
	defer cancel()

	filter := bson.M{"_id":personData.ID}
	update,err := personData.ToBSON()

	if err!=nil{
		return err
	}
	_ = coll.FindOneAndUpdate(ctx,filter, bson.M{"$set":update})

	return nil
}
