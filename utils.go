package fine

// 2-dimensional vector.
type Vec2 struct {
	X float64
	Y float64
}

// NewVec2 is a shorthand for Vec2{X: n, Y: n}.
func NewVec2(x, y float64) Vec2 {
	return Vec2{x, y}
}

// Subtracts vectors.
func (v *Vec2) Sub(v2 Vec2) Vec2 {
	return NewVec2(v.X-v2.X, v.Y-v2.Y)
}

// Adds vectors.
func (v *Vec2) Add(v2 Vec2) Vec2 {
	return NewVec2(v.X+v2.X, v.Y+v2.Y)
}

// Multiplies vectors.
func (v *Vec2) Mul(v2 Vec2) Vec2 {
	return NewVec2(v.X*v2.X, v.Y*v2.Y)
}

// Divides vectors.
func (v *Vec2) Div(v2 Vec2) Vec2 {
	return NewVec2(v.X/v2.X, v.Y/v2.Y)
}
