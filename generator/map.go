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

type InitialValuesBlocks struct {
	Size Range
}

type Map struct {
	BorderPoints []genmath.Point
	Center       genmath.Point
	Roads        []Road
	Blocks       []Block
}

type Road struct {
	Points []genmath.Point
}

type Block struct {
	Points []genmath.Point
}
