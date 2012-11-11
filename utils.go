package main

import "math"
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

type RGB struct {r,g,b float64}

func(v1 Vector) Add(v2 Vector) Vector {
  return Vector{v1.x + v2.x, v1.y + v2.y}
}

func(v1 Vector) Subtract(v2 Vector) Vector {
  return Vector{v1.x - v2.x, v1.y - v2.y}
}

func(v Vector) Mult(s float64) Vector {
  return Vector{v.x * s, v.y * s}
}

func(v Vector) Length() float64 {
  return math.Sqrt((v.x * v.x) + (v.y * v.y))
}

func(this *Vector) Normalize() {
  length := this.Length()
  this.x = this.x / length
  this.y = this.y / length
}

func(this Vector) DistanceTo(other Vector) float64 {
  x := other.x - this.x
  y := other.y - this.y
  return math.Sqrt((x*x) + (y*y))
}
