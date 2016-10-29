package clip

import "log"

type ptsInt64Type [4]int64

func (pt ptsInt64Type) x1() int64 {
	return pt[0]
}
func (pt ptsInt64Type) x2() int64 {
	return pt[2]
}
func (pt ptsInt64Type) y1() int64 {
	return pt[1]
}
func (pt ptsInt64Type) y2() int64 {
	return pt[3]
}
func (pt ptsInt64Type) float64() (pts [4]float64) {

	for i, v := range pt {
		pts[i] = float64(v)
	}
	return pts
}

type Region [4]float64

func (r Region) MinX() float64 {
	return r[0]
}
func (r Region) MaxX() float64 {
	return r[2]
}
func (r Region) MinY() float64 {
	return r[1]
}
func (r Region) MaxY() float64 {
	return r[3]
}

func (r Region) Contains(x, y float64) bool {
	var goodX, goodY bool
	if r.MaxX() >= r.MinX() {
		goodX = r.MinX() < x && x < r.MaxX()
	} else {
		goodX = r.MaxX() < x && x < r.MinX()

	}
	if r.MaxY() >= r.MinY() {
		goodY = r.MinY() < y && y < r.MaxY()
	} else {
		goodY = r.MaxY() < y && y < r.MinY()
	}
	return goodX && goodY
}

func (r Region) Axises(idx uint8, winding WindingOrder) *Axises {
	cc := func(idxs ...int) *Axises {
		a := &Axises{
			Winding: winding,
			Idx:     int(idx) % 4,
		}
		for i := 0; i < 4; i++ {
			a.Pts[i] = int64(r[idxs[i]])
		}
		return a
	}
	return cc(axisesMap[winding][idx%4]...)
}

func (r Region) ClipCoord(idx int, w WindingOrder) (float64, float64) {
	/*
		        1
		1pt   _____  2pt
		     |     |
		   0 |     | 2
		     |_____|
		0pt     3    3pt

		  Counter Clockwise

		        3
		0pt   _____  3pt
		     |     |
		   0 |     | 2
		     |_____|
		1pt     1    2pt
	*/
	// i should only be in the range of [0,3]
	i := idx % 4

	switch {
	case i == 0 && w.IsClockwise():
		return r.MinX(), r.MaxY()
	case i == 0 && w.IsCounterClockwise():
		return r.MinX(), r.MinY()
	case i == 1 && w.IsClockwise():
		return r.MinX(), r.MinY()
	case i == 1 && w.IsCounterClockwise():
		return r.MinX(), r.MaxY()
	case i == 2 && w.IsClockwise():
		return r.MaxX(), r.MinY()
	case i == 2 && w.IsCounterClockwise():
		return r.MaxX(), r.MaxY()
	case i == 3 && w.IsClockwise():
		return r.MaxX(), r.MaxY()
	case i == 3 && w.IsCounterClockwise():
		return r.MaxX(), r.MinY()
	}
	// Should never get here. This is dead code.
	return 0, 0
}

func (r Region) Clipper(winding WindingOrder) (c []point) {
	c = make([]point, 5, 5)
	for i, _ := range c {
		c[i].X, c[i].Y = r.ClipCoord(i, winding)
		c[i].Type = Clipper
		if i >= 3 {
			c[i].ClipNext = &c[0]
		} else {
			c[i].ClipNext = &c[i+1]
		}
		log.Printf("Clip[%v] x:%3v y:%3v : %v", i, c[i].X, c[i].Y, winding)
	}

	return c
}

func (r Region) DoesCross(pts [4]float64) bool {
	x1, y1 := pts[0], pts[1]
	x2, y2 := pts[2], pts[3]

	x1out := x1 <= r.MinX() || x1 >= r.MaxX()
	x2out := x2 <= r.MinX() || x2 >= r.MaxX()
	y1out := y1 <= r.MinY() || y1 >= r.MaxY()
	y2out := y2 <= r.MinY() || y2 >= r.MaxY()

	// Check diag.
	if x1out && y1out && x2out && y2out {
		//log.Println("diag")
		return true
	}

	// Is pt1/pt2 outside the x min bound and pt2/pt1 outside the x max bound
	if ((x1 <= r.MinX() && x2 >= r.MaxX()) || (x2 <= r.MinX() && x1 >= r.MaxX())) &&
		((y1 > r.MinY() && y1 < r.MaxY()) || (y2 > r.MinY() && y2 < r.MaxY())) {
		//log.Println("X min")
		return true
	}
	// Check y
	if ((y1 <= r.MinY() && y2 >= r.MaxY()) || (y2 <= r.MinY() && y1 >= r.MaxY())) &&
		((x1 > r.MinX() && x1 < r.MaxX()) || (x2 > r.MinX() && x2 < r.MaxX())) {
		//log.Println("Y min")
		return true
	}
	return false
}

/*
FindCrossingsFunc calls the provided function with the each interaction point for each border of the region.
*/
func (r Region) FindCrossingsFunc(winding WindingOrder, pts [4]int64, work func(idx int, x, y int64, inward bool) bool) {
	didCross := false
	for i := 0; i < 4; i++ {
		// log.Printf("Looking at Axises %v of region (%v).", i, r)
		a := r.Axises(uint8(i), winding)
		x, y, ok := a.DoesCross(pts)
		if !ok {
			continue
		} else {
			log.Printf("Pts( $%v ) does cross %v in region( %v ).\n", pts, a, r)
		}
		didCross = true
		if !work(i, x, y, a.IsInwardBound(pts)) {
			return
		}
	}
	if !didCross {
		log.Printf("Pts(%v) did not cross region %v\n", pts, r)
	}
}

func (r Region) FindCrossings(winding WindingOrder, pts [4]int64) (ipts [][2]int64) {

	r.FindCrossingsFunc(winding, pts, func(i int, x, y int64, _ bool) bool {
		ipts = append(ipts, [2]int64{x, y})
		return true
	})
	return ipts
}
func (r Region) FindFirstCrossing(winding WindingOrder, pts [4]int64) (ipts [][2]int64) {

	r.FindCrossingsFunc(winding, pts, func(i int, x, y int64, _ bool) bool {
		ipts = append(ipts, [2]int64{x, y})
		return false
	})
	return ipts
}
