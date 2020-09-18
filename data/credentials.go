package data

import (
	"encoding/json"
	"io"
)

type Credentials struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

func (p *Credentials) FromJSON(r io.Reader)error{
	en := json.NewDecoder(r)
	return en.Decode(p)
}

