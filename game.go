package main

import (
	"github.com/jimarnold/gl"
	"github.com/jimarnold/gltext"
)

type Game struct {
	gameState GameState
	elements *Elements
	player *Player
	program gl.Program
	positionAttrib gl.AttribLocation
	cameraToClipMatrixUniform, colorUniform gl.UniformLocation 
	cameraToClipMatrix Matrix4x4
	vao gl.VertexArray
	text *gltext.Font
	totalTime float64
	elapsedSpawnTime float64
}

func(this *Game) start() {
	if this.gameState == running {
		return
	}
	this.createElements()
    this.gameState = running
}

func(this *Game) createElements() {
	this.player = &Player{Thing{location : Vector2{width / 2, height / 2}, targetSize : initialSize, size : 16}}
	this.elements = &Elements{make([]Element, 0)}
	this.elements.Add(this.player)
	this.elements.Add(createEnemy())
}

func(this *Game) update(elapsed float64) {
	game.totalTime += elapsed
	switch this.gameState {
		case running:
			this.run(elapsed)
		case initialized, won, lost:
			waitForReset()
	}
}

const enemySpawnInterval float64 = 0.5

func(this *Game) run(elapsed float64) {
	game.elapsedSpawnTime += elapsed
	if game.elapsedSpawnTime > enemySpawnInterval {
		game.elapsedSpawnTime = 0
		this.elements.Add(createEnemy())
	}

	this.elements.Each(func(_ int,e Element) {
		e.update(elapsed)
	})

	this.collide()
	this.prune()
	this.win()
}

func(this *Game) prune() {
	deadThings := make([]Element, 0)
	this.elements.Each(func(_ int,e Element) {
		if e.isDead() {
			deadThings = append(deadThings, e)
		}
	})
	for _,e := range deadThings {
		this.elements.Delete(e)
	}
}

func createEnemy() Element {
	size := random(8, 32)
	location := Vector2{width, random(0,1) * height}
	direction := Vector2{random(-3,-0.5), random(-0.1,0.1)}
	if random(0,1) > 0.9 {
		return &Shrinker{Thing{location:location, direction:direction, size:6, targetSize:6}}
	} else {
		return NewThing(location, direction, size)
	}
	panic("wut?")
}

func(game *Game) collide() {
	game.elements.Each(func(i int, e Element) {
		if e == game.player {
			return
		}

		if game.player.intersects(e) {
			switch t := e.(type) {
				default:
				debugf("unexpected type %T", t)
			case *Thing:
				game.player.absorb(e)
				debug("absorb")
			case *Shrinker:
				game.player.burst()
				debug("burst")
			}
			e.die()
		}
	})
}

func(game *Game) win() {
	if game.player.isDead() {
		game.gameState = lost
		return
	}

	if game.elements.Any(func(e Element) bool {
		return e != game.player && !e.isDead()
	}) {
		game.gameState = running
		return
	}
	game.gameState = won
}
