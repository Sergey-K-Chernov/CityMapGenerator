package main

import (
	"fmt"
	"math"

	"chirrwick.com/projects/city/generator/genmath"
)

func generateRadialRandomPoint(angle_min, angle_max, radius_min, raduis_max float64) genmath.Point {
	angle := genmath.RandFloat(angle_min, angle_max)
	radius := genmath.RandFloat(radius_min, raduis_max)

	return genmath.Point{radius * math.Cos(angle), radius * math.Sin(angle)}
}

func GenerateBorders(n_points int, size_rough_min, size_rough_max float64) {
	angle_step := genmath.DegToRad(360. / float64(n_points))
	angle_variation := angle_step / 2

	points := make([]genmath.Point, n_points)

	for i := 0; i < n_points; i++ {
		angle := angle_step * float64(i)
		points[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, size_rough_min/2, size_rough_max/2)

		point := generateRadialRandomPoint(0, 2*math.Pi, 0, (size_rough_min+size_rough_max)/20.0)
		points[i].Add(point)

		fmt.Printf("%7.1f\t%7.1f\n", points[i].X, points[i].Y)
	}
}

func main() {
	GenerateBorders(6, 2000., 3000.)
}
