package main

import "fmt"
import "os"
import glfw "github.com/go-gl/glfw3"
import "github.com/go-gl/gl"

func initGlfw() {
	if !glfw.Init() {
                panic("Can't init glfw");
		os.Exit(1)
	}
}

func errorCallback(err glfw.ErrorCode, desc string) {
    fmt.Printf("%v: %v\n", err, desc)
}


func createWindow(width, height int) *glfw.Window {

    glfw.SetErrorCallback(errorCallback)
    glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True);
    glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 2)
	window, err := glfw.CreateWindow(width, height, "Gon", nil, nil)
        if err != nil {
            fmt.Println(err)
            panic("Unable to create window")
        }
        window.MakeContextCurrent()
	if gl.Init() != 0 {
		fmt.Println("error initializing OpenGL")
	}
	window.SetSizeCallback(onResize)
	window.SetKeyCallback(onKey)
	glfw.SwapInterval(1)
	gl.LineWidth(2)
        return window
}

func terminateGlfw() {
	glfw.Terminate()
}

func onResize(window *glfw.Window, w, h int) {
	gl.Viewport(0, 0, w, h)
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

const KeyA int = 65
const KeyS int = 83
const KeyD int = 68
const KeyW int = 87
const KeySpace int = 32

