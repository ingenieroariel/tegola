package maths

type Segment [4]float64

func (s Segment) X1() float64 { return s[0] }
func (s Segment) Y1() float64 { return s[1] }
func (s Segment) X2() float64 { return s[2] }
func (s Segment) Y2() float64 { return s[3] }
func (s Segment) Intersect(l Segment) (x, y float64, ok bool) {
	return Intersect([4]float64(s), [4]float64(l))
}
func (s Segment) SlopeIntercept() (m, b float64, slopeDefined bool) {
	return SlopeIntercept([4]float64(s))
}

func (s Segment) ContainsPoint(x, y float64) bool {
	return InBetween(s[0], x, s[2]) == 0 &&
		InBetween(s[1], y, s[3]) == 0
}

func (s Segment) Line() [4]int64 {
	return [4]int64{int64(s[0]), int64(s[1]), int64(s[2]), int64(s[3])}
}

func segments(p []float64) (segs []Segment) {
	if len(p)%2 != 0 {
		panic("number of elements of polygon array should be even.")
	}
	// Last point in the polygon
	lx, ly := p[len(p)-2], p[len(p)-1]
	// Iterating through the pairs of points
	for i, j := 0, 1; i < len(p)-1; i, j = i+2, j+2 {
		segs = append(segs, Segment{lx, ly, p[i], p[j]})
		// Store last point.
		lx, ly = p[i], p[j]
	}
	return segs
}
