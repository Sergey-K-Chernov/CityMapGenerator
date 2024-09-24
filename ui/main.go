package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"

	"chirrwick.com/projects/city/generator"
)

type mapData struct {
	cityMap     generator.Map
	channel     chan generator.Map
	invalidator func()
}

type GC = layout.Context
type Dims = layout.Dimensions

const UI_WIDTH int = 150

func main() {
	go startGui()

	app.Main()
}

func startGui() {
	window := new(app.Window)
	window.Option(app.Title("Random city"))
	err := run(window)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run(window *app.Window) error {
	theme := material.NewTheme()

	var data mapData
	data.channel = make(chan generator.Map)
	data.invalidator = func() {
		window.Invalidate()
	}

	var ops op.Ops

	bordersPage := createBordersPage()
	centersRoadsPage := createCentersAndRoadsPage()

	var ui uiPages
	ui.pages = make([]uiLayouter, 2)
	ui.pages[genBordersPage] = &bordersPage
	ui.pages[genCentersAndRoadsPage] = &centersRoadsPage
	ui.currentPage = genBordersPage

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			// This graphics context is used for managing the rendering state
			gtx := app.NewContext(&ops, e)

			//processGenerateButton(gtx, &bordersPage, data.channel, )

			//ui.currentPage = processAcceptButton(gtx, &bordersPage, data.cityMap)

			btnProcessor := ui.pages[ui.currentPage].(uiButtonProcessor)

			btnProcessor.ProcessButtons(gtx, &ui, &data)
			ui.pages[ui.currentPage].Layout(gtx, theme)

			mapConstraints := gtx.Constraints.Max
			mapConstraints.X -= UI_WIDTH
			tryDrawMap(&ops, mapConstraints, &data.cityMap, data.channel)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
