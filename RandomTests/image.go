package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"encoding/base64"
	"strconv"
)

func main() {
	width, height := getImageDimension("n.png")
	fmt.Println("Width:", width, "Height:", height)

	//fil,_ := os.Open("n.png")
	fl2 ,_ := os.Create("nn.png")//os.Open("nn.png")


	// Shuvayan test
	f, _ := os.Open( "n.png")
	readeree := bufio.NewReader(f)

	content , _ := ioutil.ReadAll(readeree)

	encoded := base64.StdEncoding.EncodeToString(content)

	fmt.Println("Encoded : " + strconv.Itoa( len(encoded) ))

	//Method 1
	////byt := new(bytes.Buffer)
	//Image,s,_ := image.Decode(fil)
	//fmt.Println(s)
	////op := &jpeg.Options{
	////
	////}
	//jpeg.Encode(fl2,Image,nil)

	//method 2
	//info,_  :=fil.Stat()
	//byt:= make([]byte,info.Size())
	byt ,_ := ioutil.ReadFile("n.png")
	reader := bytes.NewReader(byt)
	 n,_ := io.Copy(fl2,reader)
	 fmt.Println(n)
	 fl2.Close()
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return img.Width, img.Height
}
