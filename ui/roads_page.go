package main

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func createCentersAndRoadsPage() (lay uiCentersAndRoadsPage) {
	lay.nCenters.field.SingleLine = true
	lay.nCenters.field.Alignment = text.End
	lay.nCenters.label = "Centers"
	lay.nCenters.defaultValue = "2"

	lay.minRadius.field.SingleLine = true
	lay.minRadius.field.Alignment = text.End
	lay.minRadius.label = "Min Radius"
	lay.minRadius.defaultValue = "500"

	lay.maxRadius.field.SingleLine = true
	lay.maxRadius.field.Alignment = text.End
	lay.maxRadius.label = "Max Radius"
	lay.maxRadius.defaultValue = "1000"

	lay.branching.field.SingleLine = true
	lay.branching.field.Alignment = text.End
	lay.branching.label = "Branching"
	lay.branching.defaultValue = "2"

	lay.btnGenerate.label = "Generate"
	lay.btnAccept.label = "Accept"

	return lay
}

func (l *uiCentersAndRoadsPage) Layout(gtx GC, theme *material.Theme) {
	// Fixed ui hack...
	var uiFlexWeight float32
	totalWidth := gtx.Constraints.Max.X
	mapWidth := totalWidth - UI_WIDTH

	uiFlexWeight = float32(UI_WIDTH) / float32(totalWidth)

	layout.Flex{}.Layout(gtx,
		layout.Flexed(uiFlexWeight, func(gtx GC) Dims {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,

				makeLabel(theme, l.nCenters.label),
				makeFlexInput(gtx, theme, &l.nCenters.field, l.nCenters.defaultValue),

				makeLabel(theme, l.minRadius.label),
				makeFlexInput(gtx, theme, &l.minRadius.field, l.minRadius.defaultValue),

				makeLabel(theme, l.maxRadius.label),
				makeFlexInput(gtx, theme, &l.maxRadius.field, l.maxRadius.defaultValue),

				makeLabel(theme, l.branching.label),
				makeFlexInput(gtx, theme, &l.branching.field, l.branching.defaultValue),

				makeButton(gtx, theme, &l.btnGenerate.button, l.btnGenerate.label),
				makeButton(gtx, theme, &l.btnAccept.button, l.btnAccept.label),
			)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(mapWidth)}.Layout),
	)
}

func (l *uiCentersAndRoadsPage) ProcessButtons(gtx GC, ui *uiPages, data *mapData) {
	println("TO DO")
}
