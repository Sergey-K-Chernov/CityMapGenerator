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

	points := data.cityMap.BorderPoints
	var max_map genmath.Point

	for _, p := range points {
		max_map.X = math.Max(max_map.X, p.X)
		max_map.Y = math.Max(max_map.Y, p.Y)
	}
	max_map.X += 100
	max_map.Y += 100

	mapConstraints := gtx.Constraints.Max
	mapConstraints.X -= UI_WIDTH

	scale := math.Min(float64(mapConstraints.X)/max_map.X, float64(mapConstraints.Y)/max_map.Y)

	drawPoints := make([]f32.Point, 0)
	for i := 0; i < len(points); i++ {
		x := float32(points[i].X*scale + float64(UI_WIDTH))
		y := float32(points[i].Y * scale)
		drawPoints = append(drawPoints, f32.Pt(x, y))
	}

	dark_red := color.NRGBA{R: 0x60, A: 0xFF}
	var path clip.Path
	path.Begin(ops)

	path.MoveTo(drawPoints[0])
	for _, p := range drawPoints {
		path.LineTo(p)
	}

	path.Close()
	paint.FillShape(ops, dark_red,
		clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op())

	if len(data.cityMap.Roads) > 0 {
		for _, rd := range data.cityMap.Roads {
			dark_blue := color.NRGBA{B: 0x60, A: 0xFF}
			var path clip.Path

			rdPoints := make([]f32.Point, 0)

			for i := 0; i < len(rd.Points); i++ {
				x := float32(rd.Points[i].X*scale + float64(UI_WIDTH))
				y := float32(rd.Points[i].Y * scale)
				rdPoints = append(rdPoints, f32.Pt(x, y))
			}

			path.Begin(ops)

			path.MoveTo(rdPoints[0])
			for _, p := range rdPoints {
				path.LineTo(p)
			}
			path.Close()

			paint.FillShape(ops, dark_blue,
				clip.Stroke{
					Path:  path.End(),
					Width: 2,
				}.Op())

		}
	}

	centerX := int(data.cityMap.Center.X*scale + float64(UI_WIDTH))
	centerY := int(data.cityMap.Center.Y * scale)
	r := 2

	defer clip.Ellipse{Min: image.Point{X: centerX - r, Y: centerY - r},
		Max: image.Point{X: centerX + r, Y: centerY + r}}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 0x60, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
