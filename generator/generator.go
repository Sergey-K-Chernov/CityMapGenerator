package generator

import (
	"fmt"
	"math"

	"chirrwick.com/projects/city/generator/genmath"
)

func generateRadialRandomPoint(angle_min, angle_max, radius_min, raduis_max float64) genmath.Point {
	angle := genmath.RandFloat(angle_min, angle_max)
	radius := genmath.RandFloat(radius_min, raduis_max)

	return genmath.Point{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
}

func shiftPoints(points []genmath.Point) (shift genmath.Point) {
	for _, point := range points {
		shift.X = math.Min(shift.X, point.X)
		shift.Y = math.Min(shift.Y, point.Y)
	}
	shift.X = -shift.X + 100
	shift.Y = -shift.Y + 100

	for i := range points {
		points[i].X += shift.X
		points[i].Y += shift.Y
	}

	return
}

func GenerateBorders(chan_map chan Map, initials InitialValues) {
	nPoints := initials.NumSides
	rMin := initials.Raduis.Min
	rMax := initials.Raduis.Max

	angle_step := genmath.DegToRad(360. / float64(nPoints))
	angle_variation := angle_step / 2

	var cityMap Map
	cityMap.BorderPoints = make([]genmath.Point, initials.NumSides)

	for i := 0; i < nPoints; i++ {
		angle := angle_step * float64(i)
		cityMap.BorderPoints[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, rMin, rMax)

		point := generateRadialRandomPoint(0, 2*math.Pi, 0, initials.VertexShift)
		cityMap.BorderPoints[i].Add(point)
	}

	shift := shiftPoints(cityMap.BorderPoints)

	cityMap.Center = shift

	for _, point := range cityMap.BorderPoints {
		fmt.Printf("%7.1f\t%7.1f\n", point.X, point.Y)
	}

	fmt.Printf("\nShift: %7.1f\t%7.1f\n", shift.X, shift.Y)

	chan_map <- cityMap

}
