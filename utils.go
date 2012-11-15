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
    return max
  } else if f < min {
    return min
  }
  return f
}

type Rect struct {
  top, bottom, left, right float64
}

func(this Rect) cannotIntersect(other Rect) bool {
  if (this.left > other.right || this.right < other.left || this.top < other.bottom || this.bottom > other.top) {
    return true
  }
  return false
}
