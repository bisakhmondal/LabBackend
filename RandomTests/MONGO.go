package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
	"time"
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

	dbs,_ := client.ListDatabaseNames(ctx,bson.M{})
	fmt.Println("Current Database: ",dbs)
	// fmt.Println("C Done")
	coll := client.Database("users").Collection("info")
	crsr,err :=coll.Find(ctx,bson.M{}) //query bson.M{"name":bson.M{"$eq":"Shuvayan"}}

	if err !=nil{
		log.Fatal(err)
	}
	var data Users

	if err = crsr.All(ctx, &data); err!=nil{
		log.Fatal(err)
	}
	log.Println(data)
	//Disconnect

}