package generator

import (
	"errors"
	"math"
	"slices"

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

func calcPolygonArea(polygon []gm.Point, center gm.Point) float64 {
	area := 0.0
	for i := range polygon {
		i_plus_1 := i + 1
		if i_plus_1 == len(polygon) {
			i_plus_1 = 0
		}

		// Heron's formula
		a := polygon[i].Sub(polygon[i_plus_1]).Length()
		b := polygon[i].Sub(center).Length()
		c := polygon[i_plus_1].Sub(center).Length()
		p := (a + b + c) / 2

		area += math.Sqrt(p * (p - a) * (p - b) * (p - c))
	}
	return area
}

func intersectSegmentWithFigure(segment gm.LineSegment, figure []gm.Point) (gm.Point, gm.LineSegment, error) {

	for i := 0; i < len(figure); i++ {
		i_plus_1 := i + 1
		if i_plus_1 == len(figure) {
			i_plus_1 = 0
		}
		figureSegment := gm.LineSegment{Begin: figure[i], End: figure[i_plus_1]}

		if p, ok := segment.Intersect(figureSegment); ok {
			return p, figureSegment, nil
		}
	}

	return gm.Point{}, gm.LineSegment{}, errors.New("no intersection with figure")
}

func removePointsInsideFigure(points, figure []gm.Point) []gm.Point {
	for i := 0; i < len(points); {
		if checkPointInsidePolygon(points[i], figure) {
			points = slices.Delete(points, i, i+1)
		} else {
			i++
		}
	}
	return points
}

func checkPointInsidePolygon(p gm.Point, poly []gm.Point) bool {
	sum := 0
	for i := range poly {
		i_plus_1 := i + 1
		if i_plus_1 == len(poly) {
			i_plus_1 = 0
		}

		sum += checkXRayIntersectsSegment(p, poly[i], poly[i_plus_1])
	}
	return sum%2 > 0
}

func checkXRayIntersectsSegment(p, a, b gm.Point) int {
	ax := a.X - p.X
	ay := a.Y - p.Y

	bx := b.X - p.X
	by := b.Y - p.Y

	intersect_axis_x := (ay*by <= 0)
	intersect_axis_y := (ax*bx <= 0)

	// Send ray to the right, test for intersecting positive x axis .

	// No intersection at all:
	if !intersect_axis_x {
		return 0
	}

	if !intersect_axis_y {
		// Both to the left - no intersection
		if ax < 0 {
			return 0
		}
		// Both to the right - have intersection
		return 1
	}

	// find intersection point
	fraction := math.Abs(ay / (by - ay))
	x := ax + fraction*(bx-ax)

	if x < 0 {
		return 0
	}

	return 1
}

type Intersection struct {
	i, i_next int
	point     gm.Point
}

func cutFigure(center gm.Point, figure []gm.Point, max_radius float64, segment gm.LineSegment) []gm.Point {
	// Находим перпендикуляр из центра на отрезок
	np := segment.GetNormalPoint(center)
	n := np.Sub(center)

	// Если больше расстояния до самого дальнего угла, игнорим
	normal_length := n.Length()
	if normal_length > max_radius {
		return figure
	}

	// Находим пересечения отрезка с обеими сторонами
	intersections := make([]Intersection, 0)
	for i, p := range figure {
		i_next := i + 1
		if i_next == len(figure) {
			i_next = 0
		}

		s := gm.LineSegment{Begin: p, End: figure[i_next]}
		point, ok := s.IntersectLine(segment)
		if ok {
			intersections = append(intersections, Intersection{i: i, i_next: i_next, point: point})
		}
	}

	if len(intersections) < 2 {
		return figure
	}

	// Проверяем, находятся ли точки перечечения внутри отрезка
	line_vec := segment.End.Sub(segment.Begin)
	line_vec_ort := line_vec.GetNormalized()
	ok := false
	for i := 0; i < len(intersections); i++ {
		isec_vec := intersections[i].point.Sub(segment.Begin)

		projection := line_vec_ort.Dot(isec_vec)
		if projection > 0 && projection < line_vec.Length() {
			ok = true
			break
		}
	}

	if !ok {
		return figure
	}

	slices.SortFunc(intersections, func(a, b Intersection) int {
		if a.i_next < b.i_next {
			return 1
		}
		if a.i_next > b.i_next {
			return -1
		}
		return 0
	})

	// Вставляем точки в фигуру
	for _, isec := range intersections {
		figure = slices.Insert(figure, isec.i_next, isec.point)
	}

	// Находим проекции каждой точки на перпендикуляр. Если отрицательные, игнорим

	ortho := gm.Line{Origin: center, Vector: n.GetNormalized()}
	nlen := n.Length()
	for i := len(figure) - 1; i >= 0; i-- {
		vec_to_corner := figure[i].Sub(center)
		projection := ortho.Vector.Dot(vec_to_corner)
		if gm.AlmostEqual(projection, nlen, 0.0000001) {
			continue
		}
		if projection > nlen { // < 0 is automacitally < length
			figure = slices.Delete(figure, i, i+1)
		}
	}

	return figure
}
