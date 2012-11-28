package main

import (
	"log"
	"math"
	"reflect"
	"github.com/go-gl/glfw"
)

type GameState byte

const width float64 = 800
const height float64 = 600
const initialized GameState = 0
const running GameState = 1
const won GameState = 2
const lost GameState = 3

var gameState GameState
var elements []Element
var player *Player
var program Program
var positionAttrib AttribLocation
var modelToCameraMatrixUniform UniformLocation 
var vao VertexArray
var screenCenter Vector = Vector{width/2.0,height/2.0}

func init() {
  initGlfw(int(width),int(height))
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
  defer terminateGlfw()
  previousFrameTime := glfw.Time()

  init_resources()
  defer free_resources()

  for glfw.WindowParam(glfw.Opened) == 1 {
    now := glfw.Time()
    elapsed := now - previousFrameTime
    previousFrameTime = now
    update(elapsed)
    render()
    glfw.SwapBuffers()
  }
}

func init_resources() bool {
	vao = glGenVertexArray()
	vao.Bind()

	verts := make([]vec4,100)

	sides := len(verts)
	vert_scale := 1.0 / float64(sides)
	TWO_PI := math.Pi * 2.0
	for i := 0; i < sides; i++ {
		angle := float64(i) * TWO_PI * vert_scale
		verts[i] = vec4{float32(math.Cos(angle)), float32(math.Sin(angle)),0,1}
	}

	vbo := glGenBuffer()
	vbo.Bind(GL_ARRAY_BUFFER)
	glBufferData(GL_ARRAY_BUFFER, int(reflect.TypeOf(vec4{}).Size()) * len(verts), verts, GL_STATIC_DRAW)

	vs_source := `#version 150
    in vec4 position;
	uniform mat4 cameraToClipMatrix;
	uniform mat4 modelToCameraMatrix;	
    void main()
    {
		vec4 cameraPos = modelToCameraMatrix * position;
    	gl_Position = cameraToClipMatrix * cameraPos;
    }`
	vs := glCreateShader(GL_VERTEX_SHADER)
	vs.Source(vs_source)
	vs.Compile()
	compile_ok := vs.Get(GL_COMPILE_STATUS)
	if  compile_ok == 0 {
		log.Printf("Error in vertex shader\n")
		log.Println(vs.GetInfoLog())
		return false
	}

	fs := glCreateShader(GL_FRAGMENT_SHADER)
	fs_source := `#version 150
      out vec4 outputF;
      void main(void) {
        outputF = vec4(1.0,0.0,1.0,0.5);
      }`
	fs.Source(fs_source)
	fs.Compile()
	compile_ok = fs.Get(GL_COMPILE_STATUS)
	if compile_ok == 0 {
		log.Printf("Error in fragment shader\n")
		return false
	}

	program = glCreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	link_ok := program.Get(GL_LINK_STATUS)
	if link_ok == 0 {
		log.Printf("glLinkProgram:")
		return false
	}

	positionAttrib := program.GetAttribLocation("position")
	positionAttrib.AttribPointer(4, GL_FLOAT, false, 0, nil)
	positionAttrib.EnableArray()

	cameraToClipMatrixUniform := program.GetUniformLocation("cameraToClipMatrix")
	modelToCameraMatrixUniform = program.GetUniformLocation("modelToCameraMatrix")

	zNear := 1.0
	zFar := 45.0
	frustumScale := float32(CalcFrustumScale(45.0))
	cameraToClipMatrix := Newmat4(0.0)
    cameraToClipMatrix[0].x = frustumScale
    cameraToClipMatrix[1].y = frustumScale
    cameraToClipMatrix[2].z = float32((zFar + zNear) / (zNear - zFar))
    cameraToClipMatrix[2].w = -1.0
    cameraToClipMatrix[3].z = float32((2 * zFar * zNear) / (zNear - zFar))

	program.Use()
	cameraToClipMatrixUniform.UniformMatrix4fv(cameraToClipMatrix)
	program.Unuse()
	vao.Unbind()
	return true
}

func CalcFrustumScale(fovDeg float64) float64 {
    const degToRad = 3.14159 * 2.0 / 360.0
    fovRad := fovDeg * degToRad
    return 1.0 / math.Tan(fovRad / 2.0)
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
		case running, initialized:
			program.Use()
			vao.Bind()	
			for _, e := range elements {
				if e.isDead() {
					continue
				}
				location := e.Location()
				translateMatrix := Newmat4(1.0)
				translateMatrix[3] = vec4{float32((location.x/width) - 0.5), -float32((location.y/height) - 0.5), -1, 1}
				
				theScale := vec4{float32((1.0 / width)*e.Size()), float32((1.0 / height)*e.Size()), 1, 1}
				scaleMatrix := Newmat4(1.0)
				scaleMatrix[0].x = theScale.x
				scaleMatrix[1].y = theScale.y
				scaleMatrix[2].z = theScale.z
				scaleMatrix[3] = vec4{0.0,0.0,0.0, 1.0}
				modelToCameraMatrix := translateMatrix.mult(scaleMatrix)
				modelToCameraMatrixUniform.UniformMatrix4fv(modelToCameraMatrix)
				glDrawArrays(GL_LINE_LOOP, 0, 100)
			}
			vao.Unbind()
			program.Unuse()
		case won:
		case lost:
	}
}

