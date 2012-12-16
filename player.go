package main

import (
	"math"
	"github.com/go-gl/glfw"
)

type Player struct {
	Thing
}

const speed = 10 
const maxSpeed float64 = 40.0
const brake = .95
const initialSize = 24

func (this *Player) update(elapsed float64) {
	this.Thing.update(elapsed)
	if(keyDown(KeyW)) {
      this.up(elapsed)
    }
	if(keyDown(KeyS)) {
      this.down(elapsed)
    }
	if(keyDown(KeyA)) {
      this.left(elapsed)
    }
	if(keyDown(KeyD)) {
      this.right(elapsed)
    }

	this.direction = this.direction.scale(brake)
}

func(this *Thing) absorb(other Element) {
	if this.isDead() {
		panic("Dead things can't absorb!")
	}
	this.targetSize += other.Size()
}

func (this * Player) burst() {
	if this.targetSize > initialSize {
		debugf("size before: %f", this.targetSize)
		this.targetSize -= math.Min(this.size - initialSize, this.targetSize / 2)
		debugf("size after: %f", this.targetSize)
	}
}

func(this *Player) up(elapsed float64) {
	this.thrust(Vector2{0,speed * elapsed})
}

func(this *Player) down(elapsed float64) {
	this.thrust(Vector2{0,-speed * elapsed})
}

func(this *Player) left(elapsed float64) {
	this.thrust(Vector2{-speed * elapsed,0})
}

func(this *Player) right(elapsed float64) {
	this.thrust(Vector2{speed * elapsed,0})
}

func(this *Player) thrust(v Vector2) {
	this.direction = this.direction.Add(v).clampedTo(maxSpeed)
}

func keyDown(key int) bool {
	return glfw.Key(key) == 1
}
