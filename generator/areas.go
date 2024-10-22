package generator

import (
	"math"

	gm "chirrwick.com/projects/city/generator/genmath"
)

func GenerateAreas(cityMap Map, chanMap chan Map, initials InitialValuesAreas) (areas []Area) {
	total_area := calcArea(cityMap.BorderPoints, cityMap.Center)
	areas = append(areas, generateIndustrialAreas(cityMap, initials, total_area)...)
	areas = append(areas, generateParksAreas(cityMap, initials, total_area)...)
	return
}

func generateIndustrialAreas(cityMap Map, initials InitialValuesAreas, total_area float64) (areas []Area) {
	target_areas := calcTargetAreas(total_area*initials.AreaIndustrial/100, initials.NumIndustrial)

	for i := 0; i < len(target_areas); i++ {
		area := generateArea(target_areas[i])
		shiftIndustrialArea(&area, cityMap)
		area.Type = AreaIndustrial

		areas = append(areas, area)
	}

	return
}

func generateParksAreas(cityMap Map, initials InitialValuesAreas, total_area float64) (areas []Area) {
	target_areas := calcTargetAreas(total_area*initials.AreaParks/100, initials.NumParks)

	for i := 0; i < len(target_areas); i++ {
		area := generateArea(target_areas[i])
		shiftParkArea(&area, cityMap)
		area.Type = AreaPark

		areas = append(areas, area)
	}

	return
}

func shiftParkArea(area *Area, cityMap Map) {
	for i := range area.Points {
		area.Points[i].AddInPlace(cityMap.Center)
	}

	max_length := 0.0
	for _, p := range cityMap.BorderPoints {
		max_length = max(max_length, p.Sub(cityMap.Center).Length())
	}

	angle := gm.RandFloat(0, 2*math.Pi)
	end_point := gm.Point{X: max_length * 2, Y: 0}
	end_point.Rotate(angle)
	s := gm.LineSegment{Begin: cityMap.Center, End: cityMap.Center.Add(end_point)}

	if p_border, _, err := intersectWithFigure(s, cityMap.BorderPoints); err == nil {
		distance_border := p_border.Sub(cityMap.Center).Length()

		shift := gm.RandFloat(distance_border*0.2, distance_border)
		shift_point := gm.Point{X: shift, Y: 0}
		shift_point.Rotate(angle)

		for i := range area.Points {
			area.Points[i].AddInPlace(shift_point)
		}

	}

}

func shiftIndustrialArea(area *Area, cityMap Map) {
	for i := range area.Points {
		area.Points[i].AddInPlace(cityMap.Center)
	}

	max_length := 0.0
	for _, p := range cityMap.BorderPoints {
		max_length = max(max_length, p.Sub(cityMap.Center).Length())
	}

	angle := gm.RandFloat(0, 2*math.Pi)
	end_point := gm.Point{X: max_length * 2, Y: 0}
	end_point.Rotate(angle)
	s := gm.LineSegment{Begin: cityMap.Center, End: cityMap.Center.Add(end_point)}

	if p_border, _, err := intersectWithFigure(s, cityMap.BorderPoints); err == nil {
		distance_border := p_border.Sub(cityMap.Center).Length()

		if p_area, _, err := intersectWithFigure(s, area.Points); err == nil {
			distance_area := p_area.Sub(cityMap.Center).Length()
			shift := distance_border - distance_area
			shift_point := gm.Point{X: shift, Y: 0}
			shift_point.Rotate(angle)

			for i := range area.Points {
				area.Points[i].AddInPlace(shift_point)
			}
		}
	}
}

func calcTargetAreas(total_area float64, number_of_areas int) []float64 {
	target_areas := make([]float64, number_of_areas)
	min_area := total_area / float64(number_of_areas) * 0.8
	max_area := total_area / float64(number_of_areas) * 1.2

	areas_sum := 0.0
	for i := 0; i < number_of_areas-1; i++ {
		target_areas[i] = gm.RandFloat(min_area, max_area)
		areas_sum += target_areas[i]
	}
	target_areas[number_of_areas-1] = total_area - areas_sum
	return target_areas
}

func generateArea(target_area float64) Area {
	var area Area
	n_points := gm.RandInt(4, 6)
	for i := 0; i < n_points; i++ {
		angle := float64(i) * (2 * math.Pi / float64(n_points))
		angle_delta := math.Pi / float64(n_points) / 2

		angle += gm.RandFloat(-angle_delta, angle_delta)

		area.Points = append(area.Points, generateRadialRandomPoint(angle-angle_delta, angle+angle_delta, 8, 12))
	}

	current_area_area := calcArea(area.Points, gm.Point{X: 0, Y: 0})
	scale := math.Sqrt(target_area / current_area_area)

	for i := range area.Points {
		area.Points[i].Scale(scale)
	}

	area.Area = calcArea(area.Points, gm.Point{X: 0, Y: 0})
	return area
}
