package main

import "math"
import "github.com/go-gl/gl"

const TWO_PI = 2.0 * math.Pi

func drawCircle(location Vector, radius float64, color RGB) {
  gl.Color3d(color.r, color.g, color.b)
  gl.Begin(gl.LINE_LOOP)
    sides := radius * 2.0
    scale := 1.0 / sides
    for i := 0.0; i < sides; i++ {
      angle := i * TWO_PI * scale
      gl.Vertex2d(location.x + (math.Cos(angle) * radius), location.y + (math.Sin(angle) * radius))
    }
  gl.End()
}
