package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main(){
	pass := "Fuck Peeps"

	hash,_ :=bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost)
	fmt.Println(string(hash))
	err := bcrypt.CompareHashAndPassword(hash,[]byte(pass+"a"))
	fmt.Println(err)
}