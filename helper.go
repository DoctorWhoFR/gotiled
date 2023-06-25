package generator_2d

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

/*
LoadImage
(path string) ([]byte, image.Image)

Currently, helper function, used to simply load an image
based on the current+folder + provided path.
--

Params:
  - path string -> the path to your image like ("\\assets\\image.png"), you should always use "\\" instead of "/"

Response:
  - Return two thing , an array of bytes representation of the got image, and also an image.Image linked to the image path provider.

// TODO LoadImage error handling
This function need error handling
*/
func LoadImage(path string) ([]byte, image.Image) {
	// Open the PNG image file
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
	}

	file, err := os.Open(currentDir + path)
	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	// Decode the PNG image
	baseImg, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	// Encode the image as PNG and write it to the buffer
	err = png.Encode(buf, baseImg)
	if err != nil {
		fmt.Println("Failed to encode image:", err)
	}

	// Get the encoded image data as a byte slice
	encodedImage := buf.Bytes()

	return encodedImage, baseImg
}
