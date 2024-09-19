package main

import (
	"fmt"
	"math"

	"chirrwick.com/projects/city/generator/genmath"
)

func generateRadialRandomPoint(angle_min, angle_max, radius_min, raduis_max float64) genmath.Point {
	angle := genmath.RandFloat(angle_min, angle_max)
	radius := genmath.RandFloat(radius_min, raduis_max)

	var res genmath.Point
	res.X = radius * math.Cos(angle)
	res.Y = radius * math.Sin(angle)
	return res
}

func shiftPoints(points []genmath.Point) (shift genmath.Point) {
	for _, point := range points {
		shift.X = math.Min(shift.X, point.X)
		shift.Y = math.Min(shift.Y, point.Y)
	}
	shift.X = -shift.X + 100
	shift.Y = -shift.Y + 100

	for i, _ := range points {
		points[i].X += shift.X
		points[i].Y += shift.Y
	}

	return
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
	}

	shift := shiftPoints(points)

	for _, point := range points {
		fmt.Printf("%7.1f\t%7.1f\n", point.X, point.Y)
	}

	fmt.Printf("\nShift: %7.1f\t%7.1f\n", shift.X, shift.Y)
}

func main() {
	GenerateBorders(6, 2000., 3000.)
}
