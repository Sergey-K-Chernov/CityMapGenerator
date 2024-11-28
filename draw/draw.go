package draw

import (
	"image"
	"image/color"
	"math"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"chirrwick.com/projects/city/city_map"
	_ "chirrwick.com/projects/city/generator/genmath"
)

const IMAGE_SIZE = 512

type drawData struct {
	Image *image.RGBA
	Gc draw2d.GraphicContext
	Map city_map.Map
	Scale float64
}

func Draw(city_map city_map.Map) *image.RGBA {
	var data drawData
	data.Image = image.NewRGBA(image.Rect(0, 0, IMAGE_SIZE, IMAGE_SIZE))
	data.Gc = draw2dimg.NewGraphicContext(data.Image)
	data.Map = city_map
	data.Scale = calcScale(city_map)

	drawBorders(&data)
	drawAreas(&data)
	drawBlocks(&data)
	drawRoads(&data)

	return data.Image
}

func calcScale(city_map city_map.Map) float64 {
	maximum := 0.0
	for _, pnt := range city_map.BorderPoints {
		maximum = math.Max(maximum, pnt.X)
		maximum = math.Max(maximum, pnt.Y)
	}
	maximum *= 1.05 // margin
	return IMAGE_SIZE/maximum
}

func drawBorders(data *drawData) {
	data.Gc.SetFillColor(color.RGBA{0xdd, 0xff, 0xdd, 0xff})
	data.Gc.SetLineWidth(0)

	data.Gc.BeginPath()

	end_pnt := data.Map.BorderPoints[len(data.Map.BorderPoints) - 1]
	data.Gc.MoveTo(end_pnt.X * data.Scale, end_pnt.Y * data.Scale)
	for _, pnt := range data.Map.BorderPoints {
		data.Gc.LineTo(pnt.X * data.Scale, pnt.Y * data.Scale)
	}
	data.Gc.Close()
	data.Gc.FillStroke()
}


func drawAreas(data *drawData) {

}

func drawBlocks(data *drawData) {

}

func drawRoads(data *drawData) {

}

