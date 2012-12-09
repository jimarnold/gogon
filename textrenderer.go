package main

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/go-gl/glh"
	"image"
	"io/ioutil"
	"log"
	"reflect"
)

type TextRenderer struct {
	program        Program
	vs, fs         Shader
	positionAttrib AttribLocation
	colorUniform   UniformLocation
	offsetUniform  UniformLocation
	vao            VertexArray
	vbo            Buffer
	offsets        []float32
}

func NewTextRenderer(fontPath string, scale int32, dpi float64) TextRenderer {
	font := loadFont(fontPath)
	coords, texture, offsets := generateAtlas(font, scale, dpi)
	program := createProgram()

	vao := glGenVertexArray()
	vao.Bind()

	vbo := glGenBuffer()
	vbo.Bind(GL_ARRAY_BUFFER)
	glBufferData(GL_ARRAY_BUFFER, int(reflect.TypeOf(Vector4{}).Size())*len(coords), coords, GL_STATIC_DRAW)

	positionAttrib := program.GetAttribLocation("position")
	positionAttrib.AttribPointer(4, GL_FLOAT, false, 0, nil)
	positionAttrib.EnableArray()
	vbo.Unbind(GL_ARRAY_BUFFER)

	textureUniform := program.GetUniformLocation("tex")
	offsetUniform := program.GetUniformLocation("offset")
	colorUniform := program.GetUniformLocation("color")

	glActiveTexture(GL_TEXTURE0)
	tex := glGenTexture()
	tex.Bind(GL_TEXTURE_2D)
	textureUniform.Uniform1i(0)

	/* We require 1 byte alignment when uploading texture data */
	glPixelStorei(GL_UNPACK_ALIGNMENT, 1)

	/* Clamping to edges is important to prevent artifacts when scaling */
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE)

	/* Linear filtering usually looks best for text */
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR)

	glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, texture.Bounds().Dx(), texture.Bounds().Dy(), 0, GL_RGBA, GL_UNSIGNED_BYTE, texture.Pix)

	vao.Unbind()
	return TextRenderer {
	program:program,
    vao:vao,
    vbo:vbo,
    positionAttrib:positionAttrib,
    offsetUniform:offsetUniform,
    colorUniform:colorUniform,
    offsets:offsets}
}

func loadFont(fontPath string) *truetype.Font {
	b, err := ioutil.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	font, err := freetype.ParseFont(b)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return font
}

func generateAtlas(font *truetype.Font, scale int32, dpi float64) ([]Vector4, *image.RGBA, []float32) {
	var low rune = 32
	var high rune = 127
	glyphCount := int32(high-low+1)
	offsets := make([]float32, glyphCount)

	bounds := font.Bounds(scale)
	gw := bounds.XMax - bounds.XMin
	gh := bounds.YMax - bounds.YMin
	imageWidth := glh.Pow2(uint32(gw * glyphCount))
	imageHeight := glh.Pow2(uint32(gh))
	imageBounds := image.Rect(0, 0, int(imageWidth), int(imageHeight))
	sx := 2.0 / width
	sy := 2.0 / height
	w := float32(float64(gw) * sx)
	h := float32(float64(gh) * sy)
	img := image.NewRGBA(imageBounds)
	c := freetype.NewContext()
	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetSrc(image.White)
	c.SetDPI(dpi)
	c.SetFontSize(float64(scale))
	c.SetFont(font)

	var gi, gx, gy int32
	verts := make([]Vector4, 0)
	texWidth := float32(img.Bounds().Dx())
	texHeight := float32(img.Bounds().Dy())

	for ch := low; ch <= high; ch++ {
		index := font.Index(ch)
		metric := font.HMetric(scale, index)

		//the offset is used when drawing a string of glyphs - we will advance a glyph's quad by the width of all previous glyphs in the string
		offsets[gi] = float32(metric.AdvanceWidth) * float32(sx)

		//draw the glyph into the atlas at the correct location
		pt := freetype.Pt(int(gx), int(gy)+int(c.PointToFix32(float64(scale))>>8))
		c.DrawString(string(ch), pt)

		tx1 := float32(gx) / texWidth
		ty1 := float32(gy) / texHeight
		tx2 := (float32(gx) + float32(gw)) / texWidth
		ty2 := (float32(gy) + float32(gh)) / texHeight
		
		//the x,y coordinates are the same for each quad; only the texture coordinates (stored in z,w) change.
		//an optimization would be to only store texture coords, but I haven't figured that out yet
		verts = append(verts, Vector4{-1, 1, tx1, ty1},
		Vector4{-1 + (w), 1, tx2, ty1},
		Vector4{-1, 1 - (h), tx1, ty2},
		Vector4{-1 + (w), 1 - (h), tx2, ty2})

		gx += gw
		gi++
	}
	return verts, img, offsets
}

func (this *TextRenderer) Draw(s string, location, color Vector4) {
	glEnable(GL_BLEND)
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)

	this.program.Use()
	this.vao.Bind()

	this.colorUniform.Uniform4fv(color)
	totalOffset := float32(0)

	for _, ch := range s {
		index := int(ch-32)
		offset := this.offsets[index]
		this.offsetUniform.Uniform4fv(Vector4{totalOffset, 0, 0, 0}.Add(location))
		glDrawArrays(GL_TRIANGLE_STRIP, index * 4, 4)
		totalOffset += offset
	}
	this.vao.Unbind()
	this.program.Unuse()
	glDisable(GL_BLEND)
}

func (this *TextRenderer) Delete() {
	this.vs.Delete()
	this.fs.Delete()
	this.program.Delete()
	this.vbo.Delete()
	this.vao.Delete()
	log.Println(glGetError())
}

func createProgram() Program {
	vs,err := NewShader(GL_VERTEX_SHADER,`#version 150
    in vec4 position;
    out vec2 texpos;	
    uniform vec4 offset;
    void main() {
        gl_Position = vec4(position.xy, 0, 1) + offset;
		texpos = position.zw;
    }`)

	if err != nil {
		log.Printf("Error in vertex shader\n")
		log.Println(err)
	}

	fs,err := NewShader(GL_FRAGMENT_SHADER,
	`#version 150
    in vec2 texpos;
    uniform sampler2D tex;
    uniform vec4 color;
    out vec4  fragColor;
    void main(void) {
        fragColor = texture(tex, texpos) * color;	
    }`)

	if err != nil {
		log.Printf("Error in fragment shader\n")
		log.Println(err)
	}

	return NewProgram(vs, fs)
}
