package main

import (
	"errors"
	"log"
  "github.com/jimarnold/gl"
)

func NewProgram(vs, fs gl.Shader) gl.Program {
	program := gl.CreateProgram()

	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	link_ok := program.Get(gl.LINK_STATUS)
	if link_ok == 0 {
		log.Printf("glLinkProgram:")
	}

	return program
}

func NewShader(shaderType gl.GLenum, source string) (gl.Shader,error) {
	s := gl.CreateShader(shaderType)
	s.Source(source)
	s.Compile()
	compile_ok := s.Get(gl.COMPILE_STATUS)
	if compile_ok == 0 {
		return gl.Shader(0),errors.New(s.GetInfoLog())
	}
	return s, nil
}
