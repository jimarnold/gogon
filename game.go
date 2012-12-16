package main

import (
	"github.com/jimarnold/gl"
	"github.com/jimarnold/gltext"
)

type Game struct {
	gameState GameState
	elements []Element
	player *Player
	program gl.Program
	positionAttrib gl.AttribLocation
	cameraToClipMatrixUniform, colorUniform gl.UniformLocation 
  cameraToClipMatrix Matrix4x4
	vao gl.VertexArray
	text *gltext.Font
}

func(this *Game) start() {
	this.createElements()
    this.gameState = running
}

func(this *Game) createElements() {
  this.player = &Player{Thing{location : Vector2{width / 2, height / 2}, targetSize : 10, size : 10}}
  this.elements = make([]Element, 0)
  this.AddElement(this.player)
  this.AddElement(createEnemy())
}

func(this *Game) update(elapsed float64) {
  switch this.gameState {
    case running:
      this.run(elapsed)
    case initialized, won, lost:
      waitForReset()
  }
}

func(this *Game) run(elapsed float64) {
  if random(0,1) > (60 * elapsed) {
    this.elements = append(this.elements, createEnemy())
  }

  for _,e := range this.elements {
    e.update(elapsed)
  }

  this.collide()
	this.prune()
  this.win()
}

func(this *Game) prune() {
	deadThings := make([]Element, 0)
	for _,e := range this.elements {
		if e.isDead() {
			deadThings = append(deadThings, e)
		}
	}
	for _,e := range deadThings {
		this.DeleteElement(e)
	}
}

func createEnemy() Element {
    size := random(6, 12)
    location := Vector2{width, random(0,1) * height}
    direction := Vector2{random(-1,-0.5), random(-0.1,0.1)}
    e := NewThing(location, direction, size)
	return &e
}

func(this *Game) DeleteElement(e Element) {
	i := IndexOf(this.elements, e)
	if i > -1 {
	  this.elements = append(this.elements[:i], this.elements[i+1:]...)
	}
}

type ElementList []Element

func IndexOf(elements []Element, e Element) int {
	for i := 0; i < len(elements); i++ {
    	if (elements[i]) == e {
			return i
		}
	}
	debugf("!Could not find element %v", e)
	return -1
}

func(this *Game) AddElement(e Element) {
  this.elements = append(this.elements, e)
}

type ElementAction func(Element)

func(this *Game) EachElement(action ElementAction) {
  for _,e := range this.elements {
    action(e)
  }
}

func(game *Game) collide() {
  for i,a:= range game.elements {
    for j,b := range game.elements {
      if j <= i {
        continue
      }

      if a.intersects(b) {
        if a.biggerThan(b) {
          a.absorb(b)
        } else {
          b.absorb(a)
        }
      }
    }
  }
}

func(game *Game) win() {
  if game.player.isDead() {
    game.gameState = lost
	return
  }

  for _, e := range game.elements {
    if e != game.player && !e.isDead() {
      game.gameState = running
		return
    }
  }
  game.gameState = won
}
