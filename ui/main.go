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
	go startGui()

	app.Main()
}

func startGui() {
	window := new(app.Window)
	window.Option(app.Title("Random city"))
	window.Option(app.MaxSize(800, 600))
	window.Option(app.MinSize(800, 600))
	err := run(window)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

const UI_SHIFT int = 200

type uiLayout struct {
	minRadius      widget.Editor
	maxRadius      widget.Editor
	nPoints        widget.Editor
	pointVariation widget.Editor

	btnGenerate widget.Clickable
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var cityMap generator.Map
	chan_map := make(chan generator.Map)

	var ops op.Ops

	lay := initWidgets()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			// This graphics context is used for managing the rendering state
			gtx := app.NewContext(&ops, e)

			processGenerateButton(gtx, &lay, chan_map, func() {
				window.Invalidate()
			})
			layoutUI(gtx, theme, &lay)

			mx := gtx.Constraints.Max
			mx.X -= UI_SHIFT
			tryDrawMap(&ops, mx, &cityMap, chan_map)

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

	lay.pointVariation.SingleLine = true
	lay.pointVariation.Alignment = text.End

	return lay
}

func layoutUI(gtx GC, theme *material.Theme, lay *uiLayout) {
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		makeFlexLabel(theme, "Corners"),
		makeFlexInput(gtx, theme, &lay.nPoints, "3"),

		makeFlexLabel(theme, "Radius min"),
		makeFlexInput(gtx, theme, &lay.minRadius, "2000"),

		makeFlexLabel(theme, "Radius max"),
		makeFlexInput(gtx, theme, &lay.maxRadius, "3000"),

		makeFlexLabel(theme, "Variation"),
		makeFlexInput(gtx, theme, &lay.pointVariation, "300"),

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

func processGenerateButton(gtx GC, lay *uiLayout, chan_map chan generator.Map, callback func()) {
	if !lay.btnGenerate.Clicked(gtx) {
		return
	}

	var initials generator.InitialValues

	inputString := lay.nPoints.Text()
	inputString = strings.TrimSpace(inputString)
	nSides, _ := strconv.ParseInt(inputString, 10, 32)
	if nSides < 3 {
		nSides = 3
	}
	initials.NumSides = int(nSides)

	inputString = lay.minRadius.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Min, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Min <= 0 {
		initials.Raduis.Min = 2000.0
	}

	inputString = lay.maxRadius.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Max, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Max <= 0 {
		initials.Raduis.Max = 3000.0
	}

	inputString = lay.pointVariation.Text()
	inputString = strings.TrimSpace(inputString)
	initials.VertexShift, _ = strconv.ParseFloat(inputString, 32)
	if initials.VertexShift < 0.0 {
		initials.VertexShift = initials.Raduis.Max / 10.0
	}

	go generateBorders(chan_map, initials, callback)
}

func generateBorders(chan_map chan generator.Map, initials generator.InitialValues, callback func()) {
	generator.GenerateBorders(chan_map, initials)
	callback()
}

func tryDrawMap(ops *op.Ops, mx image.Point, cityMap *generator.Map, chan_map chan generator.Map) {
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

	scale := math.Min(float64(mx.X)/max_map.X, float64(mx.Y)/max_map.Y)

	dark_red := color.NRGBA{R: 0x60, A: 0xFF}
	var path clip.Path
	path.Begin(ops)

	path.MoveTo(f32.Pt(float32(points[0].X*scale+float64(UI_SHIFT)), float32(points[0].Y*scale)))
	for _, p := range points {
		path.LineTo(f32.Pt(float32(p.X*scale+float64(UI_SHIFT)), float32(p.Y*scale)))
	}
	path.LineTo(f32.Pt(float32(points[len(points)-1].X*scale+float64(UI_SHIFT)), float32(points[len(points)-1].Y*scale)))

	path.Close()
	paint.FillShape(ops, dark_red,
		clip.Stroke{
			Path:  path.End(),
			Width: 2,
		}.Op())

	defer clip.Ellipse{Min: image.Point{X: int(cityMap.Center.X*scale + float64(UI_SHIFT) - 2), Y: int(cityMap.Center.Y*scale - 2)},
		Max: image.Point{X: int(cityMap.Center.X*scale + float64(UI_SHIFT) + 2), Y: int(cityMap.Center.Y*scale + 2)}}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 0x60, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}
