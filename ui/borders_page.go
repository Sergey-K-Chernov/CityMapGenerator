package main

import (
	"strconv"
	"strings"

	"chirrwick.com/projects/city/generator"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func createBordersPage() (lay uiBordersPage) {
	lay.nPoints.field.SingleLine = true
	lay.nPoints.field.Alignment = text.End
	lay.nPoints.label = "Corners"
	lay.nPoints.defaultValue = "3"

	lay.minRadius.field.SingleLine = true
	lay.minRadius.field.Alignment = text.End
	lay.minRadius.label = "Min radius"
	lay.minRadius.defaultValue = "1000"

	lay.maxRadius.field.SingleLine = true
	lay.maxRadius.field.Alignment = text.End
	lay.maxRadius.label = "Max radius"
	lay.maxRadius.defaultValue = "3000"

	lay.pointVariation.field.SingleLine = true
	lay.pointVariation.field.Alignment = text.End
	lay.pointVariation.label = "Variation"
	lay.pointVariation.defaultValue = "300"

	lay.btnGenerate.label = "Generate"
	lay.btnAccept.label = "Accept borders"

	return
}

func (l *uiBordersPage) Layout(gtx GC, theme *material.Theme) {
	// Fixed ui hack...
	var uiFlexWeight float32
	totalWidth := gtx.Constraints.Max.X
	mapWidth := totalWidth - UI_WIDTH

	uiFlexWeight = float32(UI_WIDTH) / float32(totalWidth)
	//mapFlexWeight = float32(mapWidth) / float32(totalWidth)

	layout.Flex{}.Layout(gtx,
		layout.Flexed(uiFlexWeight, func(gtx GC) Dims {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,

				makeLabel(theme, l.nPoints.label),
				makeFlexInput(gtx, theme, &l.nPoints.field, l.nPoints.defaultValue),

				makeLabel(theme, l.minRadius.label),
				makeFlexInput(gtx, theme, &l.minRadius.field, l.minRadius.defaultValue),

				makeLabel(theme, l.maxRadius.label),
				makeFlexInput(gtx, theme, &l.maxRadius.field, l.maxRadius.defaultValue),

				makeLabel(theme, l.pointVariation.label),
				makeFlexInput(gtx, theme, &l.pointVariation.field, l.pointVariation.defaultValue),

				makeButton(gtx, theme, &l.btnGenerate.button, l.btnGenerate.label),
				makeButton(gtx, theme, &l.btnAccept.button, l.btnAccept.label),
			)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(mapWidth)}.Layout),
	)
}

func (l *uiBordersPage) ProcessButtons(gtx GC, ui *uiPages, data *mapData) {
	if l.btnGenerate.button.Clicked(gtx) {
		l.processGenerateButton(gtx, data)
	}
	if l.btnAccept.button.Clicked(gtx) {
		ui.currentPage = l.processAcceptButton(gtx, data)
	}
}

func (l *uiBordersPage) processGenerateButton(gtx GC, data *mapData) {
	var initials generator.InitialValuesMap

	inputString := l.nPoints.field.Text()
	inputString = strings.TrimSpace(inputString)
	nSides, _ := strconv.ParseInt(inputString, 10, 32)
	if nSides < 3 {
		nSides = 3
	}
	initials.NumSides = int(nSides)

	inputString = l.minRadius.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Min, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Min <= 0 {
		initials.Raduis.Min = 2000.0
	}

	inputString = l.maxRadius.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Max, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Max <= 0 {
		initials.Raduis.Max = 3000.0
	}

	inputString = l.pointVariation.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.VertexShift, _ = strconv.ParseFloat(inputString, 32)
	if initials.VertexShift < 0.0 {
		initials.VertexShift = initials.Raduis.Max / 10.0
	}

	go generateBorders(data.channel, initials, data.invalidator)
}

func (l *uiBordersPage) processAcceptButton(gtx GC, data *mapData) uiPage {
	if len(data.cityMap.BorderPoints) == 0 {
		return genBordersPage
	}

	return genCentersAndRoadsPage
}

func generateBorders(chanMap chan generator.Map, initials generator.InitialValuesMap, invalidator func()) {
	generator.GenerateBorders(chanMap, initials)
	invalidator()
}
