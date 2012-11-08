package main

import "fmt"
import "github.com/go-gl/gl"
import "github.com/go-gl/glfw"

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
  gl.Enable(gl.BLEND)
  gl.Enable(gl.LINE_SMOOTH)
  const GL_MULTISAMPLE_ARB = 0x809D
  gl.Enable(GL_MULTISAMPLE_ARB)
  gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST);
  gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
  gl.LineWidth(1)
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

