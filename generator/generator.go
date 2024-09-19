package main

import (
	"fmt"
	"math"
	"math/rand"
)

func almostEqual(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

func degToRad(value float64) float64 {
	return value * math.Pi / 180.
}

func randFloat(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

type Point struct {
	X float64
	Y float64
}

func generateRadialRandomPoint(angle_min, angle_max, radius_min, raduis_max float64) Point {
	angle := randFloat(angle_min, angle_max)
	radius := randFloat(radius_min, raduis_max)

	return Point{radius * math.Cos(angle), radius * math.Sin(angle)}
}

func generateBorders(n_points int, size_rough_min, size_rough_max float64) {
	angle_step := degToRad(360. / float64(n_points))
	angle_variation := angle_step / 2

	points := make([]Point, n_points)

	for i := 0; i < n_points; i++ {
		angle := angle_step * float64(i)
		points[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, size_rough_min/2, size_rough_max/2)

		fmt.Printf("%7.1f\t%7.1f\n", points[i].X, points[i].Y)
	}
}

func main() {
	generateBorders(6, 2000., 3000.)
}
