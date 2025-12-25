package main

import "math"

var (
	LUM_TABLE = Block{
		{16, 11, 10, 16, 24, 40, 51, 61},
		{12, 12, 14, 19, 26, 58, 60, 55},
		{14, 13, 16, 24, 40, 57, 69, 56},
		{14, 17, 22, 29, 51, 87, 80, 62},
		{18, 22, 37, 56, 68, 109, 103, 77},
		{24, 35, 55, 64, 81, 104, 113, 92},
		{49, 64, 78, 87, 103, 121, 120, 101},
		{72, 92, 95, 98, 112, 100, 103, 99},
	}

	CHRO_TABLE = Block{
		{17, 18, 24, 47, 99, 99, 99, 99},
		{18, 21, 26, 66, 99, 99, 99, 99},
		{24, 26, 56, 99, 99, 99, 99, 99},
		{47, 66, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
		{99, 99, 99, 99, 99, 99, 99, 99},
	}
)

func Quantize(cb DCTColorBlock, factor float64) DCTColorBlock {
	var yq, cbq, crq DCTBlockOfImg
	for _, v := range []*DCTBlockOfImg{&yq, &cbq, &crq} {
		*v = make(DCTBlockOfImg, len(cb.Y))
	}

	lumq := ScaledQuantTable(LUM_TABLE, 100)
	chroq := ScaledQuantTable(CHRO_TABLE, factor)

	for by := range cb.Y {
		for _, v := range []*DCTBlockOfImg{&yq, &cbq, &crq} {
			(*v)[by] = make([]DCTBlock, len(cb.Y[by]))
		}

		for bx := range cb.Y[by] {
			yq[by][bx] = quantBlock(cb.Y[by][bx], lumq)
			cbq[by][bx] = quantBlock(cb.Cb[by][bx], chroq)
			crq[by][bx] = quantBlock(cb.Cr[by][bx], chroq)
		}
	}

	return DCTColorBlock{
		Y: yq, Cb: cbq, Cr: crq,
		ow: cb.ow, oh: cb.oh, pw: cb.pw, ph: cb.ph,
	}
}

func quantBlock(b DCTBlock, qt Block) (out DCTBlock) {
	for y := range scale {
		for x := range scale {
			out[y][x] = math.Round((b[y][x]) / float64(qt[y][x]))
		}
	}

	return out
}

func ScaledQuantTable(qt Block, factor float64) (sqt Block) {
	scale := scaling(factor)

	for y := range qt {
		for x := range qt {
			sqt[y][x] = clamp((int(scale)*qt[y][x] + 50) / 100)
		}
	}

	return sqt
}

func scaling(factor float64) (scale float64) {
	if factor < 50 {
		scale = 5000 / factor
	} else {
		scale = 200 - factor*2
	}
	return scale
}
