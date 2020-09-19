package handlers

import (
	"bytes"
	"encoding/base64"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

//Image Parser to base64 encoding
func ParseImage(r io.ReadCloser, filename string) (string,error){
	var image image.Image
	var err error
	//for Png and JPEG seperate reader
	if strings.HasSuffix(filename,"png"){
		image,err =png.Decode(r)
	}else {
		image, err = jpeg.Decode(r)
	}

	if err!=nil{
		return "", err
	}

	//Resize to reduce network overhead.
	reimg := resize.Resize(512, 512, image, resize.Bilinear)

	byt := new(bytes.Buffer)
	err =jpeg.Encode(
		byt,
		reimg,
		nil,
	)
	if err != nil{
		log.Println("err2")
		return "", err
	}
	//Encoding the binary to base64
	encoded := base64.StdEncoding.EncodeToString(byt.Bytes())
	FIL,_ := os.Create("def.txt")

	b := bytes.NewReader([]byte(encoded))
	io.Copy(FIL,b)
	return encoded, nil
}
