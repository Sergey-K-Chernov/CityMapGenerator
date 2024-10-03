package generator

import (
	"math"
	"sort"

	gm "chirrwick.com/projects/city/generator/genmath"
)

type Fan struct {
	root  gm.Point
	outer []Fan
}

func makeFan(root gm.Point) (f Fan) {
	f.root = root
	f.outer = make([]Fan, 0)
	return
}

func (f *Fan) split() {
	n := len(f.outer)
	if n < 2 {
		return
	}

	new_outer := make([]Fan, 0)
	for i := 0; i < n-1; i += 2 {
		p1 := f.outer[i]
		p2 := f.outer[i+1]

		min_radial_ratio, max_radial_ratio := 0.3, 0.7
		min_dist_ratio, max_dist_ratio := 1.0/math.Log2(float64(n))*0.6, 1.0/math.Log2(float64(n))*0.9

		radial_ratio := gm.RandFloat(min_radial_ratio, max_radial_ratio)
		dist_ratio := 1 - gm.RandFloat(min_dist_ratio, max_dist_ratio)

		segment := gm.LineSegment{Begin: p1.root, End: p2.root}
		segment, _ = segment.Split(radial_ratio)

		segment = gm.LineSegment{Begin: f.root, End: segment.End}
		segment, _ = segment.Split(dist_ratio)

		fan := makeFan(segment.End)
		fan.outer = append(fan.outer, p1, p2)
		new_outer = append(new_outer, fan)
	}
	if len(f.outer)%2 != 0 {
		new_outer = append(new_outer, f.outer[len(f.outer)-1])
	}

	f.outer = new_outer

	f.split()
}

func (f *Fan) makeRoads() (roads []Road) {
	for _, subfan := range f.outer {
		var rd Road
		rd.Points = append(rd.Points, f.root, subfan.root)
		roads = append(roads, rd)
		roads = append(roads, subfan.makeRoads()...)
	}
	return
}

func GenerateRoads(cityMap Map, chanMap chan Map, initials InitialValuesRoads) (roads []Road) {
	centers := generateCenters(cityMap, initials)
	roads = append(roads, connectCenters(centers)...)
	exits := generateExits(cityMap, initials)
	roads = append(roads, connectCentersWithExits(centers, exits)...)
	return
}

func generateCenters(cityMap Map, initials InitialValuesRoads) []gm.Point {
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

func connectCenters(centers []gm.Point) (roads []Road) {
	for i := 0; i < len(centers)-1; i++ {
		var rd Road
		rd.Points = append(rd.Points, centers[i])
		for j := i + 1; j < len(centers); j++ {
			rd.Points = append(rd.Points, centers[j])
		}
		roads = append(roads, rd)
	}
	return
}

func generateExits(cityMap Map, initials InitialValuesRoads) []gm.Point {
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
		extend(&sg, cityMap.BorderPoints)
		if ex, _, err := intersectWithFigure(sg, cityMap.BorderPoints); err == nil {
			exit_points = append(exit_points, ex)
		} else {
			println(ex.X, ex.Y)
		}
	}

	return exit_points
}

func connectCentersWithExits(centers []gm.Point, exits []gm.Point) (roads []Road) {
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
