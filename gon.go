package main

import (
	"log"
	"math"
	"reflect"
	"github.com/go-gl/glfw"
	"github.com/jimarnold/gl"
	"github.com/jimarnold/gltext"
)
type GameState byte

const width float64 = 800
const height float64 = 600
const initialized GameState = 0
const running GameState = 1
const won GameState = 2
const lost GameState = 3

var game Game

func main() {
	game = Game {}
	initGlfw()
	createWindow(int(width), int(height))
	game.text = gltext.NewFont("./PixelCarnageMono.ttf", 18, 64, float32(width), float32(height))
	game.gameState = initialized
	defer terminateGlfw()
	previousFrameTime := glfw.Time()
	profiler := NewProfiler(game.text)
	init_resources()
	defer free_resources()

	profiler.start()
	for glfw.WindowParam(glfw.Opened) == 1 {
		now := glfw.Time()
		elapsed := now - previousFrameTime
		previousFrameTime = now
		game.update(elapsed)
		render()
		profiler.update()
		profiler.render()
		glfw.SwapBuffers()
	}
}

func init_resources() bool {

	vs, err := NewShader(gl.VERTEX_SHADER, `#version 150
    in vec4 position;
    uniform mat4 cameraToClipMatrix;
    void main()
    {
    	gl_Position = cameraToClipMatrix * position;
    }`)

	if err != nil {
		log.Printf("Error compiling vertex shader\n")
		log.Println(err)
	}

	fs, err := NewShader(gl.FRAGMENT_SHADER, `#version 150
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

	verts := make([]Vector4,100)

	sides := len(verts)
	vert_scale := 1.0 / float64(sides)
	TWO_PI := math.Pi * 2.0
	for i := 0; i < sides; i++ {
		angle := float64(i) * TWO_PI * vert_scale
		verts[i] = Vector4{float32(math.Cos(angle)), float32(math.Sin(angle)),0,1}
	}

	vao := gl.GenVertexArray()
	vao.Bind()

	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(reflect.TypeOf(Vector4{}).Size()) * len(verts), verts, gl.STATIC_DRAW)

	positionAttrib := program.GetAttribLocation("position")
	positionAttrib.AttribPointer(4, gl.FLOAT, false, 0, nil)
	positionAttrib.EnableArray()

	vbo.Unbind(gl.ARRAY_BUFFER)
	vao.Unbind()

	cameraToClipMatrixUniform := program.GetUniformLocation("cameraToClipMatrix")
	colorUniform := program.GetUniformLocation("color")

	zNear := float32(0.0)
	zFar := float32(45.0)
	cameraToClipMatrix := ortho(0, float32(width), 0, float32(height), zNear, zFar)

	game.program = program
	game.vao = vao
	game.cameraToClipMatrixUniform = cameraToClipMatrixUniform
	game.cameraToClipMatrix = cameraToClipMatrix
	game.colorUniform = colorUniform
	return true
}

func ortho(left, right, bottom, top, zNear, zFar float32) Matrix4x4 {
	result := NewMatrix4x4(1.0)
	result[0].x = float32(2.0) / (right - left)
	result[1].y = float32(2.0) / (top - bottom)
	result[2].z = float32(-2.0) / (zFar - zNear)
	result[3].x = -(right + left) / (right - left)
	result[3].y = -(top + bottom) / (top - bottom)
	result[3].z = -(zFar + zNear) / (zFar - zNear)
	return result
}

func free_resources() {
	game.program.Delete()
}

func waitForReset() {
	if keyDown(KeySpace) {
		game.start()
	}
}

func render() {
	gl.ClearColor(0.0, 0.0, 0.0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	switch game.gameState {
		case initialized:
			game.text.Printf(0.75, -1, "Hit space to play!")
		case running:
			game.program.Use()
			game.vao.Bind()	
			game.elements.Each(func(_ int, e Element) {
				if e.isDead() {
					return
				}
				location := e.Location()
				scale := float32(e.Size())

				translateMatrix := NewMatrix4x4(1.0)
				translateMatrix[3] = Vector4{float32(location.x), float32(location.y), -1, 1}
				scaleMatrix := NewMatrix4x4(1.0)
				scaleMatrix[0].x = scale
				scaleMatrix[1].y = scale
				modelToCameraMatrix := translateMatrix.mult(scaleMatrix)
				clipMatrix := game.cameraToClipMatrix.mult(modelToCameraMatrix)
				game.cameraToClipMatrixUniform.UniformMatrix4f(&(clipMatrix[0].x))

				if e == game.player {
					game.colorUniform.Uniform4fv(Vector4{0,0,1,1}.To_a())
				} else if e.biggerThan(game.player) {
					game.colorUniform.Uniform4fv(Vector4{1,0,0,1}.To_a())
				} else {
					game.colorUniform.Uniform4fv(Vector4{0,1,0,1}.To_a())
				}
				gl.DrawArrays(gl.LINE_LOOP, 0, 100)
			})
			game.vao.Unbind()
			game.program.Unuse()
		case won:
			game.text.Printf(0.75, -1, "You won! Hit space to play again.")
		case lost:
			game.text.Printf(0.5, -1, "You lost! Hit space to play again.")
	}
}

