package generator

import (
	"math"

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

	//side_1 := gm.RandFloat(initials.Size.Min, initials.Size.Max)
	//side_2 := gm.RandFloat(initials.Size.Min, initials.Size.Max)

	for _, bc := range block_centers {
		var b Block
		b.Points = append(b.Points, bc.Add(gm.Point{X: 0, Y: 10}))
		b.Points = append(b.Points, bc.Add(gm.Point{X: 10, Y: 0}))
		b.Points = append(b.Points, bc.Add(gm.Point{X: 0, Y: -10}))
		b.Points = append(b.Points, bc.Add(gm.Point{X: -10, Y: 0}))
		blocks = append(blocks, b)
	}

	//chan_map <- city_map
	return blocks
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

	points := make([]gm.Point, qty)

	for i := 0; i < qty; i++ {
		func(i int) {
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

// Consider figure is convex polygon
func check_inside_borders(p gm.Point, city_map Map) bool {
	bp := city_map.BorderPoints
	for i := range bp {
		i_plus_1 := i + 1
		if i_plus_1 == len(bp) {
			i_plus_1 = 0
		}

		tri := gm.Triangle{A: city_map.Center, B: bp[i], C: bp[i_plus_1]}

		if tri.HasPoint(p) {
			return true
		}
	}
	return false
}
