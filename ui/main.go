package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget/material"

	"chirrwick.com/projects/city/generator"
)

type mapData struct {
	cityMap     generator.Map
	channel     chan generator.Map
	invalidator func()
}

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
	ui := makeUi()
	settings := DrawSettings{greenfill: true, borders: true}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			// This graphics context is used for managing the rendering state
			gtx := app.NewContext(&ops, e)

			btnProcessor := ui.pages[ui.currentPage].(uiButtonProcessor)
			btnProcessor.ProcessButtons(gtx, &ui, &data)

			ui.pages[ui.currentPage].Layout(gtx, theme)

			if ui.currentPage <= genBigAreasPage {
				settings.greenfill = false
				settings.borders = true
				settings.center = true
			} else {
				settings.greenfill = true
				settings.borders = false
				settings.center = false
			}

			tryDrawMap(&ops, gtx, &data, settings)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
