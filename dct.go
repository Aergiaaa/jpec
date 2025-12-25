package main

import (
	"math"
)

const (
	scale = 8
)

var (
	COS_TABLE [scale][scale]float64
	CoefTable [scale]float64
)

func init() {
	for i := range scale {
		CoefTable[i] = math.Sqrt(2.0 / float64(scale))
	}
	CoefTable[0] = 1 / math.Sqrt(float64(scale))

	for y := range scale {
		for x := range scale {
			COS_TABLE[y][x] = math.Cos(float64((2*x+1)*y) * math.Pi / float64(2*scale))
		}
	}
}

type Block [scale][scale]int        // 8x8 block of bytes // raw data
type DCTBlock [scale][scale]float64 // 8x8 block of float32 // transformed data

func DctImg(cb ColorBlock) DCTColorBlock {
	var yDCT, cbDCT, crDCT DCTBlockOfImg
	for _, v := range []*DCTBlockOfImg{&yDCT, &cbDCT, &crDCT} {
		*v = make(DCTBlockOfImg, len(cb.Y))
	}
	for by := range cb.Y {
		for _, v := range []*DCTBlockOfImg{&yDCT, &cbDCT, &crDCT} {
			(*v)[by] = make([]DCTBlock, len(cb.Y[by]))
		}

		for bx := range cb.Y[by] {
			yDCT[by][bx] = dct(cb.Y[by][bx])
			cbDCT[by][bx] = dct(cb.Cb[by][bx])
			crDCT[by][bx] = dct(cb.Cr[by][bx])
		}
	}

	return DCTColorBlock{
		Y: yDCT, Cb: cbDCT, Cr: crDCT,
		ow: cb.ow, oh: cb.oh, pw: cb.pw, ph: cb.ph,
	}
}

func dct(input Block) (output DCTBlock) {
	for u := range scale {
		for v := range scale {
			sum := 0.0
			for x := range scale {
				for y := range scale {

					sum += float64(input[x][y]) * COS_TABLE[u][x] * COS_TABLE[v][y]

				}
			}

			output[u][v] = CoefTable[u] * CoefTable[v] * sum
		}
	}

	return output
}

func IDctImg(cb DCTColorBlock) ColorBlock {
	var yIDCT, cbIDCT, crIDCT BlockOfImg
	for _, v := range []*BlockOfImg{&yIDCT, &cbIDCT, &crIDCT} {
		*v = make(BlockOfImg, len(cb.Y))
	}
	for by := range cb.Y {
		for _, v := range []*BlockOfImg{&yIDCT, &cbIDCT, &crIDCT} {
			(*v)[by] = make([]Block, len(cb.Y[by]))
		}

		for bx := range cb.Y[by] {
			yIDCT[by][bx] = idct(cb.Y[by][bx])
			cbIDCT[by][bx] = idct(cb.Cb[by][bx])
			crIDCT[by][bx] = idct(cb.Cr[by][bx])
		}
	}

	return ColorBlock{
		Y: yIDCT, Cb: cbIDCT, Cr: crIDCT,
		ow: cb.ow, oh: cb.oh, pw: cb.pw, ph: cb.ph,
	}
}

func idct(input DCTBlock) (output Block) {
	for x := range scale {
		for y := range scale {
			sum := 0.0

			for u := range scale {
				for v := range scale {
					sum += CoefTable[u] * CoefTable[v] * input[u][v] * COS_TABLE[u][x] * COS_TABLE[v][y]
				}
			}

			output[x][y] = int(math.Round(sum)) // rounding
		}
	}

	return output
}
