package main

import "math"

type Vector struct {
  x,y float64
}

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

func(v Vector) Normalized() Vector {
  f := 1.0 / v.Length()
  return Vector{v.x * f, v.y * f}
}

func(v1 Vector) DistanceTo(v2 Vector) float64 {
  x := v2.x - v1.x
  y := v2.y - v1.y
  return math.Sqrt((x*x) + (y*y))
}

func(v Vector) clampedTo(max float64) Vector {
  v.x = clamp(v.x, -max, max)
  v.y = clamp(v.y, -max, max)
  return v
}

func(v Vector) scale(f float64) Vector {
  length := v.Length()
  if length == 0 {
    return v
  }
  vn := v.Normalized()
  return vn.Mult(length * f)
}

