package handlers

import (
	"auth/data"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"time"
)

type SignIn struct {
	db *data.MongoClient
	l *log.Logger
}

func NewSignIn(db * data.MongoClient,l *log.Logger) *SignIn{
	return &SignIn{
		db: db,
		l: l,
	}
}

func (l *SignIn)Signin( w http.ResponseWriter , r *http.Request){
	l.l.Printf("POST %s\n" , r.RequestURI)

	var creds data.Credentials

	err := creds.FromJSON(r.Body)

	if err != nil {
		http.Error( w, "Error Deserializing credentials" , http.StatusBadRequest)
		return
	}

	//TODO:: ReWRITE FROM DB
	//var id string
	//expectedPassword , ok  := users[creds.Username]
	//if !ok  || expectedPassword != creds.Password{
	//	http.Error( w , "Wrong Username or Password" , http.StatusBadRequest)
	//	return
	//}

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

func createToken( name string , password string ) ( string , error ){
	claims := jwt.MapClaims{}
	claims["authorized"]=true
	claims["name"]=name
	claims["password"]=password
	claims["exp"]=time.Now().Add(time.Hour * 1 ).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

