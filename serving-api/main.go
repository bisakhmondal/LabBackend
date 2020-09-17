package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"serving-api/data"
	"serving-api/handlers"
	"serving-api/server"
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
	
	if err !=nil{
		l.Fatal("Unable to Connect")
	}
	
	// defer Disconnect
	defer client.Disconnect(ctx)

	//wrapped mongo client
	dbclient := data.NewMongoClient(&ctx, client)

	pHandler := handlers.NewPersonH(dbclient,l)

	//Gorilla router
	smux := mux.NewRouter()

	getRouter := smux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", pHandler.GetData)
	getRouter.HandleFunc("/user/{route:[a-zA-Z0-9]+}", pHandler.FetchUser)

	//Create a new tls server
	Server := server.New(smux, *bindAddress)

	//starting server
	go func(){
		l.Println("starting at Port: ",*bindAddress)

		err := Server.ListenAndServeTLS(certFile, certKey)
		
		if err!=nil{
			l.Fatal("Server starting Failed",err)
		}
	}()
	
	ch := make(chan os.Signal)

	signal.Notify(ch, os.Kill)
	signal.Notify(ch, os.Interrupt)

	sig := <-ch
	
	l.Println("Shutting Down... ",sig)

	Server.Shutdown(ctx)

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
