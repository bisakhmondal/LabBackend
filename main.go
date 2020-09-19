package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"fmt"
	"os/signal"
	"serving-api/data"
	"serving-api/handlers"
	"serving-api/server"
	"time"
)


var (
	//bindAddress = env.String("BIND_ADDRESS",false,":8080","bind address for server")
	bindAddress2 = GetPort()//env.String("BIND_ADDRESS",false,":9090","bind address for server")
	certFile = os.Getenv("CertFile")
	certKey = os.Getenv("CertKey")
	
)
func main(){

	env.Parse()

	l := log.New(os.Stdout,"Serving-api ",log.LstdFlags)

	//Mongo Connect

	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	URI := os.Getenv("MONGO")
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
	//Server := server.New(smux,*bindAddress,true)
	
	//Create a new http Server
	Serverhttp :=server.New(smux,bindAddress2,false)

	/*//starting server
	go func(){
		l.Println("HTTPS starting at Port: ",*bindAddress)

		err := Server.ListenAndServeTLS(certFile, certKey)
		
		if err!=nil{
			l.Fatal("Server starting Failed",err)
		}
	}()
*/
	go func(){
		l.Println("HTTP starting at Port: ",bindAddress2)

		err := Serverhttp.ListenAndServe()

		if err!=nil{
			l.Fatal("Server starting Failed",err)
		}
	}()
	
	ch := make(chan os.Signal)

	signal.Notify(ch, os.Kill)
	signal.Notify(ch, os.Interrupt)

	sig := <-ch
	
	l.Println("Shutting Down... ",sig)

	//Server.Shutdown(ctx)
	Serverhttp.Shutdown(ctx)

}

func GetPort() string {
	 	var port = os.Getenv("PORT")
	 	// Set a default port if there is nothing in the environment
	 	if port == "" {
	 		port = "4747"
	 		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	 	}
		return ":" + port
	}
