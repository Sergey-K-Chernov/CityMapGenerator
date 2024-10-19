package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"chirrwick.com/projects/city/generator/genmath"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func tryDrawMap(ops *op.Ops, gtx GC, data *mapData) {
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

	drawBorders(ops, data.cityMap.BorderPoints, scale)
	drawRoads(ops, data, scale)
	drawBlocks(ops, data, scale)
	drawCenter(ops, data, scale)

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

func drawBorders(ops *op.Ops, points []genmath.Point, scale float64) {
	path := preparePath(ops, points, scale)

	dark_red := color.NRGBA{R: 0x60, A: 0xFF}
	paint.FillShape(ops, dark_red,
		clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op())
}

func drawRoads(ops *op.Ops, data *mapData, scale float64) {
	if len(data.cityMap.Roads) <= 0 {
		return
	}

	for _, rd := range data.cityMap.Roads {
		path := preparePath(ops, rd.Points, scale)

		dark_blue := color.NRGBA{B: 0x60, A: 0xFF}
		paint.FillShape(ops, dark_blue,
			clip.Stroke{
				Path:  path.End(),
				Width: 2,
			}.Op())

	}
}

func drawBlocks(ops *op.Ops, data *mapData, scale float64) {
	if len(data.cityMap.Blocks) <= 0 {
		return
	}

	for _, block := range data.cityMap.Blocks {
		path := preparePath(ops, block.Points, scale)

		dark_blue := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
		paint.FillShape(ops, dark_blue,
			clip.Stroke{
				Path:  path.End(),
				Width: 2,
			}.Op())

		for i, p := range block.Points {
			clr := uint8(float64(i) / float64(len(block.Points)) * 255)
			drawDebugVertex(ops, p, clr, scale)
		}
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
