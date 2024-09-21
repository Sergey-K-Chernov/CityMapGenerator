package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"chirrwick.com/projects/city/generator"
	"chirrwick.com/projects/city/generator/genmath"
)

type GC = layout.Context
type Dims = layout.Dimensions

func main() {
	chan_map := make(chan generator.Map)

	go startGui(chan_map)

	app.Main()
}

func startGui(chan_map chan generator.Map) {
	window := new(app.Window)
	window.Option(app.Title("Random city"))
	window.Option(app.MaxSize(800, 600))
	window.Option(app.MinSize(800, 600))
	err := run(window, chan_map)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

var uiShift int = 200

type uiLayout struct {
	minRadius widget.Editor
	maxRadius widget.Editor
	nPoints   widget.Editor
	//	pointVariation widget.Editor

	btnGenerate widget.Clickable
}

func run(window *app.Window, chan_map chan generator.Map) error {
	theme := material.NewTheme()

	var ops op.Ops

	lay := initWidgets()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			// This graphics context is used for managing the rendering state
			gtx := app.NewContext(&ops, e)

			processGenerateButton(gtx, &lay, chan_map)
			layoutUI(gtx, theme, &lay)

			mx := gtx.Constraints.Max
			mx.X -= uiShift
			tryDrawMap(&ops, mx, chan_map)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func initWidgets() (lay uiLayout) {
	lay.nPoints.SingleLine = true
	lay.nPoints.Alignment = text.End

	lay.minRadius.SingleLine = true
	lay.minRadius.Alignment = text.End

	lay.maxRadius.SingleLine = true
	lay.maxRadius.Alignment = text.End

	return lay
}

func layoutUI(gtx GC, theme *material.Theme, lay *uiLayout) {
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		makeFlexLabel(theme, "Corners"),
		makeFlexInput(gtx, theme, &lay.nPoints, "6"),

		makeFlexLabel(theme, "Radius min"),
		makeFlexInput(gtx, theme, &lay.minRadius, "2000"),

		makeFlexLabel(theme, "Radius max"),
		makeFlexInput(gtx, theme, &lay.maxRadius, "3000"),

		makeFlexButton(gtx, theme, lay),
	)
}

func makeFlexLabel(theme *material.Theme, label string) layout.FlexChild {
	return layout.Rigid(func(gtx GC) Dims {
		title := material.H6(theme, label)
		title.Alignment = text.Start

		return title.Layout(gtx)
	})
}

func makeFlexInput(gtx GC, theme *material.Theme, field *widget.Editor, defaultValue string) layout.FlexChild {
	return layout.Rigid(func(gtx GC) Dims {
		ed := material.Editor(theme, field, defaultValue)

		margins := layout.Inset{
			Top:    unit.Dp(4),
			Right:  unit.Dp(705),
			Bottom: unit.Dp(10),
			Left:   unit.Dp(35),
		}

		border := widget.Border{
			Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
			CornerRadius: unit.Dp(3),
			Width:        unit.Dp(2),
		}

		return margins.Layout(gtx,
			func(gtx GC) Dims {
				return border.Layout(gtx, ed.Layout)
			},
		)
	},
	)
}

func makeFlexButton(gtx GC, theme *material.Theme, lay *uiLayout) layout.FlexChild {
	return layout.Rigid(func(gtx GC) Dims {
		margins := layout.Inset{
			Top:    unit.Dp(10),
			Bottom: unit.Dp(10),
			Right:  unit.Dp(605),
			Left:   unit.Dp(35),
		}

		return margins.Layout(gtx,
			func(gtx GC) Dims {
				btn := material.Button(theme, &lay.btnGenerate, "Generate")
				return btn.Layout(gtx)
			},
		)
	},
	)
}

func processGenerateButton(gtx GC, lay *uiLayout, chan_map chan generator.Map) {
	if !lay.btnGenerate.Clicked(gtx) {
		return
	}

	var initials generator.InitialValues

	inputString := lay.nPoints.Text()
	inputString = strings.TrimSpace(inputString)
	nSides, _ := strconv.ParseInt(inputString, 10, 32)
	if nSides == 0 {
		nSides = 6
	}
	initials.NumSides = int(nSides)

	inputString = lay.minRadius.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Min, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Min <= 0 {
		initials.Raduis.Min = 2000.
	}

	inputString = lay.maxRadius.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Max, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Max <= 0 {
		initials.Raduis.Max = 3000.
	}

	go generator.GenerateBorders(chan_map, initials)
}

var cityMap generator.Map
var haveMap bool = false

func tryDrawMap(ops *op.Ops, mx image.Point, chan_map chan generator.Map) {
	select {
	case cityMap = <-chan_map:
		haveMap = true
		fmt.Print("Got new map\n")
		drawMap(ops, mx, cityMap)
	default:
		if haveMap {
			drawMap(ops, mx, cityMap)
		}
		return
	}
}

func drawMap(ops *op.Ops, mx image.Point, cityMap generator.Map) {
	points := cityMap.BorderPoints
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

	path.MoveTo(f32.Pt(float32(points[0].X*scale+float64(uiShift)), float32(points[0].Y*scale)))
	for _, p := range points {
		path.LineTo(f32.Pt(float32(p.X*scale+float64(uiShift)), float32(p.Y*scale)))
	}
	path.LineTo(f32.Pt(float32(points[len(points)-1].X*scale+float64(uiShift)), float32(points[len(points)-1].Y*scale)))

	path.Close()
	paint.FillShape(ops, dark_red,
		clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op())
}
