package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

)

func MiddlewareAuthenticate( next http.Handler ) http.Handler{
	return http.HandlerFunc(func(rw http.ResponseWriter , r *http.Request){


		c , err := r.Cookie("session_token")

		if err != nil {
			if err == http.ErrNoCookie{
				http.Error( rw ,"Unauthorized" , http.StatusUnauthorized)
				return
			}
			http.Error(rw ,"Bad Request" , http.StatusBadRequest )
			return
		}

		sessionToken := c.Value
		claims := jwt.MapClaims{}
		_ , err = jwt.ParseWithClaims(sessionToken,  claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("API_SECRET")), nil
		})

		if time.Now().Unix() > claims["expires"].(int64){
			http.Error(rw ,"Unauthorized" , http.StatusBadRequest )
			return
		}

		next.ServeHTTP(rw ,r)
	})
}

func CorsMiddleware( next http.Handler ) http.Handler{
	return http.HandlerFunc( func(rw  http.ResponseWriter , r *http.Request ){
		origin := r.Header.Get("Origin");

		rw.Header().Set("Access-Control-Allow-Origin", origin )
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")

		headers := rw.Header()
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
	})
}

// Always Fetch the userID from http-only cookie.
func (p *UpdateH)ParseCookie(r * http.Request) (*primitive.ObjectID){
	// c, err := r.Cookie("session_token")
	c := r.Header["Authorization"]

	if len(c) == 0  {
		return nil
	}
	sessionToken := strings.Split(c[0], " ")[1]

	claims := jwt.MapClaims{}
	_ , err := jwt.ParseWithClaims(sessionToken,  claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err!=nil{
		return nil
	}
	id, err := primitive.ObjectIDFromHex(claims["id"].(string))

	if err !=nil{
		return nil
	}

	return &id
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
