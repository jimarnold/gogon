package main

import (
  "fmt"
  "math"
  //"reflect"
  //"github.com/go-gl/gl"
  "github.com/go-gl/glfw"
)
type GameState byte

var gameState GameState
var elements []Element
var player *Player
var font Font
var program Program
var attribute_coord2d AttribLocation

const width float64 = 800
const height float64 = 600
var screenCenter Vector = Vector{width/2.0,height/2.0}
const initialized GameState = 0
const running GameState = 1
const won GameState = 2
const lost GameState = 3
func init() {
  initGlfw(int(width),int(height))

  var err error
  font, err = loadFont()
  if err != nil {
    fmt.Printf("error loading font: %v", err)
  }

  gameState = initialized
}

func createElements() {
  things := make([]Element, 16)
  for i := range things {
    size := random(5, 9)
    location := Vector{random(0,1) * width, random(0,1) * height}
    direction := Vector{random(-1,1), random(-1,1)}
    things[i] = NewThing(location, direction, size)
  }
  player = &Player{Thing{location : Vector{width / 2, height / 2}, targetSize : 10, size : 10}}
  elements = append(things, player)
}

var verts []Vector
func main() {
  defer terminateGlfw()
  previousFrameTime := glfw.Time()
  profiler := NewProfiler()

  //vbo := gl.GenBuffer()
  //vbo.Bind(gl.ELEMENT_ARRAY_BUFFER)
  //verts := []Vector {Vector{100,100}, Vector{150,100}, Vector{100,200}}
  //vertSize := 3 * reflect.TypeOf(Vector{}).Size()
  //gl.BufferData(gl.ARRAY_BUFFER, int(vertSize), &verts, gl.STATIC_DRAW)
  //cbo := gl.GenBuffer()
  //cbo.Bind(gl.ARRAY_BUFFER)
  //colors := []RGB {RGB{0.0,0.0,1.0},RGB{1.0,0.0,0.0},RGB{0.0,1.0,0.0}}
  //colorSize := 3 * reflect.TypeOf(RGB{}).Size()
  //gl.BufferData(gl.ARRAY_BUFFER, int(colorSize), &colors, gl.STATIC_DRAW)
  init_resources()
  defer free_resources()
  verts = make([]Vector,100)

	const TWO_PI = 2.0 * math.Pi
  sides := len(verts)
  scale := 1.0 / float64(sides)
  for i := 0; i < sides; i++ {
    angle := float64(i) * TWO_PI * scale
    verts[i] = Vector{math.Cos(angle), math.Sin(angle)}
  }

  for glfw.WindowParam(glfw.Opened) == 1 {
    profiler.start()

    now := glfw.Time()
    elapsed := now - previousFrameTime
    previousFrameTime = now
    update(elapsed)
    render()

    //vbo.Bind(gl.ARRAY_BUFFER)
    //gl.EnableVertexAttribArray(0)
    //gl.VertexAttribPointer(0, 2, gl.FLOAT, 0, nil)
    ////gl.EnableClientState(gl.COLOR_ARRAY)
    ////cbo.Bind(gl.ARRAY_BUFFER)
    ////gl.ColorPointer(3, gl.FLOAT, 0, colors);
    //indices := []uint32{0, 1, 2}
    //gl.IndexPointer(3, gl.UNSIGNED_BYTE, indices)
    //gl.Color3f(1.0,1.0,1.0)
    //gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, indices)
    //gl.DisableClientState(gl.VERTEX_ARRAY)

    profiler.stop()
    profiler.render()
    glfw.SwapBuffers()
  }
}
func init_resources() bool {
  vs := glCreateShader(GL_VERTEX_SHADER)
  vs_source := `#version 120
    attribute vec2 coord2d;
    void main()
    {
      gl_Position = gl_ModelViewProjectionMatrix * vec4(coord2d, 0.0, 1.0);
    }`
  vs.Source(vs_source)
  vs.Compile()
  compile_ok := vs.Get(GL_COMPILE_STATUS)
  if  compile_ok == 0 {
    fmt.Printf("Error in vertex shader\n")
    return false
  }

  fs := glCreateShader(GL_FRAGMENT_SHADER)
  fs_source := `#version 120
      void main(void) {
        gl_FragColor[0] = 0.0;
        gl_FragColor[1] = 0.0;
        gl_FragColor[2] = 1.0;
      }`
  fs.Source(fs_source)
  fs.Compile()
  compile_ok = fs.Get(GL_COMPILE_STATUS)
  if compile_ok == 0 {
    fmt.Printf("Error in fragment shader\n")
    return false
  }

  program = glCreateProgram()
  program.AttachShader(vs)
  program.AttachShader(fs)
  program.Link()
  link_ok := program.Get(GL_LINK_STATUS)
  if link_ok == 0 {
    fmt.Printf("glLinkProgram:")
    return false
  }

  attribute_name := "coord2d"
  attribute_coord2d := program.GetAttribLocation(attribute_name)
  if attribute_coord2d == -1 {
    fmt.Printf("Could not bind attribute %s\n", attribute_name)
    return false
  }

  return true
}

func free_resources() {
  program.Delete()
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
  glClearColor(0.0, 0.0, 0.0, 0)
  glClear(GL_COLOR_BUFFER_BIT)

  switch gameState {
    case running:
      program.Use()

      attribute_coord2d.EnableArray()
      attribute_coord2d.AttribPointer(
        2,                 // number of elements per vertex, here (x,y)
        GL_DOUBLE,          // the type of each element
        false,             // take our values as-is
        0,                 // no extra data between each position
        verts)

      for _, e := range elements {
        if e.isDead() {
          continue
        }
        //location := e.Location()
        //glPushMatrix()
        //glTranslated(location.x, location.y, 0.)
        //glScaled(e.Size(), e.Size(), 0.0)
        //glDrawArrays(gl.LINE_LOOP, 0, len(verts))
        //glPopMatrix()
      }

      attribute_coord2d.DisableArray()
      program.Unuse()
    case initialized:
      font.drawString(300, screenCenter.y,"Press Space to play!")
    case won:
      font.drawString(200, screenCenter.y,"WINNER! Press Space to play again!")
    case lost:
      font.drawString(100, screenCenter.y,"You were swallowed up :( Press Space to play again!")
  }
}

