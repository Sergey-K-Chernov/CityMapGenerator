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
	city_area := calcArea(city_map)
	n_blocks := estimateNumberOfBlocks(city_area, initials)

	block_centers := generatePointsInsideCity(n_blocks, city_map)

	for i, bc := range block_centers {
		if i > 0 { // debug limit
			break
		}
		b := generateBlock(bc, city_map, initials)
		blocks = append(blocks, b)
	}

	//chan_map <- city_map
	return blocks
}

func generateBlock(center gm.Point, city_map Map, initials InitialValuesBlocks) (b Block) {
	side_1 := gm.RandFloat(initials.Size.Min, initials.Size.Max)
	side_2 := gm.RandFloat(initials.Size.Min, initials.Size.Max)
	angle := gm.RandFloat(0, 2*math.Pi)

	b.Points = append(b.Points, gm.Point{X: side_1 / 2, Y: side_2 / 2})
	b.Points = append(b.Points, gm.Point{X: side_1 / 2, Y: -side_2 / 2})
	b.Points = append(b.Points, gm.Point{X: -side_1 / 2, Y: -side_2 / 2})
	b.Points = append(b.Points, gm.Point{X: -side_1 / 2, Y: side_2 / 2})

	for i := range b.Points {
		b.Points[i].Rotate(angle)
		b.Points[i].AddInPlace(center)
	}

	b.Points = cropBlock(center, b.Points, city_map)

	return
}

func cropBlock(center gm.Point, figure []gm.Point, city_map Map) []gm.Point {
	println("center")
	println(center.X, center.Y)
	println("figure")
	for _, p := range figure {
		println(p.X, p.Y)
	}
	for _, road := range city_map.Roads {
		for i := range len(road.Points) - 1 {
			figure = cutFigure(center, figure, gm.LineSegment{Begin: road.Points[i], End: road.Points[i+1]})
		}
	}
	return figure
}

type Intersection struct {
	i, i_next int
	point     gm.Point
}

func cutFigure(center gm.Point, figure []gm.Point, segment gm.LineSegment) []gm.Point {
	croped := false
	defer func(segment gm.LineSegment) {
		if croped {
			println("seg:")
			println(segment.Begin.X, segment.Begin.Y)
			println(segment.End.X, segment.End.Y)
		}
	}(segment)

	np := segment.GetNormalPoint(center)
	np = np.Sub(center)

	if np.LengthSq() == 0 {
		np = segment.End.Sub(center)
		np.Normalize()
		np.Rotate(math.Pi / 2)
	}

	// Смещаем и откладываем обратное смещение
	segment.Begin = segment.Begin.Sub(center)
	segment.End = segment.End.Sub(center)

	for i := 0; i < len(figure); i++ {
		figure[i] = figure[i].Sub(center)
	}

	defer func() {
		for i := 0; i < len(figure); i++ {
			figure[i] = figure[i].Add(center)
		}
	}()

	// Поворачиваем и откладываем обратный поворот
	angle := np.Angle()

	segment.Begin.Rotate(-angle)
	segment.End.Rotate(-angle)

	if segment.Begin.Y > segment.End.Y {
		segment.Begin, segment.End = segment.End, segment.Begin
	}

	for i := 0; i < len(figure); i++ {
		figure[i].Rotate(-angle)
	}

	defer func() {
		for i := 0; i < len(figure); i++ {
			figure[i].Rotate(angle)
		}
	}()

	// Пересекаем прямой. Получаем точки пересечения и индексы отрезков-сторон которые пересечены
	// Почему не отрезком: если отрезок кончается внутри фигуры, мы все равно хотим отрезать.
	intersections := make([]Intersection, 0)
	for i, p := range figure {
		i_next := i + 1
		if i_next == len(figure) {
			i_next = 0
		}

		point, ok := gm.LineSegment{Begin: p, End: figure[i_next]}.IntersectLine(segment)
		if !ok {
			continue
		}
		intersections = append(intersections, Intersection{i: i, i_next: i_next, point: point})
	}

	if len(intersections) < 2 {
		return figure
	}

	if len(intersections) > 2 {
		println("Error! Non-convex figure")
	}
	// Проверяем, пересекает ли отрезок:
	// 	Если нижняя точка отрезка выше верхней точки пересечения прямой, то нет
	//	Если верхняя точка отрезка ниже нижней точки пересечения прямой, то нет

	if intersections[0].point.Y < intersections[1].point.Y {
		intersections[0], intersections[1] = intersections[1], intersections[0]
	}

	if segment.Begin.Y > intersections[0].point.Y || segment.End.Y < intersections[1].point.Y {
		return figure
	}
	println()

	// Иначе есть пересечения. Режем:

	figure = cutPoints(figure, intersections)

	croped = true
	return figure
	// При выходе deferred обратно вращают и смещают нас.
}

func cutPoints(figure []gm.Point, intersections []Intersection) []gm.Point {
	// Sort by increase
	if intersections[0].i < intersections[1].i {
		intersections[0], intersections[1] = intersections[1], intersections[0]
	}

	remove_begin := -1
	remove_end := -1

	figure = slices.Insert(figure, intersections[0].i_next, intersections[0].point)
	if figure[intersections[0].i_next].X > intersections[0].point.X {
		remove_begin = intersections[0].i_next + 1
	} else {
		remove_end = intersections[0].i
	}

	if figure[intersections[1].i_next].X > intersections[1].point.X {
		// Второе пересечениие. Если вправо, вставляем по правому индексу
		figure = slices.Insert(figure, intersections[1].i_next, intersections[1].point)

		// На удаление помечаем следующий за правым, начало
		remove_begin = intersections[1].i_next + 1
	} else {
		// Если влево, то по левому.
		figure = slices.Insert(figure, intersections[1].i, intersections[1].point)

		// На удаление помечаем предыдущий перед левым, конец
		remove_end = intersections[1].i - 1
	}

	if remove_begin > remove_end {
		figure = slices.Delete(figure, remove_begin, len(figure))
		figure = slices.Delete(figure, 0, remove_end)
	} else {
		figure = slices.Delete(figure, remove_begin, remove_end)
	}

	return figure
	//figure = slices.Insert(figure, , )
}

func calcArea(city_map Map) float64 {
	bp := city_map.BorderPoints

	area := 0.0
	for i := range bp {
		i_plus_1 := i + 1
		if i_plus_1 == len(bp) {
			i_plus_1 = 0
		}

		// Heron's formula
		a := bp[i].Sub(bp[i_plus_1]).Length()
		b := bp[i].Sub(city_map.Center).Length()
		c := bp[i_plus_1].Sub(city_map.Center).Length()
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
	est_block_area := avg_square_block_area /* * (max_ratio / 3)*/

	return int(area / est_block_area)
}

func generatePointsInsideCity(qty int, city_map Map) []gm.Point {
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
			for !check_inside_borders(p, city_map) {
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
