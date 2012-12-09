package main

import (
	"github.com/go-gl/glfw"
)

type Player struct {
  Thing
}

const speed = 0.1
const maxSpeed float64 = 2.0
const brake = .95

func (this *Player) update(elapsed float64) {
  this.Thing.update(elapsed)
	if(keyDown(KeyW)) {
      this.up()
    }
	if(keyDown(KeyS)) {
      this.down()
    }
	if(keyDown(KeyA)) {
      this.left()
    }
	if(keyDown(KeyD)) {
      this.right()
    }

  this.direction = this.direction.scale(brake)
}

func(this *Player) up() {
  this.thrust(Vector2{0,-speed})
}

func(this *Player) down() {
  this.thrust(Vector2{0,speed})
}

func(this *Player) left() {
  this.thrust(Vector2{-speed,0})
}

func(this *Player) right() {
  this.thrust(Vector2{speed,0})
}

func(this *Player) thrust(v Vector2) {
  this.direction = this.direction.Add(v).clampedTo(maxSpeed)
}

func keyDown(key int) bool {
  return glfw.Key(key) == 1
}
