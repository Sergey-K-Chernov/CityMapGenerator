package genmath

import (
	"math"
	"testing"
)

func TestAlmostEqual(t *testing.T) {
	vals1 := [10]float64{0.0,
		0.0,
		0.1,
		0.1,
		-0.1,
		-0.1,
		22.456789,
		-33.987654321,
		2222222222.0,
		-123456798.0}

	vals2 := [10]float64{0.0,
		1.0e-308,
		0.10000001,
		0.100000001,
		-0.1000001,
		-0.10000001,
		22.456788,
		-33.987654322,
		2222222222.0,
		-123456798.0}

	thres := [10]float64{0.0,
		0.1e-307,
		0.00000001,
		0.0000000001,
		0.0000001,
		0.000000001,
		0.00001,
		0.000000005,
		0.1,
		0.1}

	answers := [10]bool{true,
		true,
		true,
		false,
		true,
		false,
		true,
		true,
		true,
		true}

	for i := 0; i < len(answers); i++ {
		answer := AlmostEqual(vals1[i], vals2[i], thres[i])
		if answer != answers[i] {
			t.Fatalf(`%d: %t != %t`, i, answer, answers[i])
		}
	}
}

func TestDegToRad(t *testing.T) {
	vals := [9]float64{0.,
		45.,
		60.,
		90.,
		135.,
		180.,
		225.,
		270.,
		360.}

	answers := [9]float64{0.,
		math.Pi / 4,
		math.Pi / 3,
		math.Pi / 2,
		3 * math.Pi / 4,
		math.Pi,
		5 * math.Pi / 4,
		3 * math.Pi / 2,
		2 * math.Pi}

	for i := 0; i < len(vals); i++ {
		answer := DegToRad(vals[i])
		if !AlmostEqual(answer, answers[i], 0.000001) {
			t.Fatalf(`%d: %f != %f`, i, answer, answers[i])
		}
	}
}

func TestIntersectLineSegments(t *testing.T) {
	l1s := [...]LineSegment{
		{Begin: Point{X: -5, Y: 0}, End: Point{X: 5, Y: 0}},
		{Begin: Point{X: 0, Y: 0}, End: Point{X: 10, Y: 0}},
		{Begin: Point{X: 0, Y: 0}, End: Point{X: 10, Y: 10}},
		{Begin: Point{X: 0, Y: 0}, End: Point{X: 10, Y: 0}},
		{Begin: Point{X: 0, Y: 0}, End: Point{X: 10, Y: 0}},
		{Begin: Point{X: 0, Y: 0}, End: Point{X: 10, Y: 10}},
	}

	l2s := [...]LineSegment{
		{Begin: Point{X: 0, Y: -5}, End: Point{X: 0, Y: 5}},
		{Begin: Point{X: 0, Y: 0}, End: Point{X: 0, Y: 10}},
		{Begin: Point{X: 0, Y: 10}, End: Point{X: 10, Y: 0}},
		{Begin: Point{X: 10, Y: 10}, End: Point{X: 10, Y: 0}},
		{Begin: Point{X: 11, Y: 10}, End: Point{X: 11, Y: 0}},
		{Begin: Point{X: 0, Y: -1}, End: Point{X: 10, Y: 9}},
	}

	type Answer struct {
		p  Point
		ok bool
	}

	answers := [...]Answer{
		{p: Point{X: 0, Y: 0}, ok: true},
		{p: Point{X: 0, Y: 0}, ok: true},
		{p: Point{X: 5, Y: 5}, ok: true},
		{p: Point{X: 10, Y: 0}, ok: true},
		{p: Point{}, ok: false},
		{p: Point{}, ok: false},
	}

	for i, l1 := range l1s {
		l2 := l2s[i]
		p, ok := l1.Intersect(l2)

		if ok != answers[i].ok {
			t.Fatalf(`case %d failed: ok = %t, must be %t`, i, ok, answers[i].ok)
		} else {
			if ok && !AlmostEqual(p.X, answers[i].p.X, 0.000001) && AlmostEqual(p.Y, answers[i].p.Y, 0.000001) {
				t.Fatalf(`case %d failed: point is (%f;%f), must be (%f;%f)`, i, p.X, p.Y, answers[i].p.X, answers[i].p.Y)
			}
		}
	}
}

func TestIntersectLines(t *testing.T) {
	l1s := [...]Line{
		{Origin: Point{X: -5, Y: 0}, Vector: Point{X: 1, Y: 0}},
		{Origin: Point{X: 0, Y: 0}, Vector: Point{X: 1, Y: 0}},
		{Origin: Point{X: 0, Y: 0}, Vector: Point{X: 0.707107, Y: 0.707107}},
		{Origin: Point{X: 0, Y: 0}, Vector: Point{X: 1, Y: 0}},
		{Origin: Point{X: 0, Y: 0}, Vector: Point{X: 1, Y: 0}},
		{Origin: Point{X: 0, Y: 0}, Vector: Point{X: 0.707107, Y: 0.707107}},
	}

	l2s := [...]Line{
		{Origin: Point{X: 0, Y: -5}, Vector: Point{X: 0, Y: 1}},
		{Origin: Point{X: 0, Y: 0}, Vector: Point{X: 0, Y: 1}},
		{Origin: Point{X: 0, Y: 10}, Vector: Point{X: 0.707107, Y: -0.707107}},
		{Origin: Point{X: 10, Y: 10}, Vector: Point{X: 0, Y: -1}},
		{Origin: Point{X: 11, Y: 10}, Vector: Point{X: 0, Y: -1}},
		{Origin: Point{X: 0, Y: -1}, Vector: Point{X: 0.707107, Y: 0.707107}},
	}

	type Answer struct {
		p  Point
		ok bool
	}

	answers := [...]Answer{
		{p: Point{X: 0, Y: 0}, ok: true},
		{p: Point{X: 0, Y: 0}, ok: true},
		{p: Point{X: 5, Y: 5}, ok: true},
		{p: Point{X: 10, Y: 0}, ok: true},
		{p: Point{X: 11, Y: 0}, ok: true},
		{p: Point{}, ok: false},
	}

	for i, l1 := range l1s {
		l2 := l2s[i]
		p, ok := l1.Intersect(l2)

		if ok != answers[i].ok {
			t.Fatalf(`case %d failed: ok = %t, must be %t`, i, ok, answers[i].ok)
		} else {
			if ok && !AlmostEqual(p.X, answers[i].p.X, 0.000001) && AlmostEqual(p.Y, answers[i].p.Y, 0.000001) {
				t.Fatalf(`case %d failed: point is (%f;%f), must be (%f;%f)`, i, p.X, p.Y, answers[i].p.X, answers[i].p.Y)
			}
		}
	}
}
