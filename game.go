package main

import (
	"log"
	"math"
	"reflect"
	"github.com/go-gl/gl"
        glfw "github.com/go-gl/glfw3"
	"github.com/jimarnold/gltext"
)

type Game struct {
	gameState GameState
	elements *Elements
	player *Player
	program gl.Program
	positionAttrib gl.AttribLocation
	cameraToClipMatrixUniform, colorUniform gl.UniformLocation 
	cameraToClipMatrix Matrix4x4
	vao gl.VertexArray
	text *gltext.Font
	totalTime float64
	elapsedSpawnTime float64
        window *glfw.Window
}

type GameState byte

const initialized GameState = 0
const running GameState = 1
const won GameState = 2
const lost GameState = 3

func NewGame(window *glfw.Window) *Game {
	game := &Game{window:window}
	game.init()
	return game
}

func(this *Game) init() {

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

	this.program = program
	this.vao = vao
	this.cameraToClipMatrixUniform = cameraToClipMatrixUniform
	this.cameraToClipMatrix = cameraToClipMatrix
	this.colorUniform = colorUniform
	this.text = gltext.NewFont("./PixelCarnageMono.ttf", 18, 64, float32(width), float32(height))
	this.gameState = initialized
}

func(this *Game) delete() {
	this.program.Delete()
}
func(this *Game) start() {
	if this.gameState == running {
		return
	}
	this.createElements()
    this.gameState = running
}

func(this *Game) createElements() {
	this.player = NewPlayer(Vector2{width / 2, height / 2}, this.window)
	this.elements = &Elements{make([]Element, 0)}
	this.elements.Add(this.player)
	this.elements.Add(createEnemy())
}

func(this *Game) update(elapsed float64) {
	this.totalTime += elapsed
	switch this.gameState {
		case running:
			this.run(elapsed)
		case initialized, won, lost:
			this.waitForReset()
	}
}

func(this *Game) waitForReset() {
	if keyDown(this.window, glfw.KeySpace) {
		this.start()
	}
}

const enemySpawnInterval float64 = 0.5

func(this *Game) run(elapsed float64) {
	this.elapsedSpawnTime += elapsed
	if this.elapsedSpawnTime > enemySpawnInterval {
		this.elapsedSpawnTime = 0
		this.elements.Add(createEnemy())
	}

	this.elements.Each(func(_ int,e Element) {
		e.update(elapsed)
	})

	this.collide()
	this.prune()
	this.win()
}

func(this *Game) prune() {
	deadThings := make([]Element, 0)
	this.elements.Each(func(_ int,e Element) {
		if e.isDead() {
			deadThings = append(deadThings, e)
		}
	})
	for _,e := range deadThings {
		this.elements.Delete(e)
	}
}

func createEnemy() Element {
	size := random(8, 32)
	location := Vector2{width, random(0,1) * height}
	direction := Vector2{random(-3,-0.5), random(-0.1,0.1)}
	r := random(0,1)
	if r > 0.0 && r < 0.1 {
		return NewPickup(location, direction)
	}
	if r > 0.92 {
		return NewShrinker(location, direction)
	}
	return NewThing(location, direction, size)
}

func(game *Game) collide() {
	game.elements.Each(func(i int, e Element) {
		if e == game.player {
			return
		}

		if game.player.intersects(e) {
			game.player.collideWith(e)			
		}
	})
}

func(game *Game) win() {
	if game.player.isDead() {
		game.gameState = lost
		return
	}

	game.gameState = running
}

type Color4f []float32

func(this *Game) render() {
	gl.ClearColor(0,0,0,0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	switch this.gameState {
		case initialized:
			this.text.Printf(0.75, -1, "Hit space to play!")
		case running:
			this.program.Use()
			this.vao.Bind()	
			this.elements.Each(func(_ int, e Element) {
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
				clipMatrix := this.cameraToClipMatrix.mult(modelToCameraMatrix)
                                var clipMatrixArray = clipMatrix.toa();
				this.cameraToClipMatrixUniform.UniformMatrix4f(false, &clipMatrixArray)
				this.colorUniform.Uniform4fv(1, e.Color())
				gl.DrawArrays(gl.LINE_LOOP, 0, 100)
			})
			this.vao.Unbind()
			this.program.Unuse()
			this.text.Printf(0.05, 0, "Score: %d", this.player.score)
			this.text.Printf(0.35, 0, "Distance travelled: %f", this.totalTime * 10)
		case won:
			this.text.Printf(0.75, -1, "You won! Hit space to play again.")
		case lost:
			this.text.Printf(0.5, -1, "You lost! Hit space to play again.")
	}
}
