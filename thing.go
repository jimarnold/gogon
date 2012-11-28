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
  Location() Vector2
  getDirection() Vector2
  Size() float64
}

type Thing struct {
  location Vector2
  direction Vector2
  size float64
  targetSize float64
}

func NewThing(location, direction Vector2, size float64) *Thing {
  return &Thing{location, direction, size, size}
}

func(this Thing) biggerThan(other Element) bool {
  return this.Size() > other.Size()
}

func(this *Thing) die() {
  this.size = 0
}

func(this Thing) isDead() bool {
  return this.size == 0
}

func(this Thing) Location() Vector2 {
  return this.location
}

func(this Thing) getDirection() Vector2 {
  return this.direction
}

func(this Thing) Size() float64 {
  return this.size
}

func(this *Thing) update(elapsed float64) {
  if this.isDead() {
    return
  }
  this.grow(elapsed * 100)

  x := this.location.x + (elapsed * this.direction.x * 100)
  y := this.location.y + (elapsed * this.direction.y * 100)
  this.location = wrapped(Vector2{x,y})
}

func wrapped(target Vector2) Vector2 {
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
  return Vector2{wrap(target.x, 0, width), wrap(target.y, 0, height)}
}

func(this *Thing) intersects(other Element) bool {
  distance := this.Location().DistanceTo(other.Location())
  return (this.size + other.Size()) >= distance
}

func(this *Thing) absorb(other Element) {
  if this.isDead() {
    panic("Dead things can't absorb!")
  }
  this.targetSize += other.Size()
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
