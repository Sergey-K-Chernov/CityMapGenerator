package genmath

import (
	"math"
	"math/rand/v2"
)

func AlmostEqual(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

func DegToRad(value float64) float64 {
	return value * math.Pi / 180.
}

func RandFloat(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

type Point struct {
	X float64
	Y float64
}

func (p1 *Point) Add(p2 Point) {
	p1.X += p2.X
	p1.Y += p2.Y
}
