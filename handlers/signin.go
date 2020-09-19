package handlers

import (
	"auth/data"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"strings"
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
	l.l.Println(creds)
	if err != nil {
		http.Error( w, "Error Deserializing credentials" , http.StatusBadRequest)
		return
	}

	

	filter := &bson.M{
		"username":bson.M{
			"$eq":creds.Username,
		}}

	var user data.Person

	//Get User Details
	err = l.db.CheckAuth( &user, filter)
	if err!=nil{
		http.Error(w ,"Not in Database" , http.StatusBadRequest)
		return
	}

	//log.Println(user)

	//Check Password
	//err = bcrypt.CompareHashAndPassword([]byte(user.PASSWORD), []byte(creds.Password))
	//
	//if err!=nil{
	//	http.Error(w, "Incorrect Password", http.StatusBadRequest)
	//	return
	//}
	if user.PASSWORD == creds.Password{
		l.l.Println("Password Match")
	}

	//All Validated: Generate Cookie
	id := strings.Split(user.ID.String(),"(\"")[1]
	id = id[:len(id)-2]

	sessionToken , err := createToken(id, user.USERNAME)
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
		HttpOnly: true,
		SameSite: 4, //SameSiteNone : https://golang.org/src/net/http/cookie.go
	})
	
    
	
	w.WriteHeader(http.StatusOK)
	//For testing
	w.Write([]byte(user.USERNAME))
}

func createToken( id string , username string ) ( string , error ){
	claims := jwt.MapClaims{}
	claims["authorized"]=true
	claims["id"]=id
	claims["username"]=username
	claims["expires"]=time.Now().Add(time.Hour * 1 ).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}
