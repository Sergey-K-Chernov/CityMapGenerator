package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"chirrwick.com/projects/city/city_map"
	"chirrwick.com/projects/city/generator/genmath"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type DrawSettings struct {
	borders   bool
	greenfill bool
	center    bool
}

func tryDrawMap(ops *op.Ops, gtx GC, data *mapData, settings DrawSettings) {
	select {
	case data.cityMap = <-data.channel:
		fmt.Print("Got new map\n")
	default:
		break
	}

	if len(data.cityMap.BorderPoints) < 3 {
		return
	}

	scale := calcScale(gtx, data.cityMap.BorderPoints)

	drawBorders(ops, data.cityMap.BorderPoints, scale, settings)
	drawAreas(ops, data, scale)
	drawBlocks(ops, data, scale)
	drawRoads(ops, data, scale)
	if settings.center {
		drawCenter(ops, data, scale)
	}

}

func calcScale(gtx GC, points []genmath.Point) float64 {
	var max_map genmath.Point

	for _, p := range points {
		max_map.X = math.Max(max_map.X, p.X)
		max_map.Y = math.Max(max_map.Y, p.Y)
	}
	max_map.X += 100
	max_map.Y += 100

	mapConstraints := gtx.Constraints.Max
	mapConstraints.X -= UI_WIDTH

	return math.Min(float64(mapConstraints.X)/max_map.X, float64(mapConstraints.Y)/max_map.Y)
}

func drawBorders(ops *op.Ops, points []genmath.Point, scale float64, settings DrawSettings) {
	if settings.greenfill {
		path := preparePath(ops, points, scale)
		light_green := color.NRGBA{R: 0xDD, G: 0xFF, B: 0xDD, A: 0xFF}
		area := clip.Outline{Path: path.End()}.Op()
		paint.FillShape(ops, light_green, area)
	}

	if settings.borders {
		path := preparePath(ops, points, scale)
		dark_red := color.NRGBA{R: 0x60, A: 0xFF}
		paint.FillShape(ops, dark_red,
			clip.Stroke{
				Path:  path.End(),
				Width: 2,
			}.Op())
	}
}

func drawRoads(ops *op.Ops, data *mapData, scale float64) {
	if len(data.cityMap.Roads) <= 0 {
		return
	}

	for _, rd := range data.cityMap.Roads {
		path := preparePath(ops, rd.Points, scale)

		black := color.NRGBA{A: 0xFF}
		paint.FillShape(ops, black,
			clip.Stroke{
				Path:  path.End(),
				Width: 2,
			}.Op())

	}
}

func drawAreas(ops *op.Ops, data *mapData, scale float64) {
	if len(data.cityMap.Areas) <= 0 {
		return
	}

	for _, area := range data.cityMap.Areas {
		path := preparePath(ops, area.Points, scale)

		grey := color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF}
		light_green := color.NRGBA{R: 0xDD, G: 0xFF, B: 0xDD, A: 0xFF}

		if area.Type == city_map.AreaIndustrial {
			paint.FillShape(ops, grey, clip.Outline{Path: path.End()}.Op())
		}
		if area.Type == city_map.AreaPark {
			paint.FillShape(ops, light_green, clip.Outline{Path: path.End()}.Op())
		}
	}

}

func drawBlocks(ops *op.Ops, data *mapData, scale float64) {
	if len(data.cityMap.Blocks) <= 0 {
		return
	}

	for _, block := range data.cityMap.Blocks {
		path := preparePath(ops, block.Points, scale)
		pale_orange := color.NRGBA{R: 0xFF, G: 0xF6, B: 0xDD, A: 0xFF}
		paint.FillShape(ops, pale_orange, clip.Outline{Path: path.End()}.Op())

		path = preparePath(ops, block.Points, scale)
		black := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
		paint.FillShape(ops, black,
			clip.Stroke{
				Path:  path.End(),
				Width: 1.5,
			}.Op())

		drawStreets(ops, block.Streets, scale)
	}
}

func drawStreets(ops *op.Ops, streets []genmath.LineSegment, scale float64) {
	points := make([]genmath.Point, 2)

	for _, str := range streets {
		points[0] = str.Begin
		points[1] = str.End

		path := preparePath(ops, points, scale)
		black := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
		paint.FillShape(ops, black,
			clip.Stroke{
				Path:  path.End(),
				Width: 1,
			}.Op())
	}

}

func drawDebugVertex(ops *op.Ops, p genmath.Point, clr uint8, scale float64) {
	x := int(p.X*scale + float64(UI_WIDTH))
	y := int(p.Y * scale)
	r := 3

	defer clip.Ellipse{Min: image.Point{X: x - r, Y: y - r},
		Max: image.Point{X: x + r, Y: y + r}}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{B: clr, A: 0xFF}}.Add(ops)

	paint.PaintOp{}.Add(ops)
}

func drawCenter(ops *op.Ops, data *mapData, scale float64) {
	centerX := int(data.cityMap.Center.X*scale + float64(UI_WIDTH))
	centerY := int(data.cityMap.Center.Y * scale)
	r := 2

	defer clip.Ellipse{Min: image.Point{X: centerX - r, Y: centerY - r},
		Max: image.Point{X: centerX + r, Y: centerY + r}}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 0x60, A: 0xFF}}.Add(ops)

	paint.PaintOp{}.Add(ops)
}

func preparePath(ops *op.Ops, points []genmath.Point, scale float64) (path clip.Path) {
	drawPoints := make([]f32.Point, 0)
	for i := 0; i < len(points); i++ {
		x := float32(points[i].X*scale + float64(UI_WIDTH))
		y := float32(points[i].Y * scale)
		drawPoints = append(drawPoints, f32.Pt(x, y))
	}

	path.Begin(ops)

	path.MoveTo(drawPoints[0])
	for _, p := range drawPoints {
		path.LineTo(p)
	}

	path.Close()
	return
}
