package city_map

import (
	"chirrwick.com/projects/city/generator/genmath"
)

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
	Center  genmath.Point
	Points  []genmath.Point
	Angle   float64
	Streets []genmath.LineSegment
}
