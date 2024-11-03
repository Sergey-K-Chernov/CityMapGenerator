package generator

import (
	"math"
	"sync"

	gm "chirrwick.com/projects/city/generator/genmath"
)

/*

Процедура двуступенчатая:
1. Крупные кварталы
2. Остатки

Генерируем по одному кварталу за раз.
Кидаем прямоугольник. Его центр должен быть в квартале, углы не обязательно.
Находим все пересечения его с дорогами и уже созданными кварталами.
Обрезаем лишнее, полученное преобразовываем в квартал. Информацию о первоначальной ориентации прямоугольника сохраняем.
Делаем так много раз.

*/

func GenerateBlocks(city_map Map, chan_map chan Map, initials InitialValuesBlocks) (blocks []Block) {
	city_area := calcPolygonArea(city_map.BorderPoints, city_map.Center)

	for _, a := range city_map.Areas {
		city_area -= a.Area
	}

	n_blocks := estimateNumberOfBlocks(city_area, initials)

	blocks_area := 0.0

	// genetare initial set of blocks, random
	block_centers := generateRandomPointsInsideCity(n_blocks, city_map)
	blocks, blocks_area = GenerateBlocksInPoints(block_centers, city_map, initials, blocks)

	// fill gaps with less randomly generated blocks
	for i_step := 1; blocks_area < city_area*0.98; i_step++ {
		block_centers = generateConcentricPointsInsideCity(city_map, initials, i_step, blocks)
		var area float64
		blocks, area = GenerateBlocksInPoints(block_centers, city_map, initials, blocks)
		blocks_area += area
	}

	for i := range blocks {
		go func() {
			blocks[i] = generateStreets(blocks[i], initials.Size.Min*2, initials.Size.Max/2)
		}()
	}

	//chan_map <- city_map
	return blocks
}

func GenerateBlocksInPoints(block_centers []gm.Point, city_map Map, initials InitialValuesBlocks, blocks []Block) ([]Block, float64) {
	area := 0.0

	for i := 0; i < len(block_centers); i++ {
		bc := block_centers[i]
		if i > 0 { // debug limit
			//break
		}
		b := generateBlock(bc, city_map, initials, blocks)
		blocks = append(blocks, b)
		block_centers = removePointsInsideFigure(block_centers, b.Points)
		area += calcPolygonArea(b.Points, b.Center)
	}

	return blocks, area
}

func generateBlock(center gm.Point, city_map Map, initials InitialValuesBlocks, blocks []Block) (b Block) {
	side_1 := gm.RandFloat(initials.Size.Min, initials.Size.Max)
	side_2 := gm.RandFloat(initials.Size.Min, initials.Size.Max)
	angle := gm.RandFloat(0, 2*math.Pi)

	b.Center = center

	b.Points = append(b.Points, gm.Point{X: side_1 / 2, Y: side_2 / 2})
	b.Points = append(b.Points, gm.Point{X: side_1 / 2, Y: -side_2 / 2})
	b.Points = append(b.Points, gm.Point{X: -side_1 / 2, Y: -side_2 / 2})
	b.Points = append(b.Points, gm.Point{X: -side_1 / 2, Y: side_2 / 2})

	for i := range b.Points {
		b.Points[i].Rotate(angle)
		b.Points[i].AddInPlace(center)
	}
	b.Angle = angle

	b.Points = cropFigureByRoads(b.Center, b.Points, city_map)
	b.Points = cropFigureByBlocks(b.Center, b.Points, blocks)

	return
}

func cropFigureByRoads(center gm.Point, figure []gm.Point, city_map Map) []gm.Point {
	/*
		println("center")
		println(center.X, center.Y)
		println("figure")
		for _, p := range figure {
			println(p.X, p.Y)
		}
	*/

	max_radius := 0.0
	for _, p := range figure {
		max_radius = max(max_radius, p.Sub(center).Length())
	}

	for _, road := range city_map.Roads {
		for i := range len(road.Points) - 1 {
			figure = cutFigure(center, figure, max_radius, gm.LineSegment{Begin: road.Points[i], End: road.Points[i+1]})
		}
	}
	return figure
}

func cropFigureByBlocks(center gm.Point, figure []gm.Point, blocks []Block) []gm.Point {
	for _, b := range blocks {
		figure = cropFigureByBlock(center, figure, b)
	}
	return figure
}

func cropFigureByBlock(center gm.Point, figure []gm.Point, block Block) []gm.Point {
	max_radius := 0.0
	for _, p := range figure {
		max_radius = max(max_radius, p.Sub(center).Length())
	}

	for i := 0; i < len(block.Points); i++ {
		i_next := i + 1
		if i_next == len(block.Points) {
			i_next = 0
		}

		s := gm.LineSegment{Begin: block.Points[i], End: block.Points[i_next]}
		figure = cutFigure(center, figure, max_radius, s)
	}
	return figure
}

func estimateNumberOfBlocks(area float64, initials InitialValuesBlocks) int {
	//max_ratio := initials.Size.Min / initials.Size.Max

	avg_block_side := (initials.Size.Max + initials.Size.Min) / 2
	avg_square_block_area := avg_block_side * avg_block_side

	// arbitrary:
	est_block_area := avg_square_block_area * 0.7 /* * (max_ratio / 3)*/

	return int(area / est_block_area)
}

func generateRandomPointsInsideCity(qty int, city_map Map) []gm.Point {
	rect := get_map_rect(city_map)

	var wg sync.WaitGroup
	wg.Add(qty)
	points := make([]gm.Point, qty)

	for i := 0; i < qty; i++ {
		go func(i int) {
			defer wg.Done()
			x := gm.RandFloat(rect.Left, rect.Right)
			y := gm.RandFloat(rect.Bottom, rect.Top)
			p := gm.Point{X: x, Y: y}
			for !check_inside_borders(p, city_map) || checkInsideAreas(p, city_map) {
				x = gm.RandFloat(rect.Left, rect.Right)
				y = gm.RandFloat(rect.Bottom, rect.Top)
				p = gm.Point{X: x, Y: y}
			}
			points[i] = p
		}(i)
	}

	wg.Wait()
	return points
}

func checkInsideAreas(point gm.Point, city_map Map) bool {
	for _, area := range city_map.Areas {
		if checkPointInsidePolygon(point, area.Points) {
			return true
		}
	}
	return false
}

func generateConcentricPointsInsideCity(city_map Map, initials InitialValuesBlocks, i_step int, blocks []Block) (points []gm.Point) {
	max_radius := 0.0
	for _, p := range city_map.BorderPoints {
		max_radius = max(max_radius, p.Sub(city_map.Center).Length())
	}

	step := (initials.Size.Min + initials.Size.Max) / float64(i_step)

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for radius := 0.0; radius < max_radius; radius += step {
		angle_step := 2 * math.Atan2(step, radius)

		for angle := 0.0; angle < 2*math.Pi; angle += angle_step {
			func(radius, angle float64) {
				wg.Add(1)
				defer wg.Done()

				point := gm.Point{X: radius, Y: 0}
				point.Rotate(angle)
				point.AddInPlace(city_map.Center)
				point.AddInPlace(generateRadialRandomPoint(0, 2*math.Pi, step/8, step/4))

				if !check_inside_borders(point, city_map) {
					return
				}

				for _, area := range city_map.Areas {
					if checkPointInsidePolygon(point, area.Points) {
						return
					}
				}

				for _, block := range blocks {
					if checkPointInsidePolygon(point, block.Points) {
						return
					}
				}

				mutex.Lock()
				points = append(points, point)
				mutex.Unlock()

			}(radius, angle)
		}
	}

	wg.Wait()

	return
}

func get_map_rect(city_map Map) (rect gm.Rect) {
	rect.Left = 1e10
	rect.Bottom = 1e10

	for _, p := range city_map.BorderPoints {
		rect.Left = math.Min(rect.Left, p.X)
		rect.Right = math.Max(rect.Right, p.X)

		rect.Bottom = math.Min(rect.Bottom, p.X)
		rect.Top = math.Max(rect.Top, p.X)
	}

	return
}

func check_inside_borders(p gm.Point, city_map Map) bool {
	return checkPointInsidePolygon(p, city_map.BorderPoints)
}

func generateStreets(block Block, min_dist, max_dist float64) Block {
	max_length := 0.0
	for _, p := range block.Points {
		max_length = max(max_length, p.Sub(block.Center).Length())
	}

	base_point := gm.Point{X: max_length * 2, Y: 0}
	p1 := base_point
	p2 := base_point
	p3 := base_point
	p4 := base_point

	p1.Rotate(block.Angle)
	p2.Rotate(block.Angle + math.Pi)

	p3.Rotate(block.Angle + math.Pi/2)
	p4.Rotate(block.Angle + 3*math.Pi/2)

	p1.AddInPlace(block.Center)
	p2.AddInPlace(block.Center)
	p3.AddInPlace(block.Center)
	p4.AddInPlace(block.Center)

	dist := gm.RandFloat(min_dist, max_dist/3)
	for shift := -max_length - min_dist/3; shift < max_length; shift += dist {
		shift_point := gm.Point{X: shift + gm.RandFloat(-shift/2, shift/2), Y: 0}
		shift_point.Rotate(block.Angle + math.Pi/2)

		street, ok := tryMakeStreet(p1, p2, shift_point, block)

		if ok {
			block.Streets = append(block.Streets, street)
		}

	}

	dist = gm.RandFloat(max_dist/3, max_dist)
	for shift := -max_length - min_dist/3; shift < max_length; shift += dist {
		shift_point := gm.Point{X: shift + gm.RandFloat(-shift/2, shift/2), Y: 0}
		shift_point.Rotate(block.Angle)

		street, ok := tryMakeStreet(p3, p4, shift_point, block)

		if ok {
			block.Streets = append(block.Streets, street)
		}

	}

	return block
}

func tryMakeStreet(p1, p2, shift_point gm.Point, block Block) (gm.LineSegment, bool) {
	street := gm.LineSegment{Begin: p1.Add(shift_point), End: p2.Add(shift_point)}

	points := make([]gm.Point, 0)
	for i := range block.Points {
		i_next := i + 1
		if i_next == len(block.Points) {
			i_next = 0
		}

		s := gm.LineSegment{Begin: block.Points[i], End: block.Points[i_next]}

		p, ok := s.Intersect(street)
		if ok {
			points = append(points, p)
		}
	}
	if len(points) == 2 {
		return gm.LineSegment{Begin: points[0], End: points[1]}, true
	}
	return gm.LineSegment{Begin: gm.Point{X: 0, Y: 0}, End: gm.Point{X: 0, Y: 0}}, false
}
