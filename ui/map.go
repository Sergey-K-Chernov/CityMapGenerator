package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"chirrwick.com/projects/city/generator"
	"chirrwick.com/projects/city/generator/genmath"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func tryDrawMap(ops *op.Ops, mapConstraints image.Point, cityMap *generator.Map, chan_map chan generator.Map) {
	select {
	case *cityMap = <-chan_map:
		fmt.Print("Got new map\n")
	default:
		break
	}

	if len(cityMap.BorderPoints) < 3 {
		return
	}

	points := cityMap.BorderPoints
	var max_map genmath.Point

	for _, p := range points {
		max_map.X = math.Max(max_map.X, p.X)
		max_map.Y = math.Max(max_map.Y, p.Y)
	}
	max_map.X += 100
	max_map.Y += 100

	scale := math.Min(float64(mapConstraints.X)/max_map.X, float64(mapConstraints.Y)/max_map.Y)

	dark_red := color.NRGBA{R: 0x60, A: 0xFF}
	var path clip.Path
	path.Begin(ops)

	path.MoveTo(f32.Pt(float32(points[0].X*scale+float64(UI_WIDTH)), float32(points[0].Y*scale)))
	for _, p := range points {
		path.LineTo(f32.Pt(float32(p.X*scale+float64(UI_WIDTH)), float32(p.Y*scale)))
	}
	path.LineTo(f32.Pt(float32(points[len(points)-1].X*scale+float64(UI_WIDTH)), float32(points[len(points)-1].Y*scale)))

	path.Close()
	paint.FillShape(ops, dark_red,
		clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op())

	defer clip.Ellipse{Min: image.Point{X: int(cityMap.Center.X*scale + float64(UI_WIDTH) - 2), Y: int(cityMap.Center.Y*scale - 2)},
		Max: image.Point{X: int(cityMap.Center.X*scale + float64(UI_WIDTH) + 2), Y: int(cityMap.Center.Y*scale + 2)}}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 0x60, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
