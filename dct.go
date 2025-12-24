package main

import "math"

const (
	scale = 8
)

type Block [scale][scale]int        // 8x8 block of bytes // raw data
type DCTBlock [scale][scale]float64 // 8x8 block of float32 // transformed data

func dct(input Block) (output DCTBlock) {
	var cU, cV float64

	for u := range scale - 1 {
		for v := range scale - 1 {
			sum := 0.0

			for x := range scale - 1 {
				for y := range scale - 1 {
					cosX := math.Cos(degCos(x, u, scale))
					cosY := math.Cos(degCos(y, v, scale))
					sum += float64(input[x][y]) * cosX * cosY
				}
			}

			if u == 0 {
				cU = 1 / math.Sqrt(float64(scale))
			} else {
				cU = math.Sqrt(2 / float64(scale))
			}

			if v == 0 {
				cV = 1 / math.Sqrt(float64(scale))
			} else {
				cV = math.Sqrt(2 / float64(scale))
			}

			output[u][v] = cU * cV * sum
		}
	}

	return output
}

func idct(input DCTBlock) (output Block) {
	var cU, cV float64

	for x := range scale - 1 {
		for y := range scale - 1 {
			sum := 0.0

			for u := range scale - 1 {
				for v := range scale - 1 {
					if u == 0 {
						cU = 1 / math.Sqrt(float64(scale))
					} else {
						cU = math.Sqrt(2 / float64(scale))
					}
					if v == 0 {
						cV = 1 / math.Sqrt(float64(scale))
					} else {
						cV = math.Sqrt(2 / float64(scale))
					}

					cosX := math.Cos(degCos(x, u, scale))
					cosY := math.Cos(degCos(y, v, scale))
					sum += cU * cV * input[u][v] * cosX * cosY
				}
			}

			output[x][y] = int(sum + 0.5) // rounding
		}
	}

	return output
}

func degCos(loopIn, loopOut, scale int) float64 {
	return float64((2*loopIn+1)*loopOut) * math.Pi / float64((2 * scale))
}
