package generator

import (
	"math"
	"slices"
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

Кварталы храним так: отдельно - набор сторон; отдельно - многоугольник, который хранит индексы сторон, из которых он состоит
Сторона также знает, каким кварталам она принадлежит.
Храним как? в слайсе? сортируем как? чтобы быстрее искать ближайшие?  И надо ли искать ближайшие?

Когда нагенерим основную массу кварталов, пробегаемся по всем сторонам и смотрим, скольки кварталам они принадлежат.
Если двум, то все хорошо. Если одному, то надо найти соседние стороны и сделать квартал из них.

Исключение - стороны, выпирающие за границу города (или лежащие на ней)

*/

func GenerateBlocks(city_map Map, chan_map chan Map, initials InitialValuesBlocks) (blocks []Block) {
	city_area := calcArea(city_map.BorderPoints, city_map.Center)

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
		area += calcArea(b.Points, b.Center)
	}

	return blocks, area
}

func removePointsInsideFigure(points, figure []gm.Point) []gm.Point {
	for i := 0; i < len(points); {
		if isPointInsideFigure(points[i], figure) {
			points = slices.Delete(points, i, i+1)
		} else {
			i++
		}
	}
	return points
}

func isPointInsideFigure(point gm.Point, figure []gm.Point) bool {
	line := gm.LineSegment{Begin: point, End: point.Add(gm.Point{X: 1, Y: 0})}

	intersections := make([]gm.Point, 0)
	for i := range figure {
		i_next := i + 1
		if i_next == len(figure) {
			i_next = 0
		}
		s := gm.LineSegment{Begin: figure[i], End: figure[i_next]}

		isec, ok := s.IntersectLine(line)
		if ok {
			intersections = append(intersections, isec)
		}
	}

	numberOfRayIntersections := 0
	for _, isec := range intersections {
		if isec.X > point.X {
			numberOfRayIntersections++
		}
	}

	return (numberOfRayIntersections%2 != 0)
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

	b.Points = cropBlockByRoads(b.Center, b.Points, city_map)
	b.Points = cropBlockByBlocks(b.Center, b.Points, blocks)

	return
}

func cropBlockByRoads(center gm.Point, figure []gm.Point, city_map Map) []gm.Point {
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

func cropBlockByBlocks(center gm.Point, figure []gm.Point, blocks []Block) []gm.Point {
	for _, b := range blocks {
		figure = cropBlockByBlock(center, figure, b)
	}
	return figure
}

func cropBlockByBlock(center gm.Point, figure []gm.Point, block Block) []gm.Point {
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

type Intersection struct {
	i, i_next int
	point     gm.Point
}

func cutFigure(center gm.Point, figure []gm.Point, max_radius float64, segment gm.LineSegment) []gm.Point {
	// Находим перпендикуляр из центра на отрезок
	np := segment.GetNormalPoint(center)
	n := np.Sub(center)

	// Если больше расстояния до самого дальнего угла, игнорим
	normal_length := n.Length()
	if normal_length > max_radius {
		return figure
	}

	// Находим пересечения отрезка с обеими сторонами
	intersections := make([]Intersection, 0)
	for i, p := range figure {
		i_next := i + 1
		if i_next == len(figure) {
			i_next = 0
		}

		s := gm.LineSegment{Begin: p, End: figure[i_next]}
		point, ok := s.IntersectLine(segment)
		if ok {
			intersections = append(intersections, Intersection{i: i, i_next: i_next, point: point})
		}
	}

	if len(intersections) < 2 {
		return figure
	}

	// Проверяем, находятся ли точки перечечения внутри отрезка
	line_vec := segment.End.Sub(segment.Begin)
	line_vec_ort := line_vec.GetNormalized()
	ok := false
	for i := 0; i < len(intersections); i++ {
		isec_vec := intersections[i].point.Sub(segment.Begin)

		projection := line_vec_ort.Dot(isec_vec)
		if projection > 0 && projection < line_vec.Length() {
			ok = true
			break
		}
	}

	if !ok {
		return figure
	}

	slices.SortFunc(intersections, func(a, b Intersection) int {
		if a.i_next < b.i_next {
			return 1
		}
		if a.i_next > b.i_next {
			return -1
		}
		return 0
	})

	// Вставляем точки в фигуру
	for _, isec := range intersections {
		figure = slices.Insert(figure, isec.i_next, isec.point)
	}

	// Находим проекции каждой точки на перпендикуляр. Если отрицательные, игнорим

	ortho := gm.Line{Origin: center, Vector: n.GetNormalized()}
	nlen := n.Length()
	for i := len(figure) - 1; i >= 0; i-- {
		vec_to_corner := figure[i].Sub(center)
		projection := ortho.Vector.Dot(vec_to_corner)
		if gm.AlmostEqual(projection, nlen, 0.0000001) {
			continue
		}
		if projection > nlen { // < 0 is automacitally < length
			figure = slices.Delete(figure, i, i+1)
		}
	}

	return figure
}

func calcArea(polygon []gm.Point, center gm.Point) float64 {
	area := 0.0
	for i := range polygon {
		i_plus_1 := i + 1
		if i_plus_1 == len(polygon) {
			i_plus_1 = 0
		}

		// Heron's formula
		a := polygon[i].Sub(polygon[i_plus_1]).Length()
		b := polygon[i].Sub(center).Length()
		c := polygon[i_plus_1].Sub(center).Length()
		p := (a + b + c) / 2

		area += math.Sqrt(p * (p - a) * (p - b) * (p - c))
	}
	return area
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

// Consider figure is convex polygon - переписать на
func check_inside_borders(p gm.Point, city_map Map) bool {
	return checkPointInsidePolygon(p, city_map.BorderPoints)
}

func checkPointInsidePolygon(p gm.Point, poly []gm.Point) bool {
	sum := 0
	for i := range poly {
		i_plus_1 := i + 1
		if i_plus_1 == len(poly) {
			i_plus_1 = 0
		}

		sum += checkXRayIntersectsSection(p, poly[i], poly[i_plus_1])
	}
	return sum%2 > 0
}

func checkXRayIntersectsSection(p, a, b gm.Point) int {
	ax := a.X - p.X
	ay := a.Y - p.Y

	bx := b.X - p.X
	by := b.Y - p.Y

	intersect_axis_x := (ay*by <= 0)
	intersect_axis_y := (ax*bx <= 0)

	// Send ray to the right, test for intersecting positive x axis .

	// No intersection at all:
	if !intersect_axis_x {
		return 0
	}

	if !intersect_axis_y {
		// Both to the left - no intersection
		if ax < 0 {
			return 0
		}
		// Both to the right - have intersection
		return 1
	}

	// find intersection point
	fraction := math.Abs(ay / (by - ay))
	x := ax + fraction*(bx-ax)

	if x < 0 {
		return 0
	}

	return 1
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
