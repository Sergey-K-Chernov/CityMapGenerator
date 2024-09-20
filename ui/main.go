package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"

	"chirrwick.com/projects/city/generator"
	"chirrwick.com/projects/city/generator/genmath"
)

func main() {
	chan_map := make(chan generator.Map)

	initials := generator.InitialValues{
		Raduis:      generator.Range{Min: 2000., Max: 3000.},
		NumSides:    6,
		VertexShift: 300.0}

	go generator.GenerateBorders(chan_map, initials)

	go startGui(chan_map)

	app.Main()
}

func startGui(chan_map chan generator.Map) {
	window := new(app.Window)
	err := run(window, chan_map)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

var points []genmath.Point

func run(window *app.Window, chan_map chan generator.Map) error {
	theme := material.NewTheme()
	points = make([]genmath.Point, 0)
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)
			mx := gtx.Constraints.Max
			drawTitle(gtx, theme)
			drawMap(&ops, mx, chan_map)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func drawMap(ops *op.Ops, mx image.Point, chan_map chan generator.Map) {
	if len(points) == 0 {
		points = append(points, genmath.Point{X: 100.0, Y: 100.0})
		points = append(points, genmath.Point{X: 100.0, Y: 200.0})
		points = append(points, genmath.Point{X: 200.0, Y: 200.0})
		points = append(points, genmath.Point{X: 200.0, Y: 100.0})
	}

	var cityMap generator.Map
	select {
	case cityMap = <-chan_map:
		fmt.Print("Got map\n")
		points = cityMap.BorderPoints
	default:
		fmt.Print("No map\n")
	}

	var max_map genmath.Point

	for _, p := range points {
		max_map.X = math.Max(max_map.X, p.X)
		max_map.Y = math.Max(max_map.Y, p.Y)
	}
	max_map.X += 100
	max_map.Y += 100

	scale := math.Min(float64(mx.X)/max_map.X, float64(mx.Y)/max_map.Y)

	dark_red := color.NRGBA{R: 0x60, A: 0xFF}
	var path clip.Path
	path.Begin(ops)

	path.MoveTo(f32.Pt(float32(points[0].X*scale), float32(points[0].Y*scale)))
	for _, p := range points {
		path.LineTo(f32.Pt(float32(p.X*scale), float32(p.Y*scale)))
	}
	path.LineTo(f32.Pt(float32(points[len(points)-1].X*scale), float32(points[len(points)-1].Y*scale)))

	path.Close()
	paint.FillShape(ops, dark_red,
		clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op())
}

func drawTitle(gtx layout.Context, theme *material.Theme) {
	// Define an large label with an appropriate text:
	title := material.H6(theme, "Map:")

	// Change the position of the label.
	title.Alignment = text.Start

	// Draw the label to the graphics context.
	title.Layout(gtx)
}
