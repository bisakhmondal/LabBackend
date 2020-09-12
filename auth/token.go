package main

import ( 
	"log"
    jwt "github.com/dgrijalva/jwt-go"
    "time"
    "os"
    "os/signal"
    "context"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    gohandlers "github.com/gorilla/handlers"

)

// to be stored in a database later
var users = map[string]string{
    "user1" :"password1",
    "user2" : "password2",
}

type Credentials struct{
    Password string `json:"password"`
    Username string `json:"username"`
}

func (l *logger) Signin( w http.ResponseWriter , r *http.Request){
    l.lo.Printf("POST %s\n" , r.RequestURI)

    var creds Credentials
    log.Println( )
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error( w, "Error Deserializing credentials" , http.StatusBadRequest)
        return
    }

    expectedPassword , ok  := users[creds.Username]
    if !ok  || expectedPassword != creds.Password{
        http.Error( w , "Wrong Username or Password" , http.StatusBadRequest)
        return
    }

    sessionToken , err := createToken(creds.Username , creds.Password)
    if err != nil{
        http.Error( w , "Cannot create token" , http.StatusBadRequest)
        return
    }


    //TODO : Implement Redis layer

    // _, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
	// if err != nil {
	// 	// If there is an error in setting the cache, return an internal server error
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
    // }
    

    http.SetCookie(w , &http.Cookie{
        Name : "session_token",
        Value : sessionToken, 
        Expires : time.Now().Add( 1 * time.Hour ),
    })


}


func (l* logger ) MiddlewareAuthenticate( next http.Handler ) http.Handler{
    return http.HandlerFunc(func(rw http.ResponseWriter , r *http.Request){
      

        _ , err := r.Cookie("session_token")
        if err != nil {
            if err == http.ErrNoCookie{
                http.Error( rw ,"Unauthorized" , http.StatusUnauthorized)
                return
            }
            http.Error(rw ,"Bad Request" , http.StatusBadRequest )
            return
        }
        
        next.ServeHTTP(rw ,r)
    })
}


func (l* logger) Refresh(w http.ResponseWriter, r *http.Request) {

    l.lo.Printf("GET %s\n" , r.RequestURI)
	c, err := r.Cookie("session_token")
	if err != nil {
        if err == http.ErrNoCookie{
            http.Error( w ,"Unauthorized" , http.StatusUnauthorized)
            return
        }
        http.Error( w ,"Bad Request" , http.StatusBadRequest )
		return
	}
    sessionToken := c.Value
    claims := jwt.MapClaims{}
    _ , errr := jwt.ParseWithClaims(sessionToken,  claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("API_SECRET")), nil
    })
    

    if  errr != nil{
        http.Error(w , "Wrong token signature",  http.StatusBadRequest)
        return
    }
    newSessionToken , err := createToken( claims["name"].(string) , claims["password"].(string))
    if err != nil{
        http.Error( w , "Cannot create token" , http.StatusBadRequest)
        return
    }


	

	// Now, create a new session token for the current user
	// newSessionToken := uuid.NewV4().String()
	// _, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s",response))
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// Delete the older session token
	// _, err = cache.Do("DEL", sessionToken)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	
	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(1 * time.Hour),
	})
}



func (l *logger) checkGet( w http.ResponseWriter, r *http.Request) {
    l.lo.Printf("GET %s\n" , r.RequestURI)

    w.Write([]byte("Hello World"))
}



func createToken( name string , password string ) ( string , error ){
    claims := jwt.MapClaims{}
    claims["authorized"]=true
    claims["name"]=name
    claims["password"]=password
    claims["exp"]=time.Now().Add(time.Hour * 1 ).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("API_SECRET")))

} 






// Test implementation of main

type logger struct{
    lo * log.Logger
}


var bindAddress = "http://localhost:9090"

func main(){
    
    l :=  log.New(os.Stdout , "test-run-of-auth" , log.LstdFlags)

    logg := &logger{l}

    
    sm := mux.NewRouter()

    loginRouter := sm.Methods( http.MethodPost).Subrouter()
    loginRouter.HandleFunc("/login" , logg.Signin )

    getRouter := sm.Methods( http.MethodGet ).Subrouter()
    getRouter.HandleFunc( "/get" , logg.checkGet )



    // CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         "localhost:9090",      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
    }
    
    // start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
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
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
    








}
