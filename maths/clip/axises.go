package clip

import (
	"fmt"

	"github.com/terranodo/tegola/maths"
)

var axisesMap map[WindingOrder][][]int

func init() {
	axisesMap = make(map[WindingOrder][][]int)

	const (
		x1 = 0
		y1 = 1
		x2 = 2
		y2 = 3
	)

	axisesMap[Clockwise] = [][]int{
		{x1, y2, x1, y1},
		{x1, y1, x2, y1},
		{x2, y1, x2, y2},
		{x2, y2, x1, y2},
	}

	axisesMap[CounterClockwise] = [][]int{
		{x1, y1, x1, y2},
		{x1, y2, x2, y2},
		{x2, y2, x2, y1},
		{x2, y1, x1, y1},
	}
}

type Axises struct {
	Pts     [4]int64
	Winding WindingOrder
	Idx     int
}

func (a *Axises) X1() int64 {
	if a == nil {
		return 0
	}
	return a.Pts[0]
}

func (a *Axises) Y1() int64 {
	if a == nil {
		return 0
	}
	return a.Pts[1]
}
func (a *Axises) X2() int64 {
	if a == nil {
		return 0
	}
	return a.Pts[2]
}
func (a *Axises) Y2() int64 {
	if a == nil {
		return 0
	}
	return a.Pts[3]
}

func (a *Axises) String() string {
	h := a.Pts[2]-a.Pts[0] != 0
	v := a.Pts[3]-a.Pts[1] != 0
	dir := "D"
	if h && !v {
		dir = "H"
	}
	if v && !h {
		dir = "V"
	}
	return fmt.Sprintf("{%v[%3v %3v,%3v %3v] Idx: %1v %v }", dir, a.Pts[0], a.Pts[1], a.Pts[2], a.Pts[3], a.Idx, a.Winding)
}

func (a *Axises) isInwardBound(idx int, pts [4]int64) bool {
	/*
	      3
	      _
	   0 | | 2
	     |_|
	      1
	*/
	x, y := a.Pts[0], a.Pts[1]
	switch idx {
	case 0:
		return pts[0] <= x && x < pts[2]
	case 2:
		return pts[2] < x && x <= pts[0]
	case 1:
		return pts[3] < y && y <= pts[1]
	default:
		return pts[1] <= y && y < pts[3]
	}
}

func (a *Axises) isInwardBoundCC(pts [4]int64) bool {
	/*
	      3
	      _
	   0 | | 2
	     |_|
	      1
	*/
	return a.isInwardBound(a.Idx, pts)
}

func (a *Axises) isInwardBoundC(pts [4]int64) bool {
	/*
	      1
	      _
	   0 | | 2
	     |_|
	      3
	*/
	switch a.Idx {
	case 0:
		return a.isInwardBound(0, pts)
	case 1:
		return a.isInwardBound(3, pts)
	case 2:
		return a.isInwardBound(2, pts)
	default:
		return a.isInwardBound(1, pts)
	}
}

func (a *Axises) IsInwardBound(points [4]int64) bool {
	if a.Winding == Clockwise {
		return a.isInwardBoundC(points)
	}
	return a.isInwardBoundCC(points)
}
func (a *Axises) DoesCross(points [4]int64) (x, y int64, ok bool) {
	pts := ptsInt64Type(points)
	if x, y, ok = maths.IntersectInt64(a.Pts, pts); !ok {
		//log.Println("We did not find an intersection point.")
		// If there isn't an intersection point we, know it does not cross.
		return x, y, ok
	}

	/*
		if (x == pts[0] && y == pts[1]) || (x == pts[2] && y == pts[3]) {
			return x, y, false
		}
	*/
	// If there is one, then we need to make sure it's witin the paramaters of
	// the axises. The nice thing is that since this is a rectangle. We know the
	// axises is going to be eighter verticle or horizontal. So, we only need to make
	// sure the point in on the line.
	if a.X1() == a.X2() {

		// We need to see if the intersection is between the bounds of the y-axises.
		//log.Println("Checking Y-Bounds:", a.Y1(), y, a.Y2(), "X for pts:", pts[0], x, pts[2], "X:", a.X1())
		return x, y, (maths.IsExclusiveInbetweenInt64(a.Y1(), y, a.Y2()) && maths.IsInclusiveInbetweenInt64(pts[0], x, pts[2]))
	}
	//log.Println("Checking X-Bounds:", a.X1(), x, a.X2(), "Y for pts:", pts[1], y, pts[3], "Y", a.Y1())
	return x, y, (maths.IsExclusiveInbetweenInt64(a.X1(), x, a.X2()) && maths.IsInclusiveInbetweenInt64(pts[1], y, pts[3]))
}
