package handlers

import (
	"bytes"
	"encoding/base64"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
)

//Image Parser to base64 encoding
func ParseImage(r io.ReadCloser) (string,error){
	image,_,err :=image.Decode(r)

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
		return "", err
	}
	//Encoding the binary to base64
	encoded := base64.StdEncoding.EncodeToString(byt.Bytes())

	return encoded, nil
}
