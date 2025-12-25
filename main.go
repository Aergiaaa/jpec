package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		println("Usage: ./compress <input_image> <quality_factor>")
	}

	path := os.Args[1]
	factor := os.Args[2]
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
	fmt.Println("original")
	for _, v := range block.Y[0][0] {
		fmt.Println(v)
	}

	dctBlock := DctImg(block)
	fmt.Println("dct")
	for _, v := range dctBlock.Y[0][0] {
		fmt.Println(v)
	}

	factorF, err := strconv.ParseFloat(factor, 64)
	if err != nil {
		fmt.Println("Error parsing quality factor:", err.Error())
	}
	quantBlock := Quantize(dctBlock, factorF)
	fmt.Println("quant")
	for _, v := range quantBlock.Y[0][0] {
		fmt.Println(v)
	}

	idctBlock := IDctImg(quantBlock)
	fmt.Println("idct")
	for _, v := range idctBlock.Y[0][0] {
		fmt.Println(v)
	}

	img = blocksToImg(idctBlock, scale)

	// save image
	out, err := os.Create("output.jpeg")
	if err != nil {
		fmt.Println("Error creating output file:", err.Error())
		return
	}
	defer out.Close()

	jpeg.Encode(out, img, &jpeg.Options{Quality: 80})
}

func clamp(value int) int {
	if value < 1 {
		return 1
	} else if value > 255 {
		return 255
	}
	return value
}
