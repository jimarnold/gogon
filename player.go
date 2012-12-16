package main

import (
	"github.com/go-gl/glfw"
)

type Player struct {
	Thing
}

const speed = 10 
const maxSpeed float64 = 40.0
const brake = .95

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
