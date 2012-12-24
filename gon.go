package main

import (
	"github.com/go-gl/glfw"
)

const width float64 = 800
const height float64 = 600


func main() {
	initGlfw()
	defer terminateGlfw()
	createWindow(int(width), int(height))
	game := NewGame()
	profiler := NewProfiler(game.text)
	defer game.delete()

	previousFrameTime := glfw.Time()
	profiler.start()
	for glfw.WindowParam(glfw.Opened) == 1 {
		now := glfw.Time()
		elapsed := now - previousFrameTime
		previousFrameTime = now
		game.update(elapsed)
		game.render()
		profiler.update()
		profiler.render()
		glfw.SwapBuffers()
	}
}

