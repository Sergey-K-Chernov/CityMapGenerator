package main

import (
	"strconv"
	"strings"

	"chirrwick.com/projects/city/city_map"
	"chirrwick.com/projects/city/generator"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type uiBlocksPage struct {
	minSize inputField
	maxSize inputField

	btnGenerate button
	btnAccept   button
	btnBack     button
}

func createBlocksPage() (lay uiBlocksPage) {
	lay.minSize.field.SingleLine = true
	lay.minSize.field.Alignment = text.End
	lay.minSize.label = "Min Size"
	lay.minSize.defaultValue = "100"

	lay.maxSize.field.SingleLine = true
	lay.maxSize.field.Alignment = text.End
	lay.maxSize.label = "Max Size"
	lay.maxSize.defaultValue = "500"

	lay.btnGenerate.label = "Generate blocks"
	lay.btnAccept.label = "Accept blocks"

	lay.btnBack.label = "Back"

	return
}

func (l *uiBlocksPage) Layout(gtx GC, theme *material.Theme) {
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
				makeLabel(theme, l.minSize.label),
				makeFlexInput(gtx, theme, &l.minSize.field, l.minSize.defaultValue),

				makeLabel(theme, l.maxSize.label),
				makeFlexInput(gtx, theme, &l.maxSize.field, l.maxSize.defaultValue),

				makeButton(gtx, theme, &l.btnGenerate.button, l.btnGenerate.label),
				makeButton(gtx, theme, &l.btnAccept.button, l.btnAccept.label),

				layout.Rigid(layout.Spacer{Height: unit.Dp(100)}.Layout),

				makeButton(gtx, theme, &l.btnBack.button, l.btnBack.label),
			)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(mapWidth)}.Layout),
	)
}

func (l *uiBlocksPage) ProcessButtons(gtx GC, ui *uiPages, data *mapData) {
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

func (l *uiBlocksPage) processGenerateButton(data *mapData) {
	var initials generator.InitialValuesBlocks

	inputString := l.minSize.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Size.Min, _ = strconv.ParseFloat(inputString, 32)
	if initials.Size.Min <= 0 {
		initials.Size.Min = 100.0
	}

	inputString = l.maxSize.field.Text()
	inputString = strings.TrimSpace(inputString)
	initials.Size.Max, _ = strconv.ParseFloat(inputString, 32)
	if initials.Size.Max <= 0 {
		initials.Size.Max = 500.0
	}

	go generateBlocks(data.cityMap, data.channel, initials, data.invalidator)
}

func (l *uiBlocksPage) processAcceptButton(data *mapData) uiPage {
	println("TO DO")
	return genBlocksPage
}

func (l *uiBlocksPage) processBackButton(data *mapData) uiPage {
	data.cityMap.Blocks = data.cityMap.Blocks[:0]
	return genBigAreasPage
}

func generateBlocks(cityMap city_map.Map, chanMap chan city_map.Map, initials generator.InitialValuesBlocks, invalidator func()) {
	blocks := generator.GenerateBlocks(cityMap, chanMap, initials)
	cityMap.Blocks = blocks
	chanMap <- cityMap
	invalidator()
}
