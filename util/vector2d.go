package util

import "math"

func Vector(x, y float64) *Vector2D {
	return &Vector2D{X: x, Y: y}
}

type Vector2D struct {
	X, Y float64
}

func (v *Vector2D) Copy() *Vector2D {
	return &Vector2D{X: v.X, Y: v.Y}
}

func (v *Vector2D) SquaredMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector2D) Direction() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v *Vector2D) Normalize() *Vector2D {
	if v.SquaredMagnitude() == 0 {
		return v
	}

	magnitude := v.Magnitude()
	v.X /= magnitude
	v.Y /= magnitude
	return v
}

func (v *Vector2D) Add(other *Vector2D) *Vector2D {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v *Vector2D) Subtract(other *Vector2D) *Vector2D {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v *Vector2D) Scale(scalar float64) *Vector2D {
	v.X *= scalar
	v.Y *= scalar
	return v
}

func (v *Vector2D) Dot(other *Vector2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v *Vector2D) Cross(other *Vector2D) float64 {
	return v.X*other.Y - v.Y*other.X
}

func (v *Vector2D) AngleBetween(other *Vector2D) float64 {
	dot := v.Dot(other)
	magnitudeProduct := v.Magnitude() * other.Magnitude()
	if magnitudeProduct == 0 {
		return 0
	}
	return math.Acos(dot / magnitudeProduct)
}

func (v *Vector2D) Rotate(angle float64) *Vector2D {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	x := v.X*cos - v.Y*sin
	y := v.X*sin + v.Y*cos
	v.X = x
	v.Y = y
	return v
}
