package generator

import (
	"errors"
	"math"
	"slices"

	gm "chirrwick.com/projects/city/generator/genmath"
)

type Range struct {
	Min float64
	Max float64
}

type InitialValuesMap struct {
	Raduis      Range
	NumSides    int
	VertexShift float64
}

type InitialValuesRoads struct {
	NumCenters int
	Raduis     Range
	Branching  int
}

type InitialValuesAreas struct {
	NumIndustrial  int
	AreaIndustrial float64
	NumParks       int
	AreaParks      float64
}

type InitialValuesBlocks struct {
	Size Range
}

type Intersection struct {
	i, i_next int
	point     gm.Point
}

type Fan struct {
	root  gm.Point
	outer []Fan
}

func makeFan(root gm.Point) (f Fan) {
	f.root = root
	f.outer = make([]Fan, 0)
	return
}

func (f *Fan) split() {
	n := len(f.outer)
	if n < 2 {
		return
	}

	new_outer := make([]Fan, 0)
	for i := 0; i < n-1; i += 2 {
		p1 := f.outer[i]
		p2 := f.outer[i+1]

		min_radial_ratio, max_radial_ratio := 0.3, 0.7
		min_dist_ratio, max_dist_ratio := 1.0/math.Log2(float64(n))*0.6, 1.0/math.Log2(float64(n))*0.9

		radial_ratio := gm.RandFloat(min_radial_ratio, max_radial_ratio)
		dist_ratio := 1 - gm.RandFloat(min_dist_ratio, max_dist_ratio)

		segment := gm.LineSegment{Begin: p1.root, End: p2.root}
		segment, _ = segment.Split(radial_ratio)

		segment = gm.LineSegment{Begin: f.root, End: segment.End}
		segment, _ = segment.Split(dist_ratio)

		fan := makeFan(segment.End)
		fan.outer = append(fan.outer, p1, p2)
		new_outer = append(new_outer, fan)
	}
	if len(f.outer)%2 != 0 {
		new_outer = append(new_outer, f.outer[len(f.outer)-1])
	}

	f.outer = new_outer

	f.split()
}

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

func extendSegment(segment *gm.LineSegment, figure []gm.Vector2D) {
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

func intersectExtendedSegmentWithFigure(figure []gm.Point, segment gm.LineSegment) (intersections []Intersection, ok bool) {
	intersections = intersectLineWithFigure(figure, segment)

	if len(intersections) < 2 {
		ok = false
		return
	}

	if !checkAnyPointBelongToSegment(segment, intersections) {
		ok = false
		return
	}

	ok = true
	return
}

func intersectLineWithFigure(figure []gm.Point, segment gm.LineSegment) (intersections []Intersection) {

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

	return
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

func cutFigure(center gm.Point, figure []gm.Point, max_radius float64, segment gm.LineSegment) []gm.Point {
	// Находим перпендикуляр из центра на отрезок
	np := segment.GetNormalPoint(center)
	n := np.Sub(center)

	// Если больше расстояния до самого дальнего угла, игнорим
	normal_length := n.Length()
	if normal_length > max_radius {
		return figure
	}

	intersections, ok := intersectExtendedSegmentWithFigure(figure, segment)
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

	// Находим проекции каждой точки на перпендикуляр. Если отрицательные, выкидываем

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

func checkAnyPointBelongToSegment(segment gm.LineSegment, intersections []Intersection) bool {
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
	return ok
}
