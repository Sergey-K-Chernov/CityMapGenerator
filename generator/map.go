package generator

import (
	"chirrwick.com/projects/city/generator/genmath"
)

type Range struct {
	Min float64
	Max float64
}

type InitialValues struct {
	Raduis      Range
	NumSides    int
	VertexShift float64
}

type Map struct {
	BorderPoints []genmath.Point
	Center       genmath.Point
}
