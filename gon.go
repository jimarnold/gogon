package main

import "math"
import "time"
import "github.com/go-gl/gl"
import "github.com/go-gl/glfw"

var elements []Element
var player Element
const width float64 = 800
const height float64 = 600

func main() {
  things := make([]Element, 64)
  for i := range things {
    aSize := random(1, 5)
    things[i] = &Thing {
      location : Vector{random(0,1) * width, random(0,1) * height},
      direction : Vector{random(-1,1), random(-1,1)},
      targetSize : aSize,
      size : aSize,
    }
  }
  player = &Player{Thing{location : Vector{width / 2, height / 2}, targetSize : 10, size : 10}}
  elements = append(things, player)

  initGlfw(int(width),int(height))
  defer terminateGlfw()

  previousFrameTime := time.Now()
  for glfw.WindowParam(glfw.Opened) == 1 {
    now := time.Now()
    elapsed := now.Sub(previousFrameTime).Seconds()
    previousFrameTime = now
    update(elapsed)
    render()
    glfw.SwapBuffers()
  }
}

func update(elapsed float64) {
  for _,e := range elements {
    e.update(elapsed)
  }

  collide()
  win()
}

func render() {
  gl.ClearColor(0.0, 0.0, 0.0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  gl.LoadIdentity()
  for _, e := range elements {
    if e.isDead() {
      continue
    }
    x:=e.getLocation().x
    y:=e.getLocation().y
    radius := e.getSize()
    if e == player {
      gl.Color3ub(0,0,255)
    } else if e.biggerThan(player) {
      gl.Color3ub(255,0,0)
    } else {
      gl.Color3ub(0,255,0)
    }

    gl.Begin(gl.LINE_LOOP)
      for i := 0.0; i < radius * 8; i++ {
        angle := i*2.0*math.Pi/(radius * 8)
        gl.Vertex2d(x + (math.Cos(angle) * radius), y + (math.Sin(angle) * radius))
      }
    gl.End()
  }
}
