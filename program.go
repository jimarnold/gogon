package main

import (
	"errors"
	"log"
)

func NewProgram(vs, fs Shader) Program {
	program := glCreateProgram()

	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	link_ok := program.Get(GL_LINK_STATUS)
	if link_ok == 0 {
		log.Printf("glLinkProgram:")
	}

	return program
}

func NewShader(shaderType GLenum, source string) (Shader,error) {
	s := glCreateShader(shaderType)
	s.Source(source)
	s.Compile()
	compile_ok := s.Get(GL_COMPILE_STATUS)
	if compile_ok == 0 {
		return Shader(0),errors.New(s.GetInfoLog())
	}
	return s, nil
}
