package main

import "fmt"
import "os"
import "github.com/go-gl/glfw"

func initGlfw(width, height int) {
  if err := glfw.Init(); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }

  glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
  glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
  if err := glfw.OpenWindow(width, height, 8, 8, 8, 8, 8, 8, glfw.Windowed); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }

  if glInit() != 0 {
   fmt.Println("gl error")
  }
  glfw.SetWindowSizeCallback(onResize)
  glfw.SetKeyCallback(onKey)
  glfw.SetSwapInterval(1)
  glLineWidth(2)
}

func terminateGlfw() {
  glfw.Terminate()
}

func onResize(w, h int) {
  glViewport(0, 0, w, h)
}

func onKey(key, state int) {
  if key == glfw.KeyEsc && state == glfw.KeyPress {
    glfw.CloseWindow()
  }
}

const KeyA int = 65
const KeyS int = 83
const KeyD int = 68
const KeyW int = 87
const KeySpace int = 32

