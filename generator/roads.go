package generator

import (
	"sort"

	"chirrwick.com/projects/city/city_map"
	gm "chirrwick.com/projects/city/generator/genmath"
)

func GenerateRoads(cityMap city_map.Map, chanMap chan city_map.Map, initials InitialValuesRoads) (roads []city_map.Road) {
	centers := generateCenters(cityMap, initials)
	roads = append(roads, connectCenters(centers)...)
	exits := generateExits(cityMap, initials)
	roads = append(roads, connectCentersWithExits(centers, exits)...)
	return
}

func generateCenters(cityMap city_map.Map, initials InitialValuesRoads) []gm.Point {
	nCenters := initials.NumCenters
	rMin := initials.Raduis.Min
	rMax := initials.Raduis.Max
	angle_step := gm.DegToRad(360. / float64(nCenters))
	angle_variation := angle_step / 4

	centers := make([]gm.Point, nCenters)

	for i := 0; i < nCenters; i++ {
		angle := angle_step * float64(i)
		centers[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, rMin, rMax)
		centers[i].AddInPlace(cityMap.Center)
	}

	return centers
}

func connectCenters(centers []gm.Point) (roads []city_map.Road) {
	for i := 0; i < len(centers)-1; i++ {
		var rd city_map.Road
		rd.Points = append(rd.Points, centers[i])
		for j := i + 1; j < len(centers); j++ {
			rd.Points = append(rd.Points, centers[j])
		}
		rd.Points = append(rd.Points, centers[i])
		roads = append(roads, rd)
	}
	return
}

func generateExits(cityMap city_map.Map, initials InitialValuesRoads) []gm.Point {
	angle_step := gm.DegToRad(360. / float64(initials.Branching))
	angle_variation := angle_step / 4

	intermediate_exit_vectors := make([]gm.Point, initials.Branching)

	for i := 0; i < initials.Branching; i++ {
		angle := angle_step/2 + angle_step*float64(i)
		intermediate_exit_vectors[i] = generateRadialRandomPoint(angle-angle_variation, angle+angle_variation, 1, 1)
		intermediate_exit_vectors[i].AddInPlace(cityMap.Center)
	}

	exit_points := make([]gm.Point, 0)
	for _, exit := range intermediate_exit_vectors {
		sg := gm.LineSegment{Begin: cityMap.Center, End: exit}
		extendSegment(&sg, cityMap.BorderPoints)
		if ex, _, err := intersectSegmentWithFigure(sg, cityMap.BorderPoints); err == nil {
			exit_points = append(exit_points, ex)
		} else {
			println(ex.X, ex.Y)
		}
	}

	return exit_points
}

func connectCentersWithExits(centers []gm.Point, exits []gm.Point) (roads []city_map.Road) {
	roadFans := make([]Fan, len(centers))
	for i := range roadFans {
		roadFans[i] = makeFan(centers[i])
	}

	for _, ep := range exits {
		i := findClosestPointIndex(ep, centers)
		roadFans[i].outer = append(roadFans[i].outer, makeFan(ep))
	}

	sort.Slice(roadFans[0].outer, func(a, b int) bool {
		angleA := roadFans[0].outer[a].root.Sub((roadFans[0].root)).Angle()
		angleB := roadFans[0].outer[b].root.Sub((roadFans[0].root)).Angle()
		return angleA < angleB
	})

	for _, fan := range roadFans {
		fan.split()
		roads = append(roads, fan.makeRoads()...)
	}
	return
}

func (f *Fan) makeRoads() (roads []city_map.Road) {
	for _, subfan := range f.outer {
		var rd city_map.Road
		rd.Points = append(rd.Points, f.root, subfan.root)
		roads = append(roads, rd)
		roads = append(roads, subfan.makeRoads()...)
	}
	return
}
