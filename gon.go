package main

import "os"
import "math"
import "time"
import "github.com/go-gl/gl"
import "github.com/go-gl/glfw"

var things []Element
var player Element
const width float64 = 800
const height float64 = 600

func main() {
  things = make([]Element, 64)
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
  things = append(things, player)
  err := initGlfw(int(width),int(height))
  if err != nil {
    os.Exit(1)
  }

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
  for _,thing := range things {
    thing.update(elapsed)
  }

  collide()
  win()
}

type Vector struct {
  x,y float64
}

func(v1 Vector) Add(v2 Vector) Vector {
  return Vector{v1.x + v2.x, v2.y + v2.y}
}

func(v Vector) Mult(s float64) Vector {
  return Vector{v.x * s, v.y * s}
}

func render() {
  gl.ClearColor(0.0, 0.0, 0.0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  gl.LoadIdentity()
  for _, thing := range things {
    if thing.isDead() {
      continue
    }
    x:=thing.getLocation().x
    y:=thing.getLocation().y
    radius := thing.getSize()
    if thing == player {
      gl.Color3ub(0,0,255)
    } else if thing.biggerThan(player) {
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
