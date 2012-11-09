package main

import "math/rand"
import "time"

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}

func random(min, max float64) float64 {
  return rand.Float64() * (max - min) + min
}

type Vector struct {
  x,y float64
}

func(v1 Vector) Add(v2 Vector) Vector {
  return Vector{v1.x + v2.x, v2.y + v2.y}
}

func(v Vector) Mult(s float64) Vector {
  return Vector{v.x * s, v.y * s}
}

