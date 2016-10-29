/*
Package math contins generic math functions that we need for doing transforms.
this package will augment the go math library.
*/
package maths

import (
	"log"
	"math"

	"github.com/terranodo/tegola"
)

const (
	WebMercator = tegola.WebMercator
	WGS84       = tegola.WGS84
	Deg2Rad     = math.Pi / 180
	Rad2Deg     = 180 / math.Pi
	PiDiv2      = math.Pi / 2.0
	PiDiv4      = math.Pi / 4.0
)

// AreaOfPolygon will calculate the Area of a polygon using the surveyor's formula
// (https://en.wikipedia.org/wiki/Shoelace_formula)
func AreaOfPolygon(p tegola.Polygon) (area float64) {
	var points []tegola.Point
	for _, l := range p.Sublines() {
		points = append(points, l.Subpoints()...)
	}
	n := len(points)
	for i := range points {
		j := (i + 1) % n
		area += points[i].X() * points[j].Y()
		area -= points[j].X() * points[i].Y()
	}
	return math.Abs(area) / 2.0
}

func RadToDeg(rad float64) float64 {
	return rad * Rad2Deg
}

func DegToRad(deg float64) float64 {
	return deg * Deg2Rad
}

func SlopeIntercept(l [4]float64) (m, b float64, slopeDefined bool) {
	dx := l[2] - l[0]
	dy := l[3] - l[1]
	//log.Println("dx", dx, "dy", dy)
	if dx == 0 || dy == 0 {
		return 0, l[1], dx != 0
	}
	m = dy / dx

	// y - y1 = m(x - x1)
	// y = m(x -x1) + y1
	// y = mx - m(x1) + y1
	// y = mx + y1 - m(x1)
	// y = mx + b
	b = l[1] - (m * l[0])
	return m, b, true
}

// First we need to get it into point slop form which means we need m and b.
func Intersect(line1 [4]float64, line2 [4]float64) (x, y float64, ok bool) {

	if line1[0] == line1[2] {
		if line1[0] == line2[0] {
			return line1[0], line2[1], true
		}
		if line1[0] == line2[2] {
			return line1[0], line2[3], true
		}
	}
	if line1[1] == line1[3] {
		if line1[1] == line2[1] {
			return line2[0], line1[1], true
		}
		if line1[1] == line2[3] {
			return line2[2], line1[1], true
		}
	}

	m1, b1, definedSlope1 := SlopeIntercept(line1)
	m2, b2, definedSlope2 := SlopeIntercept(line2)

	//log.Println("line1", line1, m1, b1, definedSlope1)
	//log.Println("line2", line2, m2, b2, definedSlope2)

	// If the slopes are the same then they are parallel so, they don't intersect.
	if definedSlope1 == definedSlope2 && m1 == m2 {
		return 0, 0, false
	}

	// line1 is Horizontal We have the value for x, need the value for y.
	if !definedSlope1 {
		// y = m2x+b
		// (y - b)/m2 = x
		x = line1[0]
		if m2 == 0 {
			return x, b2, true
		}
		y = (m2 * x) + b2
		return x, y, true
	}

	if !definedSlope2 {
		x = line2[0]
		y = line1[1]
		if m1 == 0 {
			return x, b1, true
		}
		y = (m1 * x) + b1
		return x, y, true
	}

	if m1 == 0 {
		// y = mx+b
		// (y - b)/m = x
		y = line1[1]
		x = (y - b2) / m2
		return x, y, true
	}

	if m2 == 0 {
		y = line2[1]
		x = (y - b1) / m1
		return x, y, true
	}

	// y = m1x+b1
	// y = m2x+b2
	// m1x+b1 = m2x+b2
	// m1x-m2x = b2-b1
	// x(m1-m2) = b2-b1
	// x = (b2-b1)/(m1-m2)
	dm := m1 - m2
	db := (b2 - b1)
	x = db / dm
	y = (m1 * x) + b1
	return x, y, true
}

func IntersectInt64(line1 [4]int64, line2 [4]int64) (x, y int64, ok bool) {
	l1 := [4]float64{float64(line1[0]), float64(line1[1]), float64(line1[2]), float64(line1[3])}
	l2 := [4]float64{float64(line2[0]), float64(line2[1]), float64(line2[2]), float64(line2[3])}
	fx, fy, ok := Intersect(l1, l2)
	return int64(fx), int64(fy), ok
}

func IsInclusiveInbetweenInt64(v1, v, v2 int64) bool {
	if v1 == v2 && v1 == v {
		return true
	}
	if v1 < v2 {
		return v1 <= v && v <= v2
	}
	return v2 <= v && v <= v1
}
func IsExclusiveInbetweenInt64(v1, v, v2 int64) bool {
	if v1 == v2 && v1 == v {
		return false
	}
	if v1 < v2 {
		return v1 < v && v < v2
	}
	return v2 < v && v < v1
}

func InBetween(v1, v, v2 float64) int {
	if v1 == v || v2 == v {
		return 0
	}

	if v1 > v2 {
		if v2 > v {
			return 1
		}
		return -1
	}
	if v1 > v {
		return -1
	}
	return 1
}

// Contains will tell you if the point is contained inside of the region. It does this by drawing
// a line from the pt to the extent_x, and counting the number of intersection points. If the the
// number is odd it's contained, if it's even it's not contained.
func Contains(p []float64, x, y float64) bool {
	log.Printf("Checking to see if %v contains (%v,%v)", p, x, y)
	count := 0
	for _, seg := range segments(p) {
		// We treat being on a segment as being in the region.
		if seg.ContainsPoint(x, y) {
			log.Printf("Contained in the border.")
			return true
		}
		// If it's above or blow the point, don't count it.
		// If both x values are bigger then the point don't count it.
		minx := seg.X1()
		if seg.X2() < seg.X1() {
			minx = seg.X2()
		}

		if minx > x {
			continue
		}

		ix, iy, ok := IntersectInt64(seg.Line(), [4]int64{int64(minx), int64(y), int64(x), int64(y)})
		if ok &&
			IsInclusiveInbetweenInt64(int64(seg.X1()), ix, int64(seg.X2())) &&
			IsInclusiveInbetweenInt64(int64(seg.Y1()), iy, int64(seg.Y2())) {
			log.Println("Hit segment", seg, x, y, minx, y)
			count++
		}
		/*
			log.Println(InBetween(seg.Y1(), y, seg.Y2()), "y", seg.Y1(), seg.Y2(), "pt", x, y)
			if (InBetween(seg.Y1(), y, seg.Y2()) != 0) ||
				(seg.X1() > x && seg.X2() > x) {
				continue
			}
			count++
		*/
	}
	log.Printf("Count is %v", count)
	// an odd count means it's in.
	return count%2 != 0
}
