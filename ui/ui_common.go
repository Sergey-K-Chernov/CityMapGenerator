package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const UI_WIDTH int = 150

type GC = layout.Context
type Dims = layout.Dimensions

type uiPage int

const (
	genBordersPage uiPage = iota
	genCentersAndRoadsPage
	genBlocksPage
	genStreetsPage
)

type inputField struct {
	field        widget.Editor
	label        string
	defaultValue string
}

type button struct {
	button widget.Clickable
	label  string
}

type uiLayouter interface {
	Layout(gtx GC, theme *material.Theme)
}

type uiButtonProcessor interface {
	ProcessButtons(gtx GC, ui *uiPages, data *mapData)
}

type uiPages struct {
	currentPage uiPage
	pages       []uiLayouter
}

func makeUi() uiPages {
	bordersPage := createBordersPage()
	centersRoadsPage := createCentersAndRoadsPage()
	blocksPage := createBlocksPage()

	var ui uiPages
	ui.pages = make([]uiLayouter, 3)
	ui.pages[genBordersPage] = &bordersPage
	ui.pages[genCentersAndRoadsPage] = &centersRoadsPage
	ui.pages[genBlocksPage] = &blocksPage
	ui.currentPage = genBordersPage

	return ui
}

func makeLabel(theme *material.Theme, label string) layout.FlexChild {
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
			Top:    unit.Dp(3),
			Right:  unit.Dp(4),
			Bottom: unit.Dp(6),
			Left:   unit.Dp(4),
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
	})
}

func makeButton(gtx GC, theme *material.Theme, button *widget.Clickable, label string) layout.FlexChild {
	return layout.Rigid(func(gtx GC) Dims {
		margins := layout.Inset{
			Top:    unit.Dp(3),
			Bottom: unit.Dp(6),
			Right:  unit.Dp(4),
			Left:   unit.Dp(4),
		}

		return margins.Layout(gtx,
			func(gtx GC) Dims {
				btn := material.Button(theme, button, label)
				return btn.Layout(gtx)
			})
	})
}
