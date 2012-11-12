package main

import "math"
import "github.com/go-gl/gl"
import "github.com/go-gl/glfw"

var elements []Element
var player *Player
const width float64 = 800
const height float64 = 600

func main() {
  things := make([]Element, 16)
  for i := range things {
    aSize := random(4, 8)
    things[i] = &Thing {
      location : Vector{random(0,1) * width, random(0,1) * height},
      direction : Vector{random(-1,1), random(-1,1)},
      targetSize : aSize,
      size : aSize,
    }
  }
  player = &Player{Thing{location : Vector{width / 2, height / 2}, targetSize : 8, size : 8}}
  elements = append(things, player)

  initGlfw(int(width),int(height))
  defer terminateGlfw()

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
}

func drawCircle(location Vector, radius float64, color RGB) {
  gl.Color3d(color.r, color.g, color.b)
  gl.Begin(gl.LINE_LOOP)
    for i := 0.0; i < radius * 8; i++ {
      angle := i*2.0*math.Pi/(radius * 8)
      gl.Vertex2d(location.x + (math.Cos(angle) * radius), location.y + (math.Sin(angle) * radius))
    }
  gl.End()
}
