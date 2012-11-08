package main

import "math"
import "math/rand"
import "time"
import "fmt"
import "github.com/go-gl/gl"
//import "github.com/go-gl/glu"
import "github.com/go-gl/glfw"

var things []Thing

func main() {
  things = make([]Thing, 16)
  rand.Seed(time.Now().UTC().UnixNano())
  for i:=range things {
    things[i].location.x = rand.Float64() * 800
    things[i].location.y = rand.Float64() * 600
    things[i].size = float64(i)
  }
  initGlfw(800,600)

  for glfw.WindowParam(glfw.Opened) == 1 {
    elapsed := 0.01666666
    update(elapsed)
    //draw(elapsed)
    render()
    glfw.SwapBuffers()
  }
}

func update(elapsed float64) {
  for i:= range things {
    newX := math.Mod(things[i].location.x + (elapsed * float64(100 - things[i].size) * 4.0), 800)
    newY := math.Mod(things[i].location.y + (elapsed * float64(100 - things[i].size) * 4.0), 600)
    things[i].location.x = newX
    things[i].location.y = newY
  }

  for i:= range things {
    for j := range things {
      if i == j {
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

func draw(elapsed float64) {
  for i:= range things {
    fmt.Printf("%d: (%f) %f, %f\n", i, things[i].size, things[i].location.x, things[i].location.y)
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

type Vector struct {
  x,y float64
}

type Thing struct {
  location Vector
  size float64
}

func initGlfw(width, height int) {
  var err error
  if err = glfw.Init(); err != nil {
    fmt.Printf("%v\n", err)
    return
  }

  if err = glfw.OpenWindow(width, height, 0, 0, 0, 0, 0, 0, glfw.Windowed); err != nil {
    fmt.Printf("%logv\n", err)
    return
  }

  glfw.SetWindowSizeCallback(onResize)
  glfw.SetKeyCallback(onKey)
  glfw.SetSwapInterval(1)
  gl.Enable(gl.DEPTH_TEST)
  gl.Disable(gl.LIGHTING)
  gl.ClearDepth(1)
  gl.DepthFunc(gl.LEQUAL)
  gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
}

func terminateGlfw() {
  glfw.Terminate()
}

func onResize(w, h int) {
  gl.MatrixMode(gl.PROJECTION)
  gl.LoadIdentity()
  gl.Viewport(0, 0, w, h)
  //glu.Perspective(65.0, float64(w)/float64(h), 0.1, 2000.0)
gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
  gl.MatrixMode(gl.MODELVIEW)
  gl.LoadIdentity()
}

func onKey(key, state int) {
  if glfw.Key(key) == glfw.KeyPress {
    glfw.CloseWindow()
  }
}

func render() {
  gl.ClearColor(0.0, 0.0, 0.05, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  gl.LoadIdentity()
  gl.PointSize(1)
  gl.Begin(gl.POINTS)
    for t := range things {
    x:=things[t].location.x
    y:=things[t].location.y
    radius := things[t].size
    for i := 0; i < 100; i++ {
        angle := float64(i)*2.0*math.Pi/100.0
        gl.Color3ub(255,255,255)
        gl.Vertex2f(float32(x + (math.Cos(angle) * radius)), float32(y + (math.Sin(angle) * radius)))
    }
    }
      //gl.Color3ub(uint8(k.x),uint8(k.y),uint8(k.z))//(rgb.r, rgb.g, rgb.b)
      //gl.Vertex3f(float32(k.x), float32(k.y), float32(k.z))
  gl.End()
}
