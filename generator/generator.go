package generator

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"

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

func GenerateBorders(chanMap chan Map, initials InitialValuesMap) {
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

func GenerateRoads(cityMap Map, chanMap chan Map, initials InitialValuesRoads) (roads []Road) {
	nCenters := initials.NumCenters
	rMin := initials.Raduis.Min
	rMax := initials.Raduis.Max

	angle_step := genmath.DegToRad(360. / float64(nCenters))
	angle_variation := angle_step / 4

	centers := make([]genmath.Point, nCenters)

	for i := 0; i < nCenters; i++ {
		angle := angle_step * float64(i)
		centers[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, rMin, rMax)
		centers[i].AddInPlace(cityMap.Center)
	}

	for i := 0; i < nCenters-1; i++ {
		var rd Road
		rd.Points = append(rd.Points, centers[i])
		for j := i + 1; j < nCenters; j++ {
			rd.Points = append(rd.Points, centers[j])
		}
		roads = append(roads, rd)
	}

	for i := 0; i < nCenters; i++ {
		outSegment := genmath.LineSegment{Begin: cityMap.Center, End: centers[i]}
		extend(&outSegment, cityMap.BorderPoints)
		if _, borderSeg, err := intersectWithFigure(outSegment, cityMap.BorderPoints); err != nil {
			println(err)
		} else {
			roads = append(roads, GenerateOutsideRoad(centers[i], borderSeg, initials.Branching)...)
		}
	}

	return
}

func extend(segment *genmath.LineSegment, figure []genmath.Vector2D) {
	maxLengthSq := 0.0
	for _, s := range figure {
		vecToBorder := segment.Begin.Sub(s)
		maxLengthSq = math.Max(maxLengthSq, vecToBorder.LengthSq())
	}
	maxLength := math.Sqrt(maxLengthSq)

	if maxLengthSq > segment.LengthSq() {
		vec := segment.GetVector()
		factor := maxLength / segment.Length() * 1.1 // 1.1 to be sure.
		vec.Scale(factor)
		segment.End = segment.Begin.Add(vec)
	}
}

func intersectWithFigure(segment genmath.LineSegment, figure []genmath.Point) (genmath.Point, genmath.LineSegment, error) {
	for i := 0; i < len(figure)-1; i++ {
		borderSegment := genmath.LineSegment{Begin: figure[i], End: figure[i+1]}

		if p, ok := segment.Intersect(borderSegment); ok {
			return p, borderSegment, nil
		}
	}

	return genmath.Point{}, genmath.LineSegment{}, errors.New("no intersection with figure")
}

func GenerateOutsideRoad(start genmath.Point, border genmath.LineSegment, branching int) []Road {
	roads := make([]Road, 0)

	splitRatio := genmath.RandFloat(0.3, 0.7)
	border1, border2 := border.Split(splitRatio)

	if branching == 0 {
		var rd Road
		rd.Points = append(rd.Points, start, border1.End)
		roads = append(roads, rd)
		return roads
	}

	vecToBorder := border1.End.Sub(start)
	splitDistance := vecToBorder.Length() / float64(branching)
	splitDistance += genmath.RandFloat(-splitDistance/4, splitDistance/4)

	vecToBorder = *vecToBorder.Normalize().Scale(splitDistance)

	end := start.Add(vecToBorder)

	var rd Road
	rd.Points = append(rd.Points, start, end)

	roads = append(roads, rd)

	roads = append(roads, GenerateOutsideRoad(end, border1, branching-1)...)
	if rand.Float64() > 0.5 {
		roads = append(roads, GenerateOutsideRoad(end, border2, branching-1)...)
	}

	return roads
}
