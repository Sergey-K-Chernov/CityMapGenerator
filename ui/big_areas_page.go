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

type uiBigAreasPage struct {
	nIndustrialAreas     inputField
	industrialPercentage inputField
	nParksAreas          inputField
	parksPercentage      inputField

	btnGenerate button
	btnAccept   button
	btnBack     button
}

func createBigAreasPage() (lay uiBigAreasPage) {
	lay.nIndustrialAreas.field.SingleLine = true
	lay.nIndustrialAreas.field.Alignment = text.End
	lay.nIndustrialAreas.label = "Industrial areas"
	lay.nIndustrialAreas.defaultValue = "2"

	lay.industrialPercentage.field.SingleLine = true
	lay.industrialPercentage.field.Alignment = text.End
	lay.industrialPercentage.label = "Industrial areas %"
	lay.industrialPercentage.defaultValue = "10"

	lay.nParksAreas.field.SingleLine = true
	lay.nParksAreas.field.Alignment = text.End
	lay.nParksAreas.label = "Park areas"
	lay.nParksAreas.defaultValue = "2"

	lay.parksPercentage.field.SingleLine = true
	lay.parksPercentage.field.Alignment = text.End
	lay.parksPercentage.label = "Parks %"
	lay.parksPercentage.defaultValue = "10"

	lay.btnGenerate.label = "Generate areas"
	lay.btnAccept.label = "Accept areas"
	lay.btnBack.label = "Back"

	return lay
}

func (l *uiBigAreasPage) Layout(gtx GC, theme *material.Theme) {
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

				makeLabel(theme, l.nIndustrialAreas.label),
				makeFlexInput(gtx, theme, &l.nIndustrialAreas.field, l.nIndustrialAreas.defaultValue),

				makeLabel(theme, l.industrialPercentage.label),
				makeFlexInput(gtx, theme, &l.industrialPercentage.field, l.industrialPercentage.defaultValue),

				makeLabel(theme, l.nParksAreas.label),
				makeFlexInput(gtx, theme, &l.nParksAreas.field, l.nParksAreas.defaultValue),

				makeLabel(theme, l.parksPercentage.label),
				makeFlexInput(gtx, theme, &l.parksPercentage.field, l.parksPercentage.defaultValue),

				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),

				makeButton(gtx, theme, &l.btnGenerate.button, l.btnGenerate.label),
				makeButton(gtx, theme, &l.btnAccept.button, l.btnAccept.label),

				layout.Rigid(layout.Spacer{Height: unit.Dp(100)}.Layout),

				makeButton(gtx, theme, &l.btnBack.button, l.btnBack.label),
			)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(mapWidth)}.Layout),
	)
}

func (l *uiBigAreasPage) ProcessButtons(gtx GC, ui *uiPages, data *mapData) {
	if l.btnGenerate.button.Clicked(gtx) {
		l.processGenerateButton(data)
	}
	if l.btnAccept.button.Clicked(gtx) {
		ui.currentPage = l.processAcceptButton(data)
	}
	if l.btnBack.button.Clicked(gtx) {
		ui.currentPage = l.processBackButton(data)
	}
}

func (l *uiBigAreasPage) processGenerateButton(data *mapData) {
	var initials generator.InitialValuesAreas

	inputString := l.nIndustrialAreas.field.Text()
	inputString = strings.TrimSpace(inputString)
	nIndustrialAreas, _ := strconv.ParseInt(inputString, 10, 32)
	if nIndustrialAreas < 1 {
		nIndustrialAreas = 2
	}
	initials.NumIndustrial = int(nIndustrialAreas)

	inputString = l.industrialPercentage.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.AreaIndustrial, _ = strconv.ParseFloat(inputString, 32)
	if initials.AreaIndustrial <= 0 {
		initials.AreaIndustrial = 10.0
	}

	inputString = l.nParksAreas.field.Text()
	inputString = strings.TrimSpace(inputString)
	nParksAreas, _ := strconv.ParseInt(inputString, 10, 32)
	initials.NumParks = int(nParksAreas)
	if initials.NumParks <= 1 {
		initials.NumParks = 2
	}

	inputString = l.parksPercentage.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.AreaParks, _ = strconv.ParseFloat(inputString, 32)
	if initials.AreaParks <= 0 {
		initials.AreaParks = 10.0
	}

	go generateAreas(data.cityMap, data.channel, initials, data.invalidator)
}

func (l *uiBigAreasPage) processAcceptButton(data *mapData) uiPage {
	if len(data.cityMap.Roads) == 0 {
		return genBigAreasPage
	}

	return genBlocksPage
}

func (l *uiBigAreasPage) processBackButton(data *mapData) uiPage {
	data.cityMap.Areas = data.cityMap.Areas[:0]
	return genCentersAndRoadsPage
}

func generateAreas(cityMap generator.Map, chanMap chan generator.Map, initials generator.InitialValuesAreas, invalidator func()) {
	areas := generator.GenerateAreas(cityMap, chanMap, initials)
	cityMap.Areas = areas
	chanMap <- cityMap
	invalidator()
}
