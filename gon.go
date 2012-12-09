package main

import (
	"log"
	"math"
	"reflect"
	"github.com/go-gl/glfw"
)

type Game struct {
	gameState GameState
	elements []Element
	player *Player
	program Program
	positionAttrib AttribLocation
	modelToCameraMatrixUniform, colorUniform UniformLocation 
	vao VertexArray
	text TextRenderer
}

type GameState byte

const width float64 = 800
const height float64 = 600
const initialized GameState = 0
const running GameState = 1
const won GameState = 2
const lost GameState = 3

func createElements() {
  things := make([]Element, 32)
  for i := range things {
    size := random(5, 9)
    location := Vector2{random(0,1) * width, random(0,1) * height}
    direction := Vector2{random(-1,1), random(-1,1)}
    things[i] = NewThing(location, direction, size)
  }
  game.player = &Player{Thing{location : Vector2{width / 2, height / 2}, targetSize : 10, size : 10}}
  game.elements = append(things, game.player)
}

var game Game

func main() {
	game = Game {}
  initGlfw(int(width),int(height))
	game.text = NewTextRenderer("./PixelCarnageMono.ttf", 18, 64)
  game.gameState = initialized
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

	vs, err := NewShader(GL_VERTEX_SHADER, `#version 150
    in vec4 position;
	uniform mat4 cameraToClipMatrix;
	uniform mat4 modelToCameraMatrix;	
    void main()
    {
		vec4 cameraPos = modelToCameraMatrix * position;
    	gl_Position = cameraToClipMatrix * cameraPos;
    }`)

	if err != nil {
		log.Printf("Error compiling vertex shader\n")
		log.Println(err)
	}

	fs, err := NewShader(GL_FRAGMENT_SHADER, `#version 150
		uniform vec4 color;
      out vec4 outputF;
      void main(void) {
        outputF = color;
      }`)

	if err != nil {
		log.Printf("Error compiling fragment shader\n")
		log.Println(err)
	}

	program := NewProgram(vs, fs)

	vao := glGenVertexArray()
	vao.Bind()

	verts := make([]Vector4,100)

	sides := len(verts)
	vert_scale := 1.0 / float64(sides)
	TWO_PI := math.Pi * 2.0
	for i := 0; i < sides; i++ {
		angle := float64(i) * TWO_PI * vert_scale
		verts[i] = Vector4{float32(math.Cos(angle)), float32(math.Sin(angle)),0,1}
	}

	vbo := glGenBuffer()
	vbo.Bind(GL_ARRAY_BUFFER)
	glBufferData(GL_ARRAY_BUFFER, int(reflect.TypeOf(Vector4{}).Size()) * len(verts), verts, GL_STATIC_DRAW)

	positionAttrib := program.GetAttribLocation("position")
	positionAttrib.AttribPointer(4, GL_FLOAT, false, 0, nil)
	positionAttrib.EnableArray()

	cameraToClipMatrixUniform := program.GetUniformLocation("cameraToClipMatrix")
	modelToCameraMatrixUniform := program.GetUniformLocation("modelToCameraMatrix")
	colorUniform := program.GetUniformLocation("color")

	zNear := 1.0
	zFar := 45.0
	frustumScale := float32(CalcFrustumScale(45.0))
	cameraToClipMatrix := NewMatrix4x4(0.0)
    cameraToClipMatrix[0].x = frustumScale
    cameraToClipMatrix[1].y = frustumScale
    cameraToClipMatrix[2].z = float32((zFar + zNear) / (zNear - zFar))
    cameraToClipMatrix[2].w = -1.0
    cameraToClipMatrix[3].z = float32((2 * zFar * zNear) / (zNear - zFar))

	program.Use()
	cameraToClipMatrixUniform.UniformMatrix4fv(cameraToClipMatrix)
	program.Unuse()
	vbo.Unbind(GL_ARRAY_BUFFER)
	vao.Unbind()
	game.program = program
	game.vao = vao
	game.modelToCameraMatrixUniform = modelToCameraMatrixUniform
	game.colorUniform = colorUniform
	return true
}

func CalcFrustumScale(fovDeg float64) float64 {
    const degToRad = 3.14159 * 2.0 / 360.0
    fovRad := fovDeg * degToRad
    return 1.0 / math.Tan(fovRad / 2.0)
}


func free_resources() {
  game.program.Delete()
}

func update(elapsed float64) {
  switch game.gameState {
    case running:
      run(elapsed)
    case initialized, won, lost:
      waitForReset()
  }
}

func run(elapsed float64) {
  for _,e := range game.elements {
    e.update(elapsed)
  }

  game.collide()
  game.win()
}

func waitForReset() {
  if keyDown(KeySpace) {
    createElements()
    game.gameState = running
  }
}

func render() {
	glClearColor(0.0, 0.0, 0.0, 0)
	glClear(GL_COLOR_BUFFER_BIT)
	switch game.gameState {
		case initialized:
			game.text.Draw("Hit space to play!", Vector4{0.75,-1,0,0}, Vector4{1,1,1,1})
		case running:
			game.program.Use()
			game.vao.Bind()	
			for _, e := range game.elements {
				if e.isDead() {
					continue
				}
				location := e.Location()
				translateMatrix := NewMatrix4x4(1.0)
				translateMatrix[3] = Vector4{float32((location.x/width) - 0.5), -float32((location.y/height) - 0.5), -1, 1}
				
				theScale := Vector4{float32((1.0 / width)*e.Size()), float32((1.0 / height)*e.Size()), 1, 1}
				scaleMatrix := NewMatrix4x4(1.0)
				scaleMatrix[0].x = theScale.x
				scaleMatrix[1].y = theScale.y
				scaleMatrix[2].z = theScale.z
				scaleMatrix[3] = Vector4{0.0,0.0,0.0, 1.0}
				modelToCameraMatrix := translateMatrix.mult(scaleMatrix)
				game.modelToCameraMatrixUniform.UniformMatrix4fv(modelToCameraMatrix)

				if e == game.player {
					game.colorUniform.Uniform4fv(Vector4{0,0,1,1})
				} else if e.biggerThan(game.player) {
					game.colorUniform.Uniform4fv(Vector4{1,0,0,1})
				} else {
					game.colorUniform.Uniform4fv(Vector4{0,1,0,1})
				}
				glDrawArrays(GL_LINE_LOOP, 0, 100)
			}
			game.vao.Unbind()
			game.program.Unuse()
		case won:
			game.text.Draw("You won! Hit space to play again.", Vector4{0.75,-1,0,0}, Vector4{1,1,1,1})
		case lost:
			game.text.Draw("You lost! Hit space to play again.", Vector4{0.5,-1,0,0}, Vector4{1,1,1,1})
	}
}

