package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	// "io"
)

func main() {
	width, height := getImageDimension("n.png")
	fmt.Println("Width:", width, "Height:", height)


	//fil,_ := os.Open("n.png")
	//fl2 ,_ := os.Create("nn.png")//os.Open("nn.png")


	// Shuvayan test
	f, _ := os.Open( "n.png")
	readeree := bufio.NewReader(f)

	content , _ := ioutil.ReadAll(readeree)

	encoded := base64.StdEncoding.EncodeToString(content)
	fmt.Println("Previous Encoded : " , len(encoded))

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

	//Only Compression
	//byt ,_ := ioutil.ReadFile("n.png")
	////Compress BYTES
	//wrt := new(bytes.Buffer)
	//gw,_ := gzip.NewWriterLevel(wrt, gzip.BestCompression)
	//gw.Write(byt)
	//newencoded := base64.StdEncoding.EncodeToString(wrt.Bytes())
	//fmt.Println("New Encoding", len(newencoded),"saved len: ",len(encoded)-len(newencoded))


	fil ,_ := os.Open("n.png")
	//Resize + Compression
	imgD,_,_ := image.Decode(fil)
	imgD =resize.Resize(768,768,imgD,resize.Lanczos2)
	newByt := new(bytes.Buffer)

	jpeg.Encode(newByt,imgD,nil)
	// ff , _ := os.Create("resizeN.png")
	// io.Copy(ff,bytes.NewReader(newByt.Bytes()))
	//Compress BYTES
	wrt := new(bytes.Buffer)
	gw,_ := gzip.NewWriterLevel(wrt, gzip.BestCompression)
	gw.Write(newByt.Bytes())
	newencoded := base64.StdEncoding.EncodeToString(wrt.Bytes())
	fmt.Println("New Encoding", len(newencoded),"saved len: ",len(encoded)-len(newencoded))
	// fmt.Println(newencoded)
	// reader := bytes.NewReader(byt)
	// n,_ := io.Copy(fl2,reader)
	// fmt.Println(n)
	// fl2.Close()
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
