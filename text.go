package main
//import (
  ////"github.com/go-gl/gltext"
  //"os"
	//"code.google.com/p/freetype-go/freetype"
	//"code.google.com/p/freetype-go/freetype/truetype"
	//"github.com/go-gl/glh"
	//"image"
	//"io"
	//"io/ioutil"
//)

//// A Glyph describes metrics for a single font glyph.
//// These indicate which area of a given image contains the
//// glyph data and how the glyph should be spaced in a rendered string.
//type Glyph struct {
	//X      int `json:"x"`      // The x location of the glyph on a sprite sheet.
	//Y      int `json:"y"`      // The y location of the glyph on a sprite sheet.
	//Width  int `json:"width"`  // The width of the glyph on a sprite sheet.
	//Height int `json:"height"` // The height of the glyph on a sprite sheet.

	//// Advance determines the distance to the next glyph.
	//// This is used to properly align non-monospaced fonts.
	//Advance int `json:"advance"`
//}

//// A Charset represents a set of glyph descriptors for a font.
//// Each glyph descriptor holds glyph metrics which are used to
//// properly align the given glyph in the resulting rendered string.
//type Charset []Glyph

//// A Font allows rendering of text to an OpenGL context.
//type Font struct {
	//config         *FontConfig // Character set for this font.
	//texture        Texture  // Holds the glyph texture id.
	//listbase       uint        // Holds the first display list id.
	//maxGlyphWidth  int         // Largest glyph width.
	//maxGlyphHeight int         // Largest glyph height.
//}

//const (
	//LeftToRight Direction = iota // E.g.: Latin
	//RightToLeft                  // E.g.: Arabic
	//TopToBottom                  // E.g.: Chinese
//)

//type MyFont struct {
	//*Font
//}
//func(font MyFont) drawString(x, y float64, s string) error {
  //return font.Printf(float32(x), float32(y), s)
//}

//func myLoadFont() (MyFont, error) {
  //file, err := os.Open("./ComingSoon.ttf")
  //if err != nil {
    //return MyFont{}, err
  //}

  //defer file.Close()
  //scale := int32(24)
  //font, err := LoadTruetype(file, scale, 32, 127, LeftToRight)
  //return MyFont{font}, err
//}

//func (f *Font) Printf(x, y float32, fs string, argv ...interface{}) error {
	//return glh.CheckGLError()
//}

//type Direction uint8

//type FontConfig struct {
	//// The direction determines the orientation of rendered strings and should 
	//// hold any of the pre-defined Direction constants.
	//Dir Direction `json:"direction"`

	//// Lower rune boundary
	//Low rune `json:"rune_low"`

	//// Upper rune boundary.
	//High rune `json:"rune_high"`

	//// Glyphs holds a set of glyph descriptors, defining the location,
	//// size and advance of each glyph in the sprite sheet.
	//Glyphs Charset `json:"glyphs"`
//}

//func LoadTruetype(r io.Reader, scale int32, low, high rune, dir Direction) (*Font, error) {
	//data, err := ioutil.ReadAll(r)
	//if err != nil {
		//return nil, err
	//}

	//// Read the truetype font.
	//ttf, err := truetype.Parse(data)
	//if err != nil {
		//return nil, err
	//}

	//// Create our FontConfig type.
	//var fc FontConfig
	//fc.Dir = dir
	//fc.Low = low
	//fc.High = high
	//fc.Glyphs = make(Charset, high-low+1)

	//// Create an image, large enough to store all requested glyphs.
	////
	//// We limit the image to 16 glyphs per row. Then add as many rows as
	//// needed to encompass all glyphs, while making sure the resulting image
	//// has power-of-two dimensions.
	//gc := int32(len(fc.Glyphs))
	//glyphsPerRow := int32(16)
	//glyphsPerCol := (gc / glyphsPerRow) + 1

	//gb := ttf.Bounds(scale)
	//gw := (gb.XMax - gb.XMin)
	//gh := (gb.YMax - gb.YMin) + 5
	//iw := glh.Pow2(uint32(gw * glyphsPerRow))
	//ih := glh.Pow2(uint32(gh * glyphsPerCol))

	//rect := image.Rect(0, 0, int(iw), int(ih))
	//img := image.NewRGBA(rect)

	//// Use a freetype context to do the drawing.
	//c := freetype.NewContext()
	//c.SetDPI(72)
	//c.SetFont(ttf)
	//c.SetFontSize(float64(scale))
	//c.SetClip(img.Bounds())
	//c.SetDst(img)
	//c.SetSrc(image.White)

	//// Iterate over all relevant glyphs in the truetype font and
	//// draw them all to the image buffer.
	////
	//// For each glyph, we also create a corresponding Glyph structure
	//// for our Charset. It contains the appropriate glyph coordinate offsets.
	//var gi int
	//var gx, gy int32

	//for ch := low; ch <= high; ch++ {
		//index := ttf.Index(ch)
		//metric := ttf.HMetric(scale, index)

		//fc.Glyphs[gi].Advance = int(metric.AdvanceWidth)
		//fc.Glyphs[gi].X = int(gx)
		//fc.Glyphs[gi].Y = int(gy)
		//fc.Glyphs[gi].Width = int(gw)
		//fc.Glyphs[gi].Height = int(gh)

		//pt := freetype.Pt(int(gx), int(gy)+int(c.PointToFix32(float64(scale))>>8))
		//c.DrawString(string(ch), pt)

		//if gi%16 == 0 {
			//gx = 0
			//gy += gh
		//} else {
			//gx += gw
		//}

		//gi++
	//}

	//return loadFont(img, &fc)
//}

//// loadFont loads the given font data. This does not deal with font scaling.
//// Scaling should be handled by the independent Bitmap/Truetype loaders.
//// We therefore expect the supplied image and charset to already be adjusted
//// to the correct font scale.
////
//// The image should hold a sprite sheet, defining the graphical layout for
//// every glyph. The config describes font metadata.
//func loadFont(img *image.RGBA, config *FontConfig) (f *Font, err error) {
	//f = new(Font)
	//f.config = config

	//// Resize image to next power-of-two.
	//img = glh.Pow2Image(img).(*image.RGBA)
	//ib := img.Bounds()

	//// Create the texture itself. It will contain all glyphs.
	//// Individual glyph-quads display a subset of this texture.
	//f.texture = glGenTexture()
	//f.texture.Bind(GL_TEXTURE_2D)
	//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR)
	//glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR)
	//glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, ib.Dx(), ib.Dy(), 0,
		//GL_RGBA, GL_UNSIGNED_BYTE, img.Pix)

	//// Create display lists for each glyph.
	//f.listbase = glGenLists(len(config.Glyphs))

	//texWidth := float32(ib.Dx())
	//texHeight := float32(ib.Dy())

	//for index, glyph := range config.Glyphs {
		//// Update max glyph bounds.
		//if glyph.Width > f.maxGlyphWidth {
			//f.maxGlyphWidth = glyph.Width
		//}

		//if glyph.Height > f.maxGlyphHeight {
			//f.maxGlyphHeight = glyph.Height
		//}

		//// Quad width/height
		//vw := float32(glyph.Width)
		//vh := float32(glyph.Height)

		//// Texture coordinate offsets.
		//tx1 := float32(glyph.X) / texWidth
		//ty1 := float32(glyph.Y) / texHeight
		//tx2 := (float32(glyph.X) + vw) / texWidth
		//ty2 := (float32(glyph.Y) + vh) / texHeight

		//// Advance width (or height if we render top-to-bottom)
		//adv := float32(glyph.Advance)

		//glNewList(f.listbase+uint(index), GL_COMPILE)
		//{
			//glBegin(GL_QUADS)
			//{
				//glTexCoord2f(tx1, ty2)
				//glVertex2f(0, 0)
				//glTexCoord2f(tx2, ty2)
				//glVertex2f(vw, 0)
				//glTexCoord2f(tx2, ty1)
				//glVertex2f(vw, vh)
				//glTexCoord2f(tx1, ty1)
				//glVertex2f(0, vh)
			//}
			//glEnd()

			//switch config.Dir {
			//case LeftToRight:
				//glTranslatef(adv, 0, 0)
			//case RightToLeft:
				//glTranslatef(-adv, 0, 0)
			//case TopToBottom:
				//glTranslatef(0, -adv, 0)
			//}
		//}
		//glEndList()
	//}

	//err = glh.CheckGLError()
	//return
//}
