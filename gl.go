package main

// #cgo darwin LDFLAGS: -framework OpenGL -lGLEW
// #cgo windows LDFLAGS: -lglew32 -lopengl32
// #cgo linux LDFLAGS: -lGLEW -lGL
// #include "gl.h"
import "C"
import "unsafe"
import "reflect"

type GLenum C.GLenum
type GLbitfield C.GLbitfield
type Object C.GLuint
type AttribLocation int
type UniformLocation int
type GLclampf C.GLclampf
type GLfloat C.GLfloat

const (
	GL_LINE_SMOOTH                   = C.GL_LINE_SMOOTH
	GL_VERTEX_SHADER                 = C.GL_VERTEX_SHADER
	GL_FRAGMENT_SHADER               = C.GL_FRAGMENT_SHADER
	GL_COMPILE_STATUS                = C.GL_COMPILE_STATUS
	GL_LINK_STATUS                   = C.GL_LINK_STATUS
	GL_COLOR_BUFFER_BIT              = C.GL_COLOR_BUFFER_BIT
	GL_UNSIGNED_BYTE                 = C.GL_UNSIGNED_BYTE
	GL_DOUBLE                        = C.GL_DOUBLE
	GL_FLOAT                         = C.GL_FLOAT
	GL_ARRAY_BUFFER                  = C.GL_ARRAY_BUFFER
	GL_STATIC_DRAW                   = C.GL_STATIC_DRAW
	GL_DYNAMIC_DRAW                  = C.GL_DYNAMIC_DRAW
	GL_INFO_LOG_LENGTH               = C.GL_INFO_LOG_LENGTH
	GL_TRIANGLES                     = C.GL_TRIANGLES
	GL_LINE_LOOP                     = C.GL_LINE_LOOP
	GL_TRIANGLE_STRIP                = C.GL_TRIANGLE_STRIP
	GL_TEXTURE0                      = C.GL_TEXTURE0
	GL_TEXTURE_2D                    = C.GL_TEXTURE_2D
	GL_TEXTURE_MIN_FILTER            = C.GL_TEXTURE_MIN_FILTER
	GL_TEXTURE_MAG_FILTER            = C.GL_TEXTURE_MAG_FILTER
	GL_TEXTURE_WRAP_S                = C.GL_TEXTURE_WRAP_S
	GL_TEXTURE_WRAP_T                = C.GL_TEXTURE_WRAP_T
	GL_CLAMP_TO_EDGE                 = C.GL_CLAMP_TO_EDGE
	GL_UNPACK_ALIGNMENT              = C.GL_UNPACK_ALIGNMENT
	GL_LINEAR                        = C.GL_LINEAR
	GL_RGBA                          = C.GL_RGBA
	GL_ALPHA                         = C.GL_ALPHA
	GL_FALSE                         = C.GL_FALSE
	GL_SMOOTH_LINE_WIDTH_GRANULARITY = C.GL_SMOOTH_LINE_WIDTH_GRANULARITY
	GL_MAX_DEPTH_TEXTURE_SAMPLES     = C.GL_MAX_DEPTH_TEXTURE_SAMPLES
	GL_SRC_ALPHA                     = C.GL_SRC_ALPHA
	GL_BLEND                         = C.GL_BLEND
	GL_ONE_MINUS_SRC_ALPHA           = C.GL_ONE_MINUS_SRC_ALPHA
)
func glGetError() GLenum {
	return  GLenum(C.glGetError())
}

func glGetFloatv(flag GLenum) float32 {
	var result C.GLfloat
	C.glGetFloatv(C.GLenum(flag), &result)
	return float32(result)
}

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

func glViewport(x int, y int, width int, height int) {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(width), C.GLsizei(height))
}

func glBlendFunc(sfactor GLenum, dfactor GLenum) {
	C.glBlendFunc(C.GLenum(sfactor), C.GLenum(dfactor))
}

type Program Object

func glCreateProgram() Program {
	return Program(C.glCreateProgram())
}

func (program Program) AttachShader(shader Shader) {
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

func (program Program) GetUniformLocation(name string) UniformLocation {

	cname := glString(name)
	defer freeString(cname)

	return UniformLocation(C.glGetUniformLocation(C.GLuint(program), cname))
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

func (shader Shader) Source(source string) {
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

func (shader Shader) GetInfoLog() string {
	var length C.GLint
	C.glGetShaderiv(C.GLuint(shader), C.GLenum(GL_INFO_LOG_LENGTH), &length)
	// length is buffer size including null character

	if length > 1 {
		log := C.malloc(C.size_t(length))
		defer C.free(log)
		C.glGetShaderInfoLog(C.GLuint(shader), C.GLsizei(length), nil, (*C.GLchar)(log))
		return C.GoString((*C.char)(log))
	}
	return ""
}

func (shader Shader) Delete() {
	C.glDeleteShader(C.GLuint(shader))
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

// Uniform

func (location UniformLocation) Uniform1f(x float32) {
	C.glUniform1f(C.GLint(location), C.GLfloat(x))
}

func (location UniformLocation) Uniform1i(x int) {
	C.glUniform1i(C.GLint(location), C.GLint(x))
}

func (location UniformLocation) Uniform2f(x float32, y float32) {
	C.glUniform2f(C.GLint(location), C.GLfloat(x), C.GLfloat(y))
}

func (location UniformLocation) Uniform4fv(v Vector4) {
	C.glUniform4fv(C.GLint(location), C.GLsizei(1), (*C.GLfloat)(&v.x))
}

//GLAPI void APIENTRY glUniformMatrix4fv (GLint location, GLsizei count, GLboolean transpose, const GLfloat *value);
func (location UniformLocation) UniformMatrix4fv(m Matrix4x4) {
	C.glUniformMatrix4fv(C.GLint(location), C.GLsizei(1), GL_FALSE, (*C.GLfloat)(&m[0].x))
}

// Vertex Arrays
type VertexArray Object

func glGenVertexArray() VertexArray {
	var a C.GLuint
	C.glGenVertexArrays(1, &a)
	return VertexArray(a)
}

func glGenVertexArrays(arrays []VertexArray) {
	if len(arrays) > 0 {
		C.glGenVertexArrays(C.GLsizei(len(arrays)), (*C.GLuint)(&arrays[0]))
	}
}

func (array VertexArray) Delete() {
	C.glDeleteVertexArrays(1, (*C.GLuint)(&array))
}

func (array VertexArray) Bind() {
	C.glBindVertexArray(C.GLuint(array))
}

func (array VertexArray) Unbind() {
	C.glBindVertexArray(C.GLuint(0))
}

func glDrawArrays(mode GLenum, first int, count int) {
	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
}

//Buffera

type Buffer Object

func glGenBuffer() Buffer {
	var b C.GLuint
	C.glGenBuffers(1, &b)
	return Buffer(b)
}

func (buffer Buffer) Bind(target GLenum) {
	C.glBindBuffer(C.GLenum(target), C.GLuint(buffer))
}

func glBufferData(target GLenum, size int, data interface{}, usage GLenum) {
	C.glBufferData(C.GLenum(target), C.GLsizeiptr(size), ptr(data), C.GLenum(usage))
}
func (buffer Buffer) Unbind(target GLenum) {
	C.glBindBuffer(C.GLenum(target), C.GLuint(0))
}
func (buffer Buffer) Delete() {
	C.glDeleteBuffers(1, (*C.GLuint)(&buffer))
}
//Textures

type Texture Object

// Create single texture object
func glGenTexture() Texture {
	var b C.GLuint
	C.glGenTextures(1, &b)
	return Texture(b)
}

func (texture Texture) Bind(target GLenum) {
	C.glBindTexture(C.GLenum(target), C.GLuint(texture))
}

func glTexParameteri(target GLenum, pname GLenum, param int) {
	C.glTexParameteri(C.GLenum(target), C.GLenum(pname), C.GLint(param))
}

func glTexImage2D(target GLenum, level int, internalformat int, width int, height int, border int, format, typ GLenum, pixels interface{}) {
	C.glTexImage2D(C.GLenum(target), C.GLint(level), C.GLint(internalformat),
		C.GLsizei(width), C.GLsizei(height), C.GLint(border), C.GLenum(format),
		C.GLenum(typ), ptr(pixels))
}

func glActiveTexture(texture GLenum) { C.glActiveTexture(C.GLenum(texture)) }

func glPixelStorei(pname GLenum, param int) {
	C.glPixelStorei(C.GLenum(pname), C.GLint(param))
}

func glInit() GLenum {
	return GLenum(C.goglInit())
}
