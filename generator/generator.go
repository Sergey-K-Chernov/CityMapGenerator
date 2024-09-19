package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"chirrwick.com/projects/city/generator/genmath"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func generateRadialRandomPoint(angle_min, angle_max, radius_min, raduis_max float64) genmath.Point {
	angle := genmath.RandFloat(angle_min, angle_max)
	radius := genmath.RandFloat(radius_min, raduis_max)

	var res genmath.Point
	res.X = radius * math.Cos(angle)
	res.Y = radius * math.Sin(angle)
	return res
}

func shiftPoints(points []genmath.Point) (shift genmath.Point) {
	for _, point := range points {
		shift.X = math.Min(shift.X, point.X)
		shift.Y = math.Min(shift.Y, point.Y)
	}
	shift.X = -shift.X + 100
	shift.Y = -shift.Y + 100

	for i := range points {
		points[i].X += shift.X
		points[i].Y += shift.Y
	}

	return
}

func GenerateBorders(chan_map chan []genmath.Point, n_points int, size_rough_min, size_rough_max float64) {
	angle_step := genmath.DegToRad(360. / float64(n_points))
	angle_variation := angle_step / 2

	points := make([]genmath.Point, n_points)

	for i := 0; i < n_points; i++ {
		angle := angle_step * float64(i)
		points[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, size_rough_min/2, size_rough_max/2)

		point := generateRadialRandomPoint(0, 2*math.Pi, 0, (size_rough_min+size_rough_max)/20.0)
		points[i].Add(point)
	}

	shift := shiftPoints(points)

	for _, point := range points {
		fmt.Printf("%7.1f\t%7.1f\n", point.X, point.Y)
	}

	fmt.Printf("\nShift: %7.1f\t%7.1f\n", shift.X, shift.Y)

	chan_map <- points

}

func main() {
	chan_map := make(chan []genmath.Point)

	go GenerateBorders(chan_map, 6, 2000., 3000.)
	go startGui(chan_map)

	app.Main()
}

func startGui(chan_map chan []genmath.Point) {
	window := new(app.Window)
	err := run(window, chan_map)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

var points []genmath.Point

func run(window *app.Window, chan_map chan []genmath.Point) error {
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

func drawMap(ops *op.Ops, mx image.Point, chan_map chan []genmath.Point) {
	if len(points) == 0 {
		points = append(points, genmath.Point{X: 100.0, Y: 100.0})
		points = append(points, genmath.Point{X: 100.0, Y: 200.0})
		points = append(points, genmath.Point{X: 200.0, Y: 200.0})
		points = append(points, genmath.Point{X: 200.0, Y: 100.0})
	}

	select {
	case points = <-chan_map:
		fmt.Print("Got map\n")
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

	dark_red := color.NRGBA{R: 0x50, A: 0xFF}
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
