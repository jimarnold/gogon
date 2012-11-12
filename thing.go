package main

import "math"
import "github.com/go-gl/glfw"

type Element interface {
  update(elapsed float64)
  absorb(other Element)
  intersects(other Element)bool
  grow(amount float64)
  biggerThan(other Element) bool
  die()
  isDead() bool
  Location() Vector
  getDirection() Vector
  Size() float64
  thrust(Vector)
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
  return this.size > other.Size()
}

func(this *Thing) die() {
  this.size = 0
}

func(this Thing) isDead() bool {
  return this.size == 0
}

func(this Thing) Location() Vector {
  return this.location
}

func(this Thing) getDirection() Vector {
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
  this.location = wrapped(Vector{x,y})
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
  distance := this.Location().DistanceTo(other.Location())
  return (this.size + other.Size()) >= distance
}

func(this *Thing) absorb(other Element) {
  if this.isDead() {
    return
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

const speed float64 = 0.1
const brake = 0.03

func up() {
  player.thrust(Vector{0,-0.1})
}

func down() {
  player.thrust(Vector{0,0.1})
}

func left() {
  player.thrust(Vector{-0.1,0})
}

func right() {
  player.thrust(Vector{0.1,0})
}

func(this *Thing) thrust(v Vector) {
  const maxSpeed float64 = 1.0
  this.direction = this.direction.Add(v)
  if this.direction.x > maxSpeed {
    this.direction.x = maxSpeed
  }
  if this.direction.y > maxSpeed {
    this.direction.y = maxSpeed
  }
  if this.direction.x < -maxSpeed {
    this.direction.x = -maxSpeed
  }
  if this.direction.y < -maxSpeed {
    this.direction.y = -maxSpeed
  }
}
type Key int
const KeyA Key = 65
const KeyS Key = 83
const KeyD Key = 68
const KeyW Key = 87

type KeyAction func()()
var keyActions = map[Key] KeyAction {
  KeyW : up,
  KeyS : down,
  KeyA : left,
  KeyD : right,
}
var keysToHandle = []Key {KeyW,KeyS,KeyA,KeyD}

func (this *Player) update(elapsed float64) {
  this.Thing.update(elapsed)

  handleKeys(keysToHandle)

  if this.direction.x < 0 {
    this.direction.x += brake
  }
  if this.direction.x > 0 {
    this.direction.x -= brake
  }
  if this.direction.y < 0 {
    this.direction.y += brake
  }
  if this.direction.y > 0 {
    this.direction.y -= brake
  }
}

func handleKeys(keys []Key) {
  for _,key := range keys {
    if(glfw.Key(int(key)) == 1) {
      keyActions[key]()
    }
  }
}
