package main

import (
	"image"
	"image/color"
)

type BlockOfImg struct {
	Img   [][]Block // matrix of 8x8 blocks
	width int
}

type ColorBlock struct {
	R, G, B BlockOfImg
}

func imgToBlocks(img image.Image) ColorBlock {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	bw := (w + scale - 1) / scale // number of blocks in width
	bh := (h + scale - 1) / scale // number of blocks in height

	var rb, gb, bb BlockOfImg
	for _, v := range []*BlockOfImg{&rb, &gb, &bb} {
		*v = *createBlockGrid(bh, bw)
	}

	for by := range bh {
		for bx := range bw {
			var r, g, b Block
			for py := range scale {
				for px := range scale {
					x := bx*scale + px
					y := by*scale + py

					rr, gg, bb, _ := getPixel(x, y, w, h, img).RGBA()
					// Convert from uint32 [0, 65535] to byte [0, 255], then to int16 [-128, 127]
					r[py][px] = int(rr>>8) - 128
					g[py][px] = int(gg>>8) - 128
					b[py][px] = int(bb>>8) - 128
				}
			}
			rb.Img[by][bx] = r
			gb.Img[by][bx] = g
			bb.Img[by][bx] = b
		}
	}

	return ColorBlock{R: rb, G: gb, B: bb}
}

func blocksToImg(cb ColorBlock) image.Image {
	bh := len(cb.R.Img)
	bw := cb.R.width
	w := bw * scale
	h := bh * scale

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for by := range bh {
		for bx := range bw {
			rBlock := cb.R.Img[by][bx]
			gBlock := cb.G.Img[by][bx]
			bBlock := cb.B.Img[by][bx]
			for py := range scale {
				for px := range scale {
					x := bx*scale + px
					y := by*scale + py

					r := uint8(clamp(rBlock[py][px] + 128))
					g := uint8(clamp(gBlock[py][px] + 128))
					b := uint8(clamp(bBlock[py][px] + 128))

					img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
				}
			}
		}
	}

	return img
}

func clamp(value int) int {
	if value < 0 {
		return 0
	} else if value > 255 {
		return 255
	}
	return value
}

func createBlockGrid(height, width int) *BlockOfImg {
	grid := &BlockOfImg{
		Img:   make([][]Block, height),
		width: width,
	}

	for y := range height {
		grid.Img[y] = make([]Block, width)
	}

	return grid
}

func getPixel(x, y, w, h int, img image.Image) color.Color {
	if x < 0 {
		x = 0
	} else if x >= w {
		x = w - 1
	}

	if y < 0 {
		y = 0
	} else if y >= h {
		y = h - 1
	}

	return img.At(x, y)
}
