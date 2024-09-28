package generator

import (
	"errors"
	"math"

	gm "chirrwick.com/projects/city/generator/genmath"
)

func generateRadialRandomPoint(angle_min, angle_max, radius_min, raduis_max float64) gm.Point {
	angle := gm.RandFloat(angle_min, angle_max)
	radius := gm.RandFloat(radius_min, raduis_max)

	return gm.Point{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)}
}

func shiftPoints(points []gm.Point) (shift gm.Point) {
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

func extend(segment *gm.LineSegment, figure []gm.Vector2D) {
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

func intersectWithFigure(segment gm.LineSegment, figure []gm.Point) (gm.Point, gm.LineSegment, error) {

	for i := 0; i < len(figure); i++ {
		i_plus_1 := i + 1
		if i_plus_1 == len(figure) {
			i_plus_1 = 0
		}
		borderSegment := gm.LineSegment{Begin: figure[i], End: figure[i_plus_1]}

		if p, ok := segment.Intersect(borderSegment); ok {
			return p, borderSegment, nil
		}
	}

	return gm.Point{}, gm.LineSegment{}, errors.New("no intersection with figure")
}

func findClosestPointIndex(point gm.Point, points []gm.Point) int {
	d_sq_min := math.MaxFloat64
	index := 0
	for i, p := range points {
		d_sq := p.Sub(point).LengthSq()
		if d_sq_min > d_sq {
			d_sq_min = d_sq
			index = i
		}
	}
	return index
}
