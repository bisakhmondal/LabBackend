package handlers

import (
	"auth/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type UpdateH struct{
	db *data.MongoClient
	l *log.Logger
}

func NewUpdateH(db * data.MongoClient,l *log.Logger) *UpdateH{
	return &UpdateH{
		db: db,
		l: l,
	}
}

//Update Route.
func (p *UpdateH)Update(rw http.ResponseWriter, r *http.Request){
	var user data.Person
	err := user.FromJSON(r.Body)
	//user.ID = bson.M{"$oid":}
	if err !=nil{
		http.Error(rw,"Invalid Request to Update", http.StatusBadRequest)
		p.l.Println("Invalid Request to Update", err)
		return
	}

	//id,err := ParseCookie(r)
	id, err := primitive.ObjectIDFromHex("5f5cd403a819ad84f8cdfc97")

	if err !=nil{
		http.Error(rw,"Invalid Cookie ReLOGIN", http.StatusBadRequest)
		return
	}
	user.ID = id //id

	err = p.db.UpdateDB(&user)

	if err!=nil{
		http.Error(rw, "Unable to Update", http.StatusBadRequest)
		return
	}
        rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
}

//Update Image through form data.

func (p *UpdateH)UploadImage(rw http.ResponseWriter, r* http.Request){
	err := r.ParseMultipartForm(1024*1024)

	if err!=nil{
		http.Error(rw, "Unable to Parse Form data", http.StatusBadRequest)
		return
	}

	// will check for with frontend..
	//id,err := ParseCookie(r)
	id, err := primitive.ObjectIDFromHex("5f5cd403a819ad84f8cdfc97")

	if err !=nil{
		http.Error(rw, "Unable to Parse Cookie", http.StatusNetworkAuthenticationRequired)
		p.l.Println(err)
		return
	}
	image,_,err := r.FormFile("file")

	if err!=nil{
		http.Error(rw, "Invalid file format", http.StatusBadRequest)
		p.l.Println(err)
		return
	}

	strImg, err := ParseImage(image)
	if err !=nil{
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		p.l.Println(err)
		return
	}
	//log.Println("encoded String", len(strImg))

	var user data.Person
	user.ID=id
	user.PROFILE= strImg

	p.db.UpdateDB(&user)
        w.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
}
