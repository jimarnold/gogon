package main
import (
  "github.com/go-gl/gltext"
  "github.com/go-gl/gl"
  "os"
)

type Font struct {
  *gltext.Font
}

func(font Font) drawString(x, y float64, s string) error {
  gl.Color4f(1, 1, 1, 1)
  return font.Printf(float32(x), float32(y), s)
}

func loadFont() (Font, error) {
  file, err := os.Open("./ComingSoon.ttf")
  if err != nil {
    return Font{}, err
  }

  defer file.Close()
  scale := int32(24)
  font, err := gltext.LoadTruetype(file, scale, 32, 127, gltext.LeftToRight)
  return Font{font}, err
}
