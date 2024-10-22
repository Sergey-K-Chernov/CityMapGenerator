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

type InitialValuesAreas struct {
	NumIndustrial  int
	AreaIndustrial float64
	NumParks       int
	AreaParks      float64
}

type InitialValuesBlocks struct {
	Size Range
}

type Map struct {
	BorderPoints []genmath.Point
	Center       genmath.Point
	Roads        []Road
	Areas        []Area
	Blocks       []Block
}

type Road struct {
	Points []genmath.Point
}

type AreaType int

const (
	AreaIndustrial AreaType = iota
	AreaPark
)

type Area struct {
	Points []genmath.Point
	Type   AreaType
	Area   float64
}

type Block struct {
	Center genmath.Point
	Points []genmath.Point
}
