package main

// #cgo darwin LDFLAGS: -framework OpenGL
//#include <stdlib.h>
//#include <OpenGL/gl3.h>
//#include <OpenGL/gl3ext.h>
//#include <OpenGL/OpenGL.h>
import "C"
import "unsafe"
import "reflect"

type GLenum C.GLenum
type GLbitfield C.GLbitfield
type Object C.GLuint
type AttribLocation int
type GLclampf C.GLclampf

const (
	GL_LINE_SMOOTH = C.GL_LINE_SMOOTH
	GL_VERTEX_SHADER = C.GL_VERTEX_SHADER
	GL_FRAGMENT_SHADER = C.GL_FRAGMENT_SHADER
	GL_COMPILE_STATUS = C.GL_COMPILE_STATUS
	GL_LINK_STATUS = C.GL_LINK_STATUS
	GL_COLOR_BUFFER_BIT = C.GL_COLOR_BUFFER_BIT
	GL_DOUBLE = C.GL_DOUBLE
)

func glString(s string) *C.GLchar { return (*C.GLchar)(C.CString(s)) }

func freeString(ptr *C.GLchar) { C.free(unsafe.Pointer(ptr)) }

func glBool(v bool) C.GLboolean {
	if v {
		return 1
	}			
	return 0
}

func ptr(v interface{}) unsafe.Pointer {

	if v == nil {
		return unsafe.Pointer(nil)
	}

	rv := reflect.ValueOf(v)
	var et reflect.Value
	switch rv.Type().Kind() {
	case reflect.Uintptr:
		offset, _ := v.(uintptr)
		return unsafe.Pointer(offset)
	case reflect.Ptr:
		et = rv.Elem()
	case reflect.Slice:
		et = rv.Index(0)
	default:
		panic("type must be a pointer, a slice, uintptr or nil")
	}

	return unsafe.Pointer(et.UnsafeAddr())
}

func glEnable(cap GLenum) {
	C.glEnable(C.GLenum(cap))
}

func glDisable(cap GLenum) {
	C.glDisable(C.GLenum(cap))
}

func glLineWidth(width float32) {
    C.glLineWidth(C.GLfloat(width))
}

func glClear(mask GLbitfield) {
	C.glClear(C.GLbitfield(mask))
}

func glClearColor(red GLclampf, green GLclampf, blue GLclampf, alpha GLclampf) {
    C.glClearColor(C.GLclampf(red), C.GLclampf(green), C.GLclampf(blue), C.GLclampf(alpha))
}

type Program Object

func glCreateProgram() Program {
	return Program(C.glCreateProgram())
}

func(program Program) AttachShader(shader Shader) {
	C.glAttachShader(C.GLuint(program), C.GLuint(shader))
}

func (program Program) Link() { C.glLinkProgram(C.GLuint(program)) }

func (program Program) Get(param GLenum) int {
    var rv C.GLint

    C.glGetProgramiv(C.GLuint(program), C.GLenum(param), &rv)
    return int(rv)
}

func (program Program) GetAttribLocation(name string) AttribLocation {

    cname := glString(name)
    defer freeString(cname)

    return AttribLocation(C.glGetAttribLocation(C.GLuint(program), cname))
}

func (program Program) Use() { C.glUseProgram(C.GLuint(program)) }

func (program Program) Unuse() { C.glUseProgram(C.GLuint(0)) }

func (program Program) Delete() {
	C.glDeleteProgram(C.GLuint(program))
}

type Shader Object

func glCreateShader(shaderType GLenum) Shader {
	return Shader(C.glCreateShader(C.GLenum(shaderType)))
}

func(shader Shader) Source(source string) {
	csource := glString(source)
    defer freeString(csource)

    var one C.GLint = C.GLint(len(source))

    C.glShaderSource(C.GLuint(shader), 1, &csource, &one)
}

func (shader Shader) Compile() { C.glCompileShader(C.GLuint(shader)) }

func (shader Shader) Get(param GLenum) int {
    var rv C.GLint

    C.glGetShaderiv(C.GLuint(shader), C.GLenum(param), &rv)
    return int(rv)
}

func (indx AttribLocation) AttribPointer(size uint, typ GLenum, normalized bool, stride int, pointer interface{}) {
    C.glVertexAttribPointer(C.GLuint(indx), C.GLint(size), C.GLenum(typ),
        glBool(normalized), C.GLsizei(stride), ptr(pointer))
}

func (indx AttribLocation) EnableArray() {
    C.glEnableVertexAttribArray(C.GLuint(indx))
}

func (indx AttribLocation) DisableArray() {
    C.glDisableVertexAttribArray(C.GLuint(indx))
}
