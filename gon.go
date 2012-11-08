package main

import "math"
import "math/rand"
import "time"
import "fmt"
import "github.com/go-gl/gl"
import "github.com/go-gl/glfw"

var things []Thing
const width float64 = 800
const height float64 = 600

func main() {
  things = make([]Thing, 16)
  rand.Seed(time.Now().UTC().UnixNano())
  for i:=range things {
    things[i].location = Vector{random(0,1) * width, random(0,1) * height}
    things[i].direction = Vector{random(-1,1), random(-1,1)}
    things[i].size = float64(i)
  }
  initGlfw(int(width),int(height))

  for glfw.WindowParam(glfw.Opened) == 1 {
    elapsed := 0.01666666
    update(elapsed)
    render()
    glfw.SwapBuffers()
  }
}

func update(elapsed float64) {
  for i:= range things {
    thing := things[i]
    newX := thing.location.x + (elapsed * thing.direction.x * 100)
    newY := thing.location.y + (elapsed * thing.direction.y * 100)
    thing.location.x = newX
    thing.location.y = newY
    if thing.location.x > width {
      thing.location.x = thing.location.x - width
    }
    if thing.location.x < 0 {
      thing.location.x = thing.location.x + width
    }
    if thing.location.y > height {
      thing.location.y = thing.location.y - height
    }
    if thing.location.y < 0 {
      thing.location.y = thing.location.y + height
    }

    things[i] = thing
  }

  for i:= range things {
    for j := range things {
      if i == j || (things[i].size == 0 || things[j].size == 0) {
        continue
      }

      if things[i].intersects(things[j]) {
        if things[i].size > things[j].size {
          things[i].absorb(&things[j])
        } else {
          things[j].absorb(&things[i])
        }
      }
    }
  }
}

func(this Thing) intersects(other Thing) bool {
  distance := math.Sqrt(math.Pow((other.location.y - this.location.y), 2) + math.Pow((other.location.x - this.location.x),2) )
  return (this.size + other.size) >= distance
}

func(this *Thing) absorb(other *Thing) {
  this.size += other.size
  other.size = 0
}

func random(min, max float64) float64 {
  return rand.Float64() * (max - min) + min
}

type Vector struct {
  x,y float64
}

type Thing struct {
  location Vector
  direction Vector
  size float64
}

func initGlfw(width, height int) {
  var err error
  if err = glfw.Init(); err != nil {
    fmt.Printf("%v\n", err)
    return
  }

  glfw.OpenWindowHint(glfw.FsaaSamples, 8)

  if err = glfw.OpenWindow(width, height, 8, 8, 8, 8, 8, 8, glfw.Windowed); err != nil {
    fmt.Printf("%logv\n", err)
    return
  }
  glfw.SetWindowSizeCallback(onResize)
  glfw.SetKeyCallback(onKey)
  glfw.SetSwapInterval(1)
  gl.Disable(gl.LIGHTING)
}

func terminateGlfw() {
  glfw.Terminate()
}

func onResize(w, h int) {
  gl.MatrixMode(gl.PROJECTION)
  gl.LoadIdentity()
  gl.Viewport(0, 0, w, h)
  gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
  gl.MatrixMode(gl.MODELVIEW)
  gl.LoadIdentity()
}

func onKey(key, state int) {
  if key == glfw.KeyEsc && state == glfw.KeyPress {
    glfw.CloseWindow()
  }
}

func render() {
  const GL_MULTISAMPLE_ARB = 0x809D
  gl.ClearColor(0.0, 0.0, 0.0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  gl.Enable(gl.BLEND)
  gl.Enable(gl.LINE_SMOOTH)
  gl.Enable(GL_MULTISAMPLE_ARB)
  gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST);
  gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
  gl.LineWidth(1)
  gl.LoadIdentity()
  gl.Color3ub(255,255,255)
  for t := range things {
    if things[t].size == 0 {
      continue
    }
    gl.Begin(gl.LINE_LOOP)
      x:=things[t].location.x
      y:=things[t].location.y
      radius := things[t].size
      for i := 0.0; i < radius * 8; i++ {
        angle := i*2.0*math.Pi/(radius * 8)
        gl.Vertex2d(x + (math.Cos(angle) * radius), y + (math.Sin(angle) * radius))
      }
    gl.End()
  }
}
