package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
)

//export computeHash
func computeHash(pass string) string{
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(pass),
		bcrypt.DefaultCost,
	)
	sHash := string(hash)

	return sHash
}

//export GETHASH
func GETHASH(){
//read Password
	byt, _ := ioutil.ReadFile("hash.txt")

	//fmt.Println(string(byt))
	s:= computeHash(string(byt))

	//fmt.Println("From go ", s)

//Write HASH
	ioutil.WriteFile("hash.txt",[]byte(s),os.FileMode(770))
}

func main(){ }