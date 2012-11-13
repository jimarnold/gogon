package main

import (
  "fmt"
  "os"
  "github.com/go-gl/gl"
  "github.com/go-gl/glfw"
  "github.com/go-gl/gltext"
)
type GameState string

var gameState GameState
var elements []Element
var player *Player

const width float64 = 800
const height float64 = 600
var screenCenter Vector = Vector{width/2.0,height/2.0}
const running GameState = "running"
const won GameState = "won"
const lost GameState = "lost"
const initialized GameState = "initialized"

func init() {
  gameState = initialized
}

func createElements() {
  things := make([]Element, 32)
  for i := range things {
    aSize := random(5, 9)
    things[i] = &Thing {
      location : Vector{random(0,1) * width, random(0,1) * height},
      direction : Vector{random(-1,1), random(-1,1)},
      targetSize : aSize,
      size : aSize,
    }
  }
  player = &Player{Thing{location : Vector{width / 2, height / 2}, targetSize : 8, size : 8}}
  elements = append(things, player)
}

func main() {
  initGlfw(int(width),int(height))
  defer terminateGlfw()

  initText()
  previousFrameTime := glfw.Time()
  for glfw.WindowParam(glfw.Opened) == 1 {
    now := glfw.Time()
    elapsed := now - previousFrameTime
    previousFrameTime = now
    update(elapsed)
    render()
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

var font *gltext.Font
func initText() {
  file := "./ComingSoon.ttf"
  var err error
  font, err = loadFont(file, 24)
  if err != nil {
    fmt.Printf("error loading font: %v", err)
  //defer font.Release()
  }
}

// loadFont loads the specified font at the given scale.
func loadFont(file string, scale int32) (*gltext.Font, error) {
  fd, err := os.Open(file)
  if err != nil {
    return nil, err 
  }

  defer fd.Close()

  return gltext.LoadTruetype(fd, scale, 32, 127, gltext.LeftToRight)
}

func drawString(x, y float64, str string) error {
  _, h := font.GlyphBounds()
  y += float64(h)

  gl.Color4f(1, 1, 1, 1)
  err := font.Printf(float32(x), float32(y), str)
  return err
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
    case initialized, won, lost:
      //gl.LoadIdentity()
      drawString(screenCenter.x, screenCenter.y,"Press Space to play!")
  }
}

