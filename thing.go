package main

import "math"

type Element interface {
  update(elapsed float64)
  absorb(other Element)
  intersects(other Element)bool
  grow(amount float64)
  biggerThan(other Element) bool
  die()
  isDead() bool
  getLocation() Vector
  getSize() float64
}

type Player struct {
  Thing
}

type Thing struct {
  location Vector
  direction Vector
  size float64
  targetSize float64
}

func(this Thing) biggerThan(other Element) bool {
  return this.size > other.getSize()
}

func(this *Thing) die() {
  this.size = 0
}

func(this Thing) isDead() bool {
  return this.size == 0
}

func(this Thing) getLocation() Vector {
  return this.location
}

func(this Thing) getSize() float64 {
  return this.size
}

func(this *Thing) update(elapsed float64) {
  if this.isDead() {
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

func(this *Thing) intersects(other Element) bool {
  distance := math.Sqrt(math.Pow((other.getLocation().y - this.location.y), 2) + math.Pow((other.getLocation().x - this.location.x),2) )
  return (this.size + other.getSize()) >= distance
}

func(this *Thing) absorb(other Element) {
  if this.isDead() {
    return
  }
  this.targetSize += other.getSize()
  other.die()
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
