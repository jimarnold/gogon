package main

import (
	"math"
)

type Element interface {
	update(elapsed float64)
	intersects(other Element) bool
	grow(amount float64)
	biggerThan(other Element) bool
	die()
	isDead() bool
	Location() Vector2
	getDirection() Vector2
	Size() float64
	Color() Color4f
}

type Thing struct {
	location   Vector2
	direction  Vector2
	size       float64
	targetSize float64
	color      Color4f
}

func NewThing(location, direction Vector2, size float64) *Thing {
	return &Thing{location, direction, size, size, Color4f{1, 0, 0, 1}}
}

func (this *Thing) biggerThan(other Element) bool {
	return this.Size() > other.Size()
}

func (this *Thing) die() {
	this.size = 0
}

func (this *Thing) isDead() bool {
	return this.size == 0
}

func (this *Thing) Location() Vector2 {
	return this.location
}

func (this *Thing) getDirection() Vector2 {
	return this.direction
}

func (this *Thing) Size() float64 {
	return this.size
}

func (this *Thing) Color() Color4f {
	return this.color
}

func (this *Thing) update(elapsed float64) {
	if this.isDead() {
		return
	}
	this.grow(elapsed * 100)

	x := this.location.x + (elapsed * this.direction.x * 100)
	y := this.location.y + (elapsed * this.direction.y * 100)
	if outOfBounds(this.size, this.location.x, this.location.y) {
		this.die()
	}
	this.location = (Vector2{x, y})
}

func outOfBounds(size, x, y float64) bool {
	out := func(i float64, min float64, max float64) bool {
		if i > max || i < min {
			return true
		}
		return false
	}
	return out(x, 0+size, width) || out(y, 0+size, height-size)
}

func wrapped(target Vector2) Vector2 {
	wrap := func(i float64, min float64, max float64) float64 {
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

func (this *Thing) intersects(other Element) bool {
	distance := this.Location().DistanceTo(other.Location())
	return (this.size + other.Size()) >= distance
}

func (this *Thing) grow(amount float64) {
	if math.Abs(this.size-this.targetSize) < amount {
		return
	}

	if this.size < this.targetSize {
		this.size += amount
	}

	if this.size > this.targetSize {
		this.size -= amount
	}
}
