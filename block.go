package main

import (
	"image"
	"image/color"
)

type BlockOfImg [][]Block // matrix of 8x8 blocks
type DCTBlockOfImg [][]DCTBlock

type ColorBlock struct {
	Y, Cb, Cr BlockOfImg
	ow, oh    int //original image width and height
	pw, ph    int //padded image width and height
}

type DCTColorBlock struct {
	Y, Cb, Cr DCTBlockOfImg
	ow, oh    int //original image width and height
	pw, ph    int //padded image width and height
}

func imgToBlocks(img image.Image) ColorBlock {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	bw := (w + scale - 1) / scale // number of blocks in width
	bh := (h + scale - 1) / scale // number of blocks in height

	pw := bw * scale // padded width
	ph := bh * scale // padded height

	var y, cb, cr BlockOfImg
	for _, v := range []*BlockOfImg{&y, &cb, &cr} {
		*v = *createBlockGrid(bh, bw, pw, ph)
	}

	for by := range bh {
		for bx := range bw {
			var yb, cbb, crb Block
			for py := range scale {
				for px := range scale {
					x := bx*scale + px
					y := by*scale + py

					yc := getPixel(x, y, w, h, img)
					// Convert from uint32 [0, 65535] to byte [0, 255], then to int16 [-128, 127]
					yb[py][px] = int(yc.Y) - 128
					cbb[py][px] = int(yc.Cb) - 128
					crb[py][px] = int(yc.Cr) - 128
				}
			}
			y[by][bx] = yb
			cb[by][bx] = cbb
			cr[by][bx] = crb
		}
	}

	return ColorBlock{
		Y: y, Cb: cb, Cr: cr,
		ow: w, oh: h, pw: pw, ph: ph,
	}
}

func blocksToImg(cb ColorBlock, scale int) image.Image {
	pw := cb.pw
	ph := cb.ph
	img := image.NewRGBA(image.Rect(0, 0, pw, ph))

	for by := range cb.Y {
		for bx := range cb.Y[by] {
			yBlock := cb.Y[by][bx]
			cbBlock := cb.Cb[by][bx]
			crBlock := cb.Cr[by][bx]
			for py := range scale {
				for px := range scale {
					x := bx*scale + px
					y := by*scale + py

					yc := color.YCbCr{
						Y:  uint8(clamp(yBlock[py][px] + 128)),
						Cb: uint8(clamp(cbBlock[py][px] + 128)),
						Cr: uint8(clamp(crBlock[py][px] + 128)),
					}

					img.Set(x, y, yc)
				}
			}
		}
	}

	return img.SubImage(image.Rect(0, 0, cb.ow, cb.oh))
}

func createBlockGrid(h, w int, pw, ph int) *BlockOfImg {
	grid := make(BlockOfImg, h)

	for y := range h {
		grid[y] = make([]Block, w)
	}

	return &grid
}

func getPixel(x, y, w, h int, img image.Image) color.YCbCr {
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

	return color.YCbCrModel.Convert(img.At(x, y)).(color.YCbCr)
}
