package main

import (
	"auth/data"
	"auth/handlers"
	"auth/server"
	"context"
	//gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func  checkGet( w http.ResponseWriter, r *http.Request) {
   log.Printf("GET %s\n" , r.RequestURI)

   w.Write([]byte("Hello World"))
}
//
//

// Test implementation of main

//type logger struct{
//    lo * log.Logger
//}


var (
	bindAddress = env.String("BIND_ADDRESS",false,":9090","bind address for server")
	certFile = os.Getenv("CertFile")
	certKey = os.Getenv("CertKey")
)
func main(){
    env.Parse()

    l :=  log.New(os.Stdout , "Authentication Server " , log.LstdFlags)

    //logg := &logger{l}

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


	SignInHan := handlers.NewSignIn(dbclient,l)

    sm := mux.NewRouter()

    loginRouter := sm.Methods( http.MethodPost).Subrouter()
    loginRouter.HandleFunc("/login" , SignInHan.Signin)

    getRouter := sm.Methods( http.MethodGet ).Subrouter()
    getRouter.HandleFunc( "/" , checkGet )

	server := server.New(sm,*bindAddress)

    // start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := server.ListenAndServeTLS(certFile,certKey)
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
    }()

    // trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	server.Shutdown(ctx)
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

