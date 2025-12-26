package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 5 {
		println("Usage: ./compress <input_image> <quality_factor> <x_loc> <y_loc>")
	}

	path := os.Args[1]
	factor := os.Args[2]
  atX := os.Args[3]
  atY := os.Args[4]

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

  x, _:= strconv.Atoi(atX)
  y, _:= strconv.Atoi(atY)

	// testing dct compression
	block := imgToBlocks(img)
	fmt.Println("original")
	for _, v := range block.Y[y][x] {
		fmt.Println(v)
	}

	dctBlock := DctImg(block)
	fmt.Println("dct")
	for _, v := range dctBlock.Y[y][x] {
		fmt.Printf("%.3f\n",v)
	}

  fmt.Println()
  total := energyDist(dctBlock.Y[y][x],8,8)
  fmt.Println(total)
  fmt.Println(energyDist(dctBlock.Y[y][x],3,3))
  fmt.Println(energyDist(dctBlock.Y[y][x],0,0)/total)
  fmt.Println()

	factorF, err := strconv.ParseFloat(factor, 64)
	if err != nil {
		fmt.Println("Error parsing quality factor:", err.Error())
	}
	quantBlock := Quantize(dctBlock, factorF)
	fmt.Println("quant")
	for _, v := range quantBlock.Y[y][x] {
		fmt.Println(v)
	}

  fmt.Println(spars(quantBlock,x,y))

	idctBlock := IDctImg(quantBlock)
	fmt.Println("idct")
	for _, v := range idctBlock.Y[y][x] {
		fmt.Println(v)
	}

  mse := MSE(idctBlock,block,x,y)
  fmt.Println(mse)
  fmt.Println(PSNR(mse))

	img = blocksToImg(idctBlock, scale)


  name := fmt.Sprintf("output_%s.jpeg",factor)
	// save image
	out, err := os.Create(name)
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
