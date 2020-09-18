package handlers

import (
    "log"
    "auth/data"
    "net/http"
    "os"
    "github.com/dgrijalva/jwt-go"
)

type Util struct {
	db *data.MongoClient
	l *log.Logger
}

func NewUtil(db * data.MongoClient,l *log.Logger) *SignIn{
	return &SignIn{
		db: db,
		l: l,
	}
}


func check_valid_authentication( c *http.Cookie , w http.ResponseWriter ) bool {
    sessionToken := c.Value
	claims := jwt.MapClaims{}
	_ , errr := jwt.ParseWithClaims(sessionToken,  claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})


	if  errr != nil{
		http.Error(w , "Wrong token signature",  http.StatusBadRequest)
		return false
    }
    return true
}

func (u *Util)Create( rw http.ResponseWriter , r *http.Request ){
    c, err := r.Cookie("session_token")
    if err!=nil{
        http.Error( rw , "Unauthorized" , http.StatusBadRequest)
        return 
    }

    //check if authentication is valid
    if check_valid_authentication(c , rw ){
        http.Error( rw , "Unauthorized" , http.StatusBadRequest)
        return 
    }

    // TODO

}