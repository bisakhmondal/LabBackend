package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"data"

)

type user struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	NAME string `bson:"name,omitempty"`
	STREAM string `bson:"stream,omitempty"`
}
type Users []user

func main(){
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	URI := getURI("MONGO")
	fmt.Println(URI)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	defer client.Disconnect(ctx)
	
	if err !=nil{
		log.Fatal("Unable to Connect")
	}

	dbs, _ := client.ListDatabaseNames(ctx,bson.M{})
	fmt.Println("Current Database: ",dbs)
	// fmt.Println("C Done")
	coll := client.Database("users").Collection("info")

	//Fetch
	crsr,err :=coll.Find(ctx , bson.M{}) //query bson.M{"name":bson.M{"$eq":"Shuvayan"}}
	
	//inUser := user{
	//	NAME:   "NewStudent",
	//	STREAM: "PG",
	//}
	if err !=nil{
		log.Fatal(err)
	}
	var data Users

	err = crsr.All(ctx, &data); check(err)
	log.Println(data)

	//Insert
	//inserted,err := coll.InsertOne(ctx, inUser)
	//
	//check(err)
	//log.Println(inserted)
	//selective fetch

	// crs2,err := coll.Find(ctx,bson.M{}, options.Find().SetProjection(bson.M{"_id":0,"name":1}))

	// var data2 Users
	// check(err)
	// err = crs2.All(ctx, &data2); check(err)
	// log.Println(data2)

}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}