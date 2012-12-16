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
	this.player = &Player{Thing{location : Vector2{width / 2, height / 2}, targetSize : 10, size : 10}}
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
	size := random(6, 12)
	location := Vector2{width, random(0,1) * height}
	direction := Vector2{random(-2,-0.5), random(-0.1,0.1)}
	e := NewThing(location, direction, size)
	return &e
}

func(game *Game) collide() {
	game.elements.Each(func(i int, a Element) {
		game.elements.Each(func(j int, b Element) {
			if j <= i {
				return
			}
			if a.intersects(b) {
				if a.biggerThan(b) {
					a.absorb(b)
				} else {
					b.absorb(a)
				}
			}
		})
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

type Elements struct {
	items []Element
}

func(this *Elements) Add(e Element) {
	this.items = append(this.items, e)
}

func(this *Elements) Delete(e Element) {
	i := this.IndexOf(e)
	if i > -1 {
		this.items = append(this.items[:i], this.items[i+1:]...)
	}
}

func(this *Elements) IndexOf(e Element) int {
	for i,el := range this.items {
		if el == e {
			return i
		}
	}
	debugf("!Could not find element %v", e)
	return -1
}

type ElementAction func(int, Element)
type ElementQuery func(Element) bool

func(this *Elements) Each(action ElementAction) {
	for i,e := range this.items {
		action(i,e)
	}
}

func(this *Elements) Any(query ElementQuery) bool {
	for _,e := range this.items {
		if query(e) {
			return true
		}
	}
	return false
}

func(this *Elements) Count() int {
	return len(this.items)
}
