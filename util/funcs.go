package util

import (
	"image/color"
	"math"
	"math/rand/v2"
)

// Smoothly interpolate between two values
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Smoothly interpolate between two angles in radians
func LerpAngle(a, b, t float64) float64 {
	cs := (1-t)*math.Cos(a) + t*math.Cos(b)
	sn := (1-t)*math.Sin(a) + t*math.Sin(b)
	return math.Atan2(sn, cs)
}

// Get the squared distance between two points
//
// This is more efficient than using the Distance function if you don't need the actual distance
func SquaredDistance(a, b *Vector2D) float64 {
	return (b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y)
}

// Get the distance between two points
//
// Prefer SquaredDistance for comparisons between distances
func Distance(a, b *Vector2D) float64 {
	return math.Sqrt((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y))
}

// Get the angle between two vectors in radians
//
// The result is in the range [-π, π]
func AngleBetween(source, destination *Vector2D) float64 {
	return math.Atan2(destination.Y-source.Y, destination.X-source.X)
}

// RandomRange returns a random float64 between min and max
//
// If min is greater than max, they are swapped
func RandomRange(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}

	return min + (max-min)*rand.Float64()
}

// RandomRangeInt returns a random int between min and max
//
// If min is greater than max, they are swapped
func RandomRangeInt(min, max int) int {
	if min > max {
		min, max = max, min
	}

	return min + rand.IntN(max-min+1)
}

// AngleDifference returns the smallest difference between two angles in radians
//
// The result is in the range [-π, π]
func AngleDifference(a, b float64) float64 {
	return math.Atan2(math.Sin(a-b), math.Cos(a-b))
}

func RandomFireColor() color.RGBA {
	if rand.Float64() > .5 {
		return color.RGBA{
			R: uint8(rand.IntN(96) + 96),
			G: uint8(rand.IntN(95)),
			B: uint8(rand.IntN(24)),
			A: uint8(rand.IntN(48) + 128),
		}
	}

	return color.RGBA{
		R: uint8(rand.IntN(32) + 48),
		G: uint8(rand.IntN(32) + 48),
		B: uint8(rand.IntN(32) + 48),
		A: uint8(rand.IntN(32) + 96),
	}
}

func NearColorHueOnly(col color.Color, randomThresh float64) color.RGBA {
	r16, g16, b16, a16 := col.RGBA()
	if a16 == 0 {
		return color.RGBA{R: 255, G: 255, B: 255, A: 0}
	}

	const inv65535 = 1.0 / 65535.0
	const sqrt3 = 1.7320508075688772
	a := float64(a16) * inv65535
	ra := float64(r16) * inv65535 / a
	ga := float64(g16) * inv65535 / a
	ba := float64(b16) * inv65535 / a
	L := (ra + ga + ba) / 3.0
	X := 2.0*ra - ga - ba
	Y := sqrt3 * (ga - ba)
	d := (rand.Float64() - 0.5) * randomThresh
	s, c := math.Sincos(d)
	Xp := X*c - Y*s
	Yp := X*s + Y*c
	inv2sqrt3 := 1.0 / (2.0 * sqrt3)
	r := L + Xp/3.0
	g := L - Xp/6.0 + Yp*inv2sqrt3
	b := L - Xp/6.0 - Yp*inv2sqrt3

	if r < 0 {
		r = 0
	} else if r > 1 {
		r = 1
	}
	if g < 0 {
		g = 0
	} else if g > 1 {
		g = 1
	}
	if b < 0 {
		b = 0
	} else if b > 1 {
		b = 1
	}

	R := uint8(r*a*255.0 + 0.5)
	G := uint8(g*a*255.0 + 0.5)
	B := uint8(b*a*255.0 + 0.5)
	A := uint8(a*255.0 + 0.5)

	return color.RGBA{R: R, G: G, B: B, A: A}
}

func RandomAngularVector(minStrength, maxStrength float64) *Vector2D {
	// strength := RandomRange(minStrength, maxStrength)
	// angle := rand.Float64() * 2 * math.Pi
	// return strength * math.Cos(angle), strength * math.Sin(angle)

	var (
		angle    float64 = rand.Float64() * 2 * math.Pi
		strength float64 = RandomRange(minStrength, maxStrength)
	)

	return &Vector2D{
		X: strength * math.Cos(angle),
		Y: strength * math.Sin(angle),
	}
}

func RandomRadius(radius float64) *Vector2D {
	var (
		angle    float64 = rand.Float64() * 2 * math.Pi
		distance float64 = rand.Float64() * radius
	)

	return &Vector2D{
		X: distance * math.Cos(angle),
		Y: distance * math.Sin(angle),
	}
}
