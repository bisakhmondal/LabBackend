package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Project struct {
	NAME string`bson:"name,omitempty" json:"name"`
	IMG string `bson:"img,omitempty" json:"img"`
}

type Projects []*Project

type Person struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NAME string `bson:"name,omitempty" json:"name" validate:"required"`
	EMAIL string `bson:"email,omitempty" json:"email" validate:"email,required"`
	EDUCATION []string `bson:"education,omitempty" json:"education"`
	USERNAME string `bson:"username,omitempty" json:"username"`
	PASSWORD string `bson:"password,omitempty" json:"-"`//hash String
	ROUTE string `bson:"route,omitempty" json:"route"`
	PROFILE string `bson:"profile,omitempty" json:"profile"`
	SPECIALIZATION []string `bson:"specialization,omitempty" json:"specialization"` //Specialized Field
	PROJECTS Projects `bson:"projects,omitempty" json:"projects"`
	ACHIEVEMENTS []string `bson:"achievements,omitempty" json:"achievements"`
}

func (p *Person) FromJSON(r io.Reader)error{
	en := json.NewDecoder(r)
	return en.Decode(p)
}

func (p Person) ToBSON() (*bson.M,error){
	databyte, err := bson.Marshal(p)
	if err!=nil{
		return nil,err
	}
	var update bson.M

	err = bson.Unmarshal(databyte, &update)

	if err!=nil{
		return nil,err
	}
	return &update, nil
}
