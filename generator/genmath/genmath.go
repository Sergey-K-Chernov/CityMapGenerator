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

type Vector2D = Point

type Line struct {
	Origin Point
	Vector Vector2D
}

type LineSegment struct {
	Begin Point
	End   Point
}

type Rect struct {
	Left, Top, Right, Bottom float64
}

type Triangle struct {
	A, B, C Point
}

func (p Point) LengthSq() float64 {
	return p.X*p.X + p.Y*p.Y
}

func (p Point) Length() float64 {
	return math.Sqrt(p.LengthSq())
}

func (p Point) Angle() float64 {
	return math.Atan2(p.Y, p.X)
}

func (s LineSegment) GetVector() Vector2D {
	return Vector2D{X: s.End.X - s.Begin.X, Y: s.End.Y - s.Begin.Y}
}

func (s LineSegment) LengthSq() float64 {
	p := s.GetVector()
	return p.X*p.X + p.Y*p.Y
}

func (s LineSegment) Length() float64 {
	return math.Sqrt(s.LengthSq())
}

func (s LineSegment) ToRect() Rect {
	return Rect{
		Left:   math.Min(s.Begin.X, s.End.X),
		Right:  math.Max(s.Begin.X, s.End.X),
		Top:    math.Max(s.Begin.Y, s.End.Y),
		Bottom: math.Min(s.Begin.Y, s.End.Y)}
}

func (s LineSegment) Split(ratio float64) (LineSegment, LineSegment) {
	vec := s.GetVector()
	vec.Scale(ratio)
	p := s.Begin.Add(vec)

	return LineSegment{s.Begin, p}, LineSegment{p, s.End}
}

func (r Rect) HasPoint(p Point) bool {
	return p.X >= r.Left && p.X <= r.Right && p.Y >= r.Bottom && p.Y <= r.Top
}

func (t Triangle) HasPoint(p Point) bool {
	s1 := (t.A.X-p.X)*(t.B.Y-t.A.Y) - (t.B.X-t.A.X)*(t.A.Y-p.Y)
	s2 := (t.B.X-p.X)*(t.C.Y-t.B.Y) - (t.C.X-t.B.X)*(t.B.Y-p.Y)
	s3 := (t.C.X-p.X)*(t.A.Y-t.C.Y) - (t.A.X-t.C.X)*(t.C.Y-p.Y)

	if math.Signbit(s1) == math.Signbit(s2) && math.Signbit(s1) == math.Signbit(s3) {
		return true
	}
	return false
}

func (p *Point) Normalize() *Point {
	length := math.Sqrt(p.X*p.X + p.Y*p.Y)
	p.X /= length
	p.Y /= length
	return p
}

func (p *Point) Scale(factor float64) *Point {
	p.X *= factor
	p.Y *= factor
	return p
}

func (p1 *Point) AddInPlace(p2 Point) {
	p1.X += p2.X
	p1.Y += p2.Y
}

func (p1 Point) Add(p2 Point) Point {
	return Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
}

func (p1 Point) Sub(p2 Point) Point {
	return Point{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}

func (l1 LineSegment) Intersect(l2 LineSegment) (Point, bool) {
	xy12 := l1.Begin.X*l1.End.Y - l1.Begin.Y*l1.End.X
	xy34 := l2.Begin.X*l2.End.Y - l2.Begin.Y*l2.End.X

	x12 := l1.Begin.X - l1.End.X
	x34 := l2.Begin.X - l2.End.X

	y12 := l1.Begin.Y - l1.End.Y
	y34 := l2.Begin.Y - l2.End.Y

	denominator := x12*y34 - y12*x34
	if denominator == 0 {
		return Point{}, false
	}

	x := (xy12*x34 - x12*xy34) / denominator
	y := (xy12*y34 - y12*xy34) / denominator

	p := Point{X: x, Y: y}

	if l1.ToRect().HasPoint(p) && l2.ToRect().HasPoint(p) {
		return p, true
	}

	return Point{}, false
}

func (l1 Line) Intersect(l2 Line) (Point, bool) {
	end1 := l1.Origin.Add(l1.Vector)
	end2 := l2.Origin.Add(l2.Vector)

	xy12 := l1.Origin.X*end1.Y - l1.Origin.Y*end1.X
	xy34 := l2.Origin.X*end2.Y - l2.Origin.Y*end2.X

	x12 := l1.Origin.X - end1.X
	x34 := l2.Origin.X - end2.X

	y12 := l1.Origin.Y - end1.Y
	y34 := l2.Origin.Y - end2.Y

	denominator := x12*y34 - y12*x34
	if denominator == 0 {
		return Point{}, false
	}

	x := (xy12*x34 - x12*xy34) / denominator
	y := (xy12*y34 - y12*xy34) / denominator

	return Point{X: x, Y: y}, true
}
