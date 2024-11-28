package draw

import (
	"image"
	"image/color"
	"math"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"chirrwick.com/projects/city/city_map"
	gm "chirrwick.com/projects/city/generator/genmath"
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
	if len(data.Map.Areas) <= 0 {
		return
	}

	for _, area := range data.Map.Areas {
		switch area.Type {
			case city_map.AreaIndustrial: 
				data.Gc.SetFillColor(color.RGBA{0xdd, 0xdd, 0xdd, 0xff})
			case city_map.AreaPark:
				data.Gc.SetFillColor(color.RGBA{0xcc, 0xff, 0xcc, 0xff})
		}

		data.Gc.BeginPath()

		end_pnt := area.Points[len(area.Points)-1]
		data.Gc.MoveTo(end_pnt.X * data.Scale, end_pnt.Y * data.Scale)

		for _, pnt := range area.Points {
			data.Gc.LineTo(pnt.X * data.Scale, pnt.Y * data.Scale)
		}

		data.Gc.Close()
		data.Gc.Fill()
	}
}

func drawBlocks(data *drawData) {
	if len(data.Map.Blocks) <= 0 {
		return
	}

	for _, block := range data.Map.Blocks {
		data.Gc.SetFillColor(color.RGBA{0xff, 0xf6, 0xdd, 0xff})
		data.Gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
		data.Gc.SetLineWidth(1)
		data.Gc.BeginPath()

		end_pnt := block.Points[len(block.Points) - 1]
		data.Gc.MoveTo(end_pnt.X * data.Scale, end_pnt.Y * data.Scale)
		for _, pnt := range block.Points {
			data.Gc.LineTo(pnt.X * data.Scale, pnt.Y * data.Scale)
		}

		data.Gc.Close()
		data.Gc.FillStroke()
		drawStreets(data, block.Streets)
	}
}

func drawStreets(data *drawData, streets []gm.LineSegment) {
	data.Gc.SetLineWidth(0.7)
	for _, str := range streets {
		data.Gc.BeginPath()

		data.Gc.MoveTo(str.Begin.X * data.Scale, str.Begin.Y * data.Scale)
		data.Gc.LineTo(str.End.X * data.Scale, str.End.Y * data.Scale)

		data.Gc.Close()
		data.Gc.Stroke()
	}

}

func drawRoads(data *drawData) {
	if len(data.Map.Roads) <= 0 {
		return
	}
	
	// Colors need unification with ui?
	data.Gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	data.Gc.SetLineWidth(2)

	for _, rd := range data.Map.Roads {
		data.Gc.BeginPath()
		
		pnt := rd.Points[0]
		data.Gc.MoveTo(pnt.X * data.Scale, pnt.Y * data.Scale)
		for _, pnt := range rd.Points {
			data.Gc.LineTo(pnt.X * data.Scale, pnt.Y * data.Scale)
		}
		
		data.Gc.Close()
		data.Gc.Stroke()
	}
}

