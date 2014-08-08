package main

import "math"

type Vector2 struct {
	x, y float64
}

func (v1 Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{v1.x + v2.x, v1.y + v2.y}
}

func (v1 Vector2) Subtract(v2 Vector2) Vector2 {
	return Vector2{v1.x - v2.x, v1.y - v2.y}
}

func (v Vector2) Mult(s float64) Vector2 {
	return Vector2{v.x * s, v.y * s}
}

func (v Vector2) Length() float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y))
}

func (v Vector2) Normalized() Vector2 {
	f := 1.0 / v.Length()
	return Vector2{v.x * f, v.y * f}
}

func (v1 Vector2) DistanceTo(v2 Vector2) float64 {
	x := v2.x - v1.x
	y := v2.y - v1.y
	return math.Sqrt((x * x) + (y * y))
}

func (v Vector2) clampedTo(max float64) Vector2 {
	v.x = clamp(v.x, -max, max)
	v.y = clamp(v.y, -max, max)
	return v
}

func (v Vector2) scale(f float64) Vector2 {
	length := v.Length()
	if length == 0 {
		return v
	}
	vn := v.Normalized()
	return vn.Mult(length * f)
}

type Vector4 struct {
	x, y, z, w float32
}

func (v1 Vector4) Add(v2 Vector4) Vector4 {
	return Vector4{v1.x + v2.x, v1.y + v2.y, v1.z + v2.z, v1.w + v2.w}
}

func (this Vector4) To_a() []float32 {
	return []float32{this.x, this.y, this.z, this.w}
}
