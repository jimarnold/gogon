package main

type Matrix4x4 [4]Vector4

func NewMatrix4x4(d float32) Matrix4x4 {
	return Matrix4x4{
		Vector4{d, 0, 0, 0},
		Vector4{0, d, 0, 0},
		Vector4{0, 0, d, 0},
		Vector4{0, 0, 0, d}}
}

func (m Matrix4x4) toa() [16]float32 {
	return [16]float32{
		m[0].x,
		m[0].y,
		m[0].z,
		m[0].w,
		m[1].x,
		m[1].y,
		m[1].z,
		m[1].w,
		m[2].x,
		m[2].y,
		m[2].z,
		m[2].w,
		m[3].x,
		m[3].y,
		m[3].z,
		m[3].w}
}

func (a Matrix4x4) mult(b Matrix4x4) Matrix4x4 {
	a0 := a[0]
	a1 := a[1]
	a2 := a[2]
	a3 := a[3]
	b0 := b[0]
	b1 := b[1]
	b2 := b[2]
	b3 := b[3]
	return Matrix4x4{
		Vector4{
			a0.x*b0.x + a1.x*b0.y + a2.x*b0.z + a3.x*b0.w,
			a0.y*b0.x + a1.y*b0.y + a2.y*b0.z + a3.y*b0.w,
			a0.z*b0.x + a1.z*b0.y + a2.z*b0.z + a3.z*b0.w,
			a0.w*b0.x + a1.w*b0.y + a2.w*b0.z + a3.w*b0.w},
		Vector4{
			a0.x*b1.x + a1.x*b1.y + a2.x*b1.z + a3.x*b1.w,
			a0.y*b1.x + a1.y*b1.y + a2.y*b1.z + a3.y*b1.w,
			a0.z*b1.x + a1.z*b1.y + a2.z*b1.z + a3.z*b1.w,
			a0.w*b1.x + a1.w*b1.y + a2.w*b1.z + a3.w*b1.w},
		Vector4{
			a0.x*b2.x + a1.x*b2.y + a2.x*b2.z + a3.x*b2.w,
			a0.y*b2.x + a1.y*b2.y + a2.y*b2.z + a3.y*b2.w,
			a0.z*b2.x + a1.z*b2.y + a2.z*b2.z + a3.z*b2.w,
			a0.w*b2.x + a1.w*b2.y + a2.w*b2.z + a3.w*b2.w},
		Vector4{
			a0.x*b3.x + a1.x*b3.y + a2.x*b3.z + a3.x*b3.w,
			a0.y*b3.x + a1.y*b3.y + a2.y*b3.z + a3.y*b3.w,
			a0.z*b3.x + a1.z*b3.y + a2.z*b3.z + a3.z*b3.w,
			a0.w*b3.x + a1.w*b3.y + a2.w*b3.z + a3.w*b3.w}}
}

func ortho(left, right, bottom, top, zNear, zFar float32) Matrix4x4 {
	result := NewMatrix4x4(1.0)
	result[0].x = float32(2.0) / (right - left)
	result[1].y = float32(2.0) / (top - bottom)
	result[2].z = float32(-2.0) / (zFar - zNear)
	result[3].x = -(right + left) / (right - left)
	result[3].y = -(top + bottom) / (top - bottom)
	result[3].z = -(zFar + zNear) / (zFar - zNear)
	return result
}
