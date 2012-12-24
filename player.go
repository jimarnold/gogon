package main

import (
	"math"
	"github.com/go-gl/glfw"
)

type Player struct {
	Thing
	score int
}

const speed = 10 
const maxSpeed float64 = 40.0
const brake = .95
const initialSize = 32

func NewPlayer(location Vector2) *Player {
	return &Player{Thing:Thing{location : location, targetSize : initialSize, size : 16, color: Color4f{0,0,1,1}}, score:0}
}

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

func(this *Player) collideWith(e Element) {
	switch t := e.(type) {
	default:
		debugf("unexpected type %T", t)
	case *Thing:
		this.absorb(e)
	case *Shrinker:
		this.burst()
	case *Pickup:
		this.score += 10
	}
	e.die()
}

func(this *Player) absorb(e Element) {
	if this.isDead() {
		panic("Dead things can't absorb!")
	}
	this.targetSize += e.Size()
}

func (this *Player) burst() {
	if this.targetSize > initialSize {
		this.targetSize -= math.Min(this.size - initialSize, this.targetSize / 2)
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
