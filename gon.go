package main

import "math"
import "time"
import "github.com/go-gl/gl"
import "github.com/go-gl/glfw"

var things []*Thing
const width float64 = 800
const height float64 = 600

func main() {
  things = make([]*Thing, 64)
  for i := range things {
    aSize := random(1, 5)
    things[i] = &Thing {
      location : Vector{random(0,1) * width, random(0,1) * height},
      direction : Vector{random(-1,1), random(-1,1)},
      targetSize : aSize,
      size : aSize,
    }
  }

  initGlfw(int(width),int(height))

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
}

type Vector struct {
  x,y float64
}

func render() {
  gl.ClearColor(0.0, 0.0, 0.0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  gl.LoadIdentity()
  for _, thing := range things {
    if thing.size == 0 {
      continue
    }
    x:=thing.location.x
    y:=thing.location.y
    radius := thing.size
    gl.Color3ub(255 - uint8(thing.size),255,255 - uint8(thing.size))
    gl.Begin(gl.LINE_LOOP)
      for i := 0.0; i < radius * 8; i++ {
        angle := i*2.0*math.Pi/(radius * 8)
        gl.Vertex2d(x + (math.Cos(angle) * radius), y + (math.Sin(angle) * radius))
      }
    gl.End()
  }
}
