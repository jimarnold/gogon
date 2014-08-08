package main

import (
	glfw "github.com/go-gl/glfw3"
)

const width float64 = 800
const height float64 = 600

func main() {
	initGlfw()
	defer terminateGlfw()
	window := createWindow(int(width), int(height))
	game := NewGame(window)
	profiler := NewProfiler(game.text)
	defer game.delete()

	previousFrameTime := glfw.GetTime()
	profiler.start()
	for !window.ShouldClose() {
		now := glfw.GetTime()
		elapsed := now - previousFrameTime
		previousFrameTime = now
		game.update(elapsed)
		game.render()
		profiler.update()
		profiler.render()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
