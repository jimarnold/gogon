package main

import "math/rand"
import "time"

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}

func random(min, max float64) float64 {
  return rand.Float64() * (max - min) + min
}

