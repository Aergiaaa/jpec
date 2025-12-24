package main

import (
	"fmt"
	"image/jpeg"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: ./compress <input_image>")
	}

	path := os.Args[1]

	// fetch image
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening image:", err.Error())
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err.Error())
		return
	}

	// testing dct compression
	block := imgToBlocks(img)
	fmt.Println(block.R.Img[0][0])
	img = blocksToImg(block)

	// save image
	outFile, err := os.Create("output.jpeg")
	if err != nil {
		fmt.Println("Error creating output file:", err.Error())
		return
	}
	defer outFile.Close()

	jpeg.Encode(outFile, img, &jpeg.Options{})
}
