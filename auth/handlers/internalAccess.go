package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"time"
)

func MiddlewareAuthenticate( next http.Handler ) http.Handler{
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

// Always Fetch the userID from http-only cookie.
func ParseCookie(r * http.Request) (*primitive.ObjectID,error){
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil,err
	}
	sessionToken := c.Value
	claims := jwt.MapClaims{}
	_ , err = jwt.ParseWithClaims(sessionToken,  claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err!=nil{
		return nil,err
	}
	id, err := primitive.ObjectIDFromHex(claims["id"].(string))

	if err !=nil{
		return nil, err
	}

	return &id,nil
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	log.Printf("GET %s\n" , r.RequestURI)
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
	newSessionToken , err := createToken( claims["id"].(string) , claims["username"].(string))
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

