package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	NAME string`bson:"name" json:"name"`
	IMG string `bson:"img, omitempty" json:"img"`
}

type Projects []*Project

type Person struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NAME string `bson:"name" json:"name" validate:"required"`
	EMAIL string `bson:"email" json:"email" validate:"email,required"`
	EDUCATION []string `bson:"education,omitempty" json:"education"`
	USERNAME string `bson:"username" json:"username"`
	PASSWORD string `bson:"password" json:"-"`//hash String
	ROUTE string `bson:"route" json:"route"`
	SPECIALIZATION []string `bson:"specialization,omitempty" json:"specialization"` //Specialized Field
	PROJECTS Projects `bson:"projects,omitempty" json:"projects"`
	ACHIEVEMENTS []string `bson:"achievements,omitempty" json:"achievements"`
}

