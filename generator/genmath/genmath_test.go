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
