package handlers

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"serving-api/data"
	"strings"
)

type PersonH struct{
	dbClient *data.MongoClient
	l *log.Logger
}

func NewPersonH(db *data.MongoClient, l*log.Logger) *PersonH{
	return &PersonH{
		dbClient: db,
		l: l,
	}
}


func (p *PersonH)GetData(rw http.ResponseWriter, r* http.Request) {

	personData, err := p.dbClient.GetData()

	if err!=nil {
		http.Error(rw,"Unable to Fetch Data", http.StatusInternalServerError)
		p.l.Println("Can't Fetch from DB ",err)
		return
	}

	err = personData.ToJSON(rw)

	if err!=nil{
		http.Error(rw,"Unable to Process Data", http.StatusInternalServerError)
		p.l.Println("Can't Marshal data")
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type","application/json")
}

func (p* PersonH)FetchUser(rw http.ResponseWriter, r* http.Request){
	vars := mux.Vars(r)
	route := strings.ToLower(vars["route"])

	user,err := p.dbClient.FindUser(&route)

	if err!=nil{
		http.Error(rw, "Bad Route", http.StatusBadGateway)
		return
	}
	err = user.ToJSON(rw)
	if err!=nil{
		http.Error(rw, "Unable to Process data", http.StatusInternalServerError)
		p.l.Println("unable to marshal json")
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
}
