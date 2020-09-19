package main


import "fmt"

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

func main(){
	fmt.Println(SameSiteNoneMode);
}
