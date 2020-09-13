package main

import (
	"context"
	// "github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	// "net/http"
	"os"
	// "os/signal"
	"serving-api/data"
	// "serving-api/server"
	"time"
)


var (
	bindAddress = env.String("BIND_ADDRESS",false,":8080","bind address for server")
	certFile = os.Getenv("CertFile")
	certKey = os.Getenv("CertKey")
)
func main(){

	env.Parse()

	l := log.New(os.Stdout,"Serving-api ",log.LstdFlags)

	//Mongo Connect

	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	URI := getURI("MONGO")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	defer client.Disconnect(ctx)
	
	if err !=nil{
		l.Fatal("Unable to Connect")
	}
	
	// defer Disconnect
	defer client.Disconnect(ctx)

	//wrapped mongo client
	mclient := data.NewMongoClient(&ctx, client)



	//testing mongo client 

	dbs, _ := client.ListDatabaseNames(ctx,bson.M{})
	l.Println("Current Database: ",dbs)
	// fmt.Println("C Done")
	coll := client.Database("users").Collection("persons")


	PersonList := data.PersonList

	for person := range PersonList{
		inserted,err := coll.InsertOne( ctx , person)
		if err != nil{
			l.Println(err.Error() )
			l.Fatal("Failed inserting")
			
		} else {
			l.Println("Inserted")
		}
		l.Println(inserted)
		
	}





	// test end

	//just testing
	log.Println(mclient)

	//Gorilla router
	// smux := mux.NewRouter()

	// getRouter := smux.Methods(http.MethodGet).Subrouter()
	// getRouter.HandleFunc("/",func(rw http.ResponseWriter, r* http.Request){
	// 	rw.Header().Set("Content-Type","text/plain; charset=utf-8")

	// 	l.Println("Received on ROute '/' ")
	// 	rw.WriteHeader(http.StatusOK)
	// 	rw.Write([]byte("Hello WOrLd"))
	// })

	// //Create a new tls server
	// Server := server.New(smux, *bindAddress)

	// //starting server
	// go func(){
	// 	l.Println("starting at Port: ",*bindAddress)

	// 	err := Server.ListenAndServeTLS(certFile, certKey)
		
	// 	if err!=nil{
	// 		l.Fatal("Server starting Failed",err)
	// 	}
	// }()
	
	// ch := make(chan os.Signal)

	// signal.Notify(ch, os.Kill)
	// signal.Notify(ch, os.Interrupt)

	// sig := <-ch
	
	// l.Println("Shutting Down... ",sig)

	// Server.Shutdown(ctx)

}

//Get URI from Environment yaml file
func getURI(key string) string{
	viper.SetConfigFile("config.yaml")

	err:= viper.ReadInConfig()
	
	if err !=nil {
		log.Fatal("Can't Read Env Variable",err)
	}

	value,ok := viper.Get(key).(string)
	
	if !ok{
		log.Fatal("Invalid TypeCast")
	}

	return value
}
