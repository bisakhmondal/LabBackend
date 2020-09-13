package main

import (
	"github.com/spf13/viper"
	// "fmt"
	"log"
)


// func main(){
// 	fmt.Println(getURI("MONGO"))
// }

func getURI(key string) string{
	viper.SetConfigFile("../serving-api/config.yaml")

	err:= viper.ReadInConfig()
	
	if err !=nil {
		log.Fatal("Can't Read Env Variable",err)
	}

	value,ok := viper.Get(key).(string)
	
	if !ok{
		log.Fatal("Invalid TypeCast")
	}

	return value
}
