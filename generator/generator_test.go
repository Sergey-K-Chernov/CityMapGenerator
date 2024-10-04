package generator

import (
	"testing"

	"chirrwick.com/projects/city/generator/genmath"
)

func TestCheckXRayIntersectsSection(t *testing.T) {
	a := [...]genmath.Point{
		{X: 10, Y: 10},
		{X: -10, Y: 10},
		{X: 0, Y: 10},
		{X: -10, Y: 10},
		{X: -10, Y: -10},
		{X: -10, Y: 0},
		{X: -10, Y: 10},
		{X: -10, Y: -10},
		{X: -11, Y: -10},
		{X: -9, Y: -10},
		{X: -11, Y: 10},
		{X: -9, Y: 10},
	}

	b := [...]genmath.Point{
		{X: 10, Y: -10},
		{X: -10, Y: -10},
		{X: 0, Y: -10},
		{X: 10, Y: 10},
		{X: 10, Y: -10},
		{X: 10, Y: 0},
		{X: 10, Y: -10},
		{X: 10, Y: 10},
		{X: 9, Y: 10},
		{X: 11, Y: 10},
		{X: 9, Y: -10},
		{X: 11, Y: -10},
	}

	answers := [...]int{
		1,
		0,
		1,
		0,
		0,
		1,
		1,
		1,
		0,
		1,
		0,
		1,
	}

	if len(a) != len(b) || len(a) != len(answers) {
		t.Fatalf("Wrong conditions, array lengths are not equal")
	}

	for i := 0; i < len(a); i++ {
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				addition := genmath.Point{X: float64(x), Y: float64(y)}
				if checkXRayIntersectsSection(genmath.Point{X: 0, Y: 0}.Add(addition), a[i].Add(addition), b[i].Add(addition)) != answers[i] {
					t.Fatalf("Test %d failed for x=%d, y=%d.", i, x, y)
				}
			}
		}

	}
}

func TestInsideTriangle(t *testing.T) {
	triangle := [...]genmath.Point{
		{X: 2, Y: 2},
		{X: 7, Y: 12},
		{X: 14, Y: 7},
	}

	points := [...]genmath.Point{
		{X: 7, Y: 6},
		{X: 1, Y: 9},
		{X: 8, Y: 2},
		{X: 8, Y: 2},
	}

	answers := []bool{
		true,
		false,
		false,
		false,
	}

	for i, p := range points {
		if checkPointInsidePolygon(p, triangle[:]) != answers[i] {
			t.Fatalf("Test %d falied", i)
		}
	}
}

func TestInsideRectangle(t *testing.T) {
	rectangle := [...]genmath.Point{
		{X: 5, Y: 6},
		{X: 5, Y: -6},
		{X: -5, Y: -6},
		{X: -5, Y: 6},
	}

	points := [...]genmath.Point{
		{X: 0, Y: 0},
		{X: 9, Y: -3},
		{X: -3, Y: -12},
		{X: -8, Y: 1},
		{X: 1, Y: 12},
	}

	answers := []bool{
		true,
		false,
		false,
		false,
		false,
	}

	for i, p := range points {
		if checkPointInsidePolygon(p, rectangle[:]) != answers[i] {
			t.Fatalf("Test %d falied", i)
		}
	}
}

func TestInsideM(t *testing.T) {
	m := [...]genmath.Point{
		{X: 6, Y: 9},
		{X: 10, Y: 23},
		{X: 15, Y: 12},
		{X: 19, Y: 23},
		{X: 26, Y: 1},
		{X: 19, Y: 14},
		{X: 13, Y: 9},
		{X: 11, Y: 18},
	}

	points := [...]genmath.Point{
		{X: 6, Y: 16},
		{X: 11, Y: 26},
		{X: 15, Y: 17},
		{X: 19, Y: 26},
		{X: 24, Y: 16},
		{X: 19, Y: 11},
		{X: 10, Y: 13},

		{X: 8, Y: 13},
		{X: 10, Y: 19},
		{X: 14, Y: 12},
		{X: 19, Y: 18},
		{X: 23, Y: 8},
	}

	answers := []bool{
		false,
		false,
		false,
		false,
		false,
		false,
		false,

		true,
		true,
		true,
		true,
		true,
	}

	for i, p := range points {
		if checkPointInsidePolygon(p, m[:]) != answers[i] {
			t.Fatalf("Test %d falied", i)
		}
	}
}
