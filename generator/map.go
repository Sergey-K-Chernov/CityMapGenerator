package generator

import (
	"chirrwick.com/projects/city/generator/genmath"
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

type Map struct {
	BorderPoints []genmath.Point
	Center       genmath.Point
	Roads        []Road
}

type Road struct {
	Points []genmath.Point
}
