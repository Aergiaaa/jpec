package main

import (
  "math"
)

const (
  MAX = 255
)

func energyDist(block DCTBlock, w,h int) float64 {
  out := 0.0
  for y := range h {
    for x := range w {
      out += block[y][x] * block[y][x] 
    }
  }
  return out
}

func MSE(block ColorBlock, before ColorBlock, x,y int) float64 {
  c := 64.0
  sum := 0.0
 
  for i := range block.Y[y][x]{
    for j, v := range block.Y[y][x][i]{
      sum += (float64(v) - float64(before.Y[y][x][i][j])) * (float64(v) - float64(before.Y[y][x][i][j]))
    } 
  }

  out := sum / c
  return out
}

func spars(block DCTColorBlock, x,y int) float64 {
  num := 0
  for i := range block.Y[y][x]{
    for _, v := range block.Y[y][x][i]{
      if v == 0 {
        num++
      }
    } 
  }

  return float64(num) / 64.0
}

func PSNR(MSE float64) float64 {
  return 10 * math.Log10(MAX*MAX/MSE)
}
