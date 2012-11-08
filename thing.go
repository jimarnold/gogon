package main

import "math"

type Thing struct {
  location Vector
  direction Vector
  size float64
  targetSize float64
}

func(this *Thing) update(elapsed float64) {
  if this.size == 0 {
    return
  }
  this.grow(elapsed * 100)

  newX := this.location.x + (elapsed * this.direction.x * 100)
  newY := this.location.y + (elapsed * this.direction.y * 100)
  this.location = wrapped(Vector{newX,newY})
}

func wrapped(target Vector) Vector {
  wrap := func (i float64, min float64, max float64) float64 {
    result := i
    if result > max {
      result -= max
    }
    if result < min {
      result += max
    }
    return result
  }
  return Vector{wrap(target.x, 0, width), wrap(target.y, 0, height)}
}

func(this *Thing) intersects(other *Thing) bool {
  distance := math.Sqrt(math.Pow((other.location.y - this.location.y), 2) + math.Pow((other.location.x - this.location.x),2) )
  return (this.size + other.size) >= distance
}

func(this *Thing) absorb(other *Thing) {
  if this.targetSize == 0 {
    return
  }
  this.targetSize += other.size
  other.size = 0
}

func(this *Thing) grow(amount float64) {
  if math.Abs(this.size - this.targetSize) < amount {
    return
  }

  if this.size < this.targetSize {
    this.size += amount
  }

  if this.size > this.targetSize {
    this.size -= amount
  }
}
