package generator

import (
	"fmt"
	"math"

	gm "chirrwick.com/projects/city/generator/genmath"
)

func GenerateBorders(chanMap chan Map, initials InitialValuesMap) {
	nPoints := initials.NumSides
	rMin := initials.Raduis.Min
	rMax := initials.Raduis.Max

	angle_step := gm.DegToRad(360. / float64(nPoints))
	angle_variation := angle_step / 2

	var cityMap Map
	cityMap.BorderPoints = make([]gm.Point, initials.NumSides)

	for i := 0; i < nPoints; i++ {
		angle := angle_step * float64(i)
		cityMap.BorderPoints[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, rMin, rMax)

		point := generateRadialRandomPoint(0, 2*math.Pi, 0, initials.VertexShift)
		cityMap.BorderPoints[i].AddInPlace(point)
	}

	shift := shiftPoints(cityMap.BorderPoints)

	cityMap.Center = shift

	for _, point := range cityMap.BorderPoints {
		fmt.Printf("%7.1f\t%7.1f\n", point.X, point.Y)
	}

	fmt.Printf("\nShift: %7.1f\t%7.1f\n", shift.X, shift.Y)

	chanMap <- cityMap
}
