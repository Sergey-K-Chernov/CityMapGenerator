package generator

import (
	"testing"

	"chirrwick.com/projects/city/generator/genmath"
	gm "chirrwick.com/projects/city/generator/genmath"
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
				if checkXRayIntersectsSegment(genmath.Point{X: 0, Y: 0}.Add(addition), a[i].Add(addition), b[i].Add(addition)) != answers[i] {
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

func TestCutFigure(t *testing.T) {
	center := gm.Point{X: 25, Y: 99}
	figure := []gm.Point{
		{X: 11, Y: 82},
		{X: 3, Y: 104},
		{X: 35, Y: 118},
		{X: 53, Y: 91},
	}

	segments := []gm.LineSegment{
		{Begin: gm.Point{X: 7, Y: 147}, End: gm.Point{X: 19, Y: 115}},
		{Begin: gm.Point{X: 32, Y: 85}, End: gm.Point{X: 45, Y: 51}},
		{Begin: gm.Point{X: 16, Y: 142}, End: gm.Point{X: 32, Y: 100}},
		{Begin: gm.Point{X: 32, Y: 100}, End: gm.Point{X: 16, Y: 142}},
		{Begin: gm.Point{X: 26, Y: 130}, End: gm.Point{X: 49, Y: 69}},
		{Begin: gm.Point{X: 44, Y: 96}, End: gm.Point{X: 52, Y: 75}},
		{Begin: gm.Point{X: 7, Y: 115}, End: gm.Point{X: 19, Y: 82}},
	}

	answers := [][]gm.Point{
		{
			{X: 11, Y: 82},
			{X: 3, Y: 104},
			{X: 35, Y: 118},
			{X: 53, Y: 91},
		},
		{
			{X: 11, Y: 82},
			{X: 3, Y: 104},
			{X: 35, Y: 118},
			{X: 53, Y: 91},
		},
		{
			{X: 36.754717, Y: 87.518868},
			{X: 11, Y: 82},
			{X: 3, Y: 104},
			{X: 26.55102, Y: 114.303571},
		},
		{
			{X: 36.754717, Y: 87.518868},
			{X: 11, Y: 82},
			{X: 3, Y: 104},
			{X: 26.55102, Y: 114.303571},
		},
		{
			{X: 41.624052, Y: 88.562297},
			{X: 11, Y: 82},
			{X: 3, Y: 104},
			{X: 31.158311, Y: 116.319261},
		},
		{
			{X: 46.440252, Y: 89.59434},
			{X: 11, Y: 82},
			{X: 3, Y: 104},
			{X: 35, Y: 118},
			{X: 36.444444, Y: 115.833333},
		},
		{
			{X: 18.42168, Y: 83.590361},
			{X: 9.901961, Y: 107.019608},
			{X: 35, Y: 118},
			{X: 53, Y: 91},
		},
	}

	if len(segments) != len(answers) {
		t.Fatalf("Wrong conditions, array lengths are not equal")
	}

	for i, segment := range segments {
		fig := make([]gm.Point, 4)
		copy(fig, figure)
		fig = cutFigure(center, fig, 29.12044, segment)
		for j, p := range fig {
			if !gm.AlmostEqualPoints(p, answers[i][j], 0.0001) {
				t.Fatalf("Test %d falied", i)
			}
		}
	}

	figure = []gm.Point{
		{X: 1544.199, Y: 4116.501},
		{X: 1118.902, Y: 3888.183},
		{X: 1026.936, Y: 4059.493},
		{X: 1452.233, Y: 4287.811},
	}

	center = gm.Point{X: 1285.568, Y: 4087.997}
	segment := gm.LineSegment{Begin: gm.Point{X: 1293.684, Y: 3746.878}, End: gm.Point{X: 1750.855, Y: 4933.706}}

	answer := []gm.Point{
		{X: 1481.239949, Y: 4233.778201},
		{X: 1407.8728, Y: 4043.315143},
		{X: 1118.902, Y: 3888.183},
		{X: 1026.936, Y: 4059.493},
		{X: 1452.233, Y: 4287.811},
	}

	figure = cutFigure(center, figure, 260.197484, segment)
	for j, p := range figure {
		if !gm.AlmostEqualPoints(p, answer[j], 0.0001) {
			t.Fatalf("Additional test 1 failed")
		}
	}

}

/*
func TestCropBlockByBlocks(t *testing.T) {
	block1 := []gm.Point{
		{X: 2166.407, Y: 1666.528},
		{X: 2203.204, Y: 1358.976},
		{X: 2100.28, Y: 1346.662},
		{X: 2063.483, Y: 1654.214},
	}

	center1 := gm.Point{X: 2133.343, Y: 1506.595}

	block2 := []gm.Point{
		{X: 2112.96, Y: 1096.496},
		{X: 1803.338, Y: 1426.942},
		{X: 1988.472, Y: 1600.409},
		{X: 2298.094, Y: 1269.964},
	}
	center2 := gm.Point{X: 2050.716, Y: 1348.453}

	answer := []gm.Point{
		{X: 2112.96, Y: 1096.496},
		{X: 1803.338, Y: 1426.942},
		{X: 1988.472, Y: 1600.409},

		{X: 2081.8815, Y: 1500.438},
		{X: 2100.28, Y: 1346.662},
		{X: 2203.204, Y: 1358.976},
		{X: 2201.522483, Y: 1373.030238},

		{X: 2298.094, Y: 1269.964},
	}

	println(block1)
	println(center1.X, center1.Y)
	println(block2)
	println(center2.X, center2.Y)
	println(answer)

	t.Fatalf("Test unfinished")
}
*/
