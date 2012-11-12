package main

import "math/rand"
import "time"

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}

func random(min, max float64) float64 {
  return rand.Float64() * (max - min) + min
}

type RGB struct {r,g,b float64}

func clamp(f float64, min, max float64) float64 {
  if f > max {
    f = max
  } else if f < min {
    f = min
  }
  return f
}
