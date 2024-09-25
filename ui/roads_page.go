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

type uiCentersAndRoadsPage struct {
	nCenters  inputField
	minRadius inputField
	maxRadius inputField
	branching inputField

	btnGenerate button
	btnAccept   button
}

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
	lay.branching.defaultValue = "0"

	lay.btnGenerate.label = "Generate"
	lay.btnAccept.label = "Accept roads"

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
	if l.btnGenerate.button.Clicked(gtx) {
		l.processGenerateButton(gtx, data)
	}
	if l.btnAccept.button.Clicked(gtx) {
		ui.currentPage = l.processAcceptButton(gtx, data)
	}
}

func (l *uiCentersAndRoadsPage) processGenerateButton(gtx GC, data *mapData) {
	var initials generator.InitialValuesRoads

	inputString := l.nCenters.field.Text()
	inputString = strings.TrimSpace(inputString)
	nSides, _ := strconv.ParseInt(inputString, 10, 32)
	if nSides < 1 {
		nSides = 1
	}
	initials.NumCenters = int(nSides)

	inputString = l.minRadius.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Min, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Min <= 0 {
		initials.Raduis.Min = 500.0
	}

	inputString = l.maxRadius.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Raduis.Max, _ = strconv.ParseFloat(inputString, 32)
	if initials.Raduis.Max <= 0 {
		initials.Raduis.Max = 1000.0
	}

	inputString = l.branching.field.Text()
	inputString = strings.TrimSpace(inputString)
	br, _ := strconv.ParseInt(inputString, 10, 32)
	initials.Branching = int(br)

	go generateRoads(data.cityMap, data.channel, initials, data.invalidator)
}

func (l *uiCentersAndRoadsPage) processAcceptButton(gtx GC, data *mapData) uiPage {
	println("TO DO")
	return genCentersAndRoadsPage
}

func generateRoads(cityMap generator.Map, chanMap chan generator.Map, initials generator.InitialValuesRoads, invalidator func()) {
	roads := generator.GenerateRoads(cityMap, chanMap, initials)
	cityMap.Roads = roads
	chanMap <- cityMap
	invalidator()
}
