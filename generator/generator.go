package main

import (
	"fmt"
	"math"
	"math/rand"
)

func almostEqual(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

func degToRad(value float64) float64 {
	return value * math.Pi / 180.
}

func randFloat(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

func generateBorders(n_points int, size_rough_min, size_rough_max float64) {
	angle_step := degToRad(360. / float64(n_points))
	angle_variation_min := -angle_step / 2
	angle_variation_max := angle_step / 2

	for i := 0; i < n_points; i++ {
		angle := angle_step*float64(i) + randFloat(angle_variation_min, angle_variation_max)
		radius := randFloat(size_rough_min/2, size_rough_max/2)

		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		fmt.Printf("%6.2f\t%6.1f\t%6.1f\t%6.1f\n", angle, radius, x, y)
		//fmt.Printf("%7.1f, %7.1f\n", x, y)
	}
}

func main() {
	generateBorders(6, 2000., 3000.)
}
