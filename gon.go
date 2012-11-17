package main

import (
  "fmt"
  "github.com/go-gl/gl"
  "github.com/go-gl/glfw"
)
type GameState byte

var gameState GameState
var elements []Element
var player *Player
var font Font

const width float64 = 800
const height float64 = 600
var screenCenter Vector = Vector{width/2.0,height/2.0}
const initialized GameState = 0
const running GameState = 1
const won GameState = 2
const lost GameState = 3

func init() {
  var err error
  font, err = loadFont()
  if err != nil {
    fmt.Printf("error loading font: %v", err)
  }

  initGlfw(int(width),int(height))
  defer terminateGlfw()

  gameState = initialized
}

func createElements() {
  things := make([]Element, 32)
  for i := range things {
    size := random(5, 9)
    location := Vector{random(0,1) * width, random(0,1) * height}
    direction := Vector{random(-1,1), random(-1,1)}
    things[i] = NewThing(location, direction, size)
  }
  player = &Player{Thing{location : Vector{width / 2, height / 2}, targetSize : 10, size : 10}}
  elements = append(things, player)
}

func main() {
  previousFrameTime := glfw.Time()
  profiler := NewProfiler()

  for glfw.WindowParam(glfw.Opened) == 1 {
    profiler.start()

    now := glfw.Time()
    elapsed := now - previousFrameTime
    previousFrameTime = now
    update(elapsed)
    render()

    profiler.stop()
    profiler.render()
    glfw.SwapBuffers()
  }
}

func update(elapsed float64) {
  switch gameState {
    case running:
      run(elapsed)
    case initialized, won, lost:
      waitForReset()
  }
}

func run(elapsed float64) {
  for _,e := range elements {
    e.update(elapsed)
  }

  collide()
  gameState = win()
}

func waitForReset() {
  if keyDown(KeySpace) {
    createElements()
    gameState = running
  }
}

func render() {
  gl.ClearColor(0.0, 0.0, 0.0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  switch gameState {
    case running:
      gl.LoadIdentity()
      for _, e := range elements {
        if e.isDead() {
          continue
        }
        var color RGB
        if e == player {
          color = RGB{0,0,1}
        } else if e.biggerThan(player) {
          color = RGB{1,0,0}
        } else {
          color = RGB{0,1,0}
        }

        drawCircle(e.Location(), e.Size(), color)
      }
    case initialized:
      font.drawString(300, screenCenter.y,"Press Space to play!")
    case won:
      font.drawString(200, screenCenter.y,"WINNER! Press Space to play again!")
    case lost:
      font.drawString(100, screenCenter.y,"You were swallowed up :( Press Space to play again!")
  }
}

