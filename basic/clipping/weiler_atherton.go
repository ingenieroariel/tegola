package clipping

import "fmt"

type PointType uint8

const (
	Clipper = PointType(iota)
	Subject
	Intersect
)

/*
Basics of the alogrithim.

Given:

Clipping polygon

Subject polygon

Result:

One or more polygons clipped into the clipping polygon.


for each vertex for the clipping and subject polygon create a link list.

Say you have the following:

                  k——m
        β---------|--|-ℽ
        |         |  | |
 a——————|———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——|———c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----|------|-δ
 |   e————————f      |
 |                   |
 p———————————————————n

 We will create the following linked lists:

    a → b → c → d → e → f → g → h → k → m → n → p →  a
    α → β → ℽ → δ → α


Now, we will iterate from through the vertices of the subject polygon (a to b, etc…) look for point of intersection with the
clipping polygon (α,β,ℽ,δ). When we come upon an intersection, we will insert an intersection point into both lists.

For example, examing vertex a, and b; against the line formed by (α,β). We notice we have an intersection at I.

                  k——m
        β---------|--|-ℽ
        |         |  | |
 a——————I———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——————c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----|------|-δ
 |   e————————f      |
 |                   |
 p———————————————————n

 We will add I to both lists. We will also note that I, in heading into the clipping region.
 (We will also, mark a as being outside the clipping region, and b being inside the clipping region. If the point is on the boarder of the clipping polygon, it is considered outside of the clipping region.)

    a → I → b → c → d → e → f → g → h → k → m → n → p → a
    α → I → β → ℽ → δ → α

    We will also keep track of the intersections, and weather they are inbound or outbound.
    I(i)

We will check (a,b) against the line formed by (β,ℽ). And see there isn't an intersection.
We will check (a,b) against the line formed by (ℽ,δ). And see there isn't an intersection.
We will check (a,b) against the line formed by (δ,α). And see there isn't an intersection.

When we look at (b,c) we notice that they are both inside the clipping region. And move on to the next set of vertices.

We looking at (c,d), we notice that c is inside and d is outside. This means that there is an intersection point head out.
We check against the line formed by (α,β), and add a Point J, after checking to see we don't already have another equi to J; and adjust
the pointers accordingly. The point c in the subject will now point to J, and J will point to d. And for the intersecting line, α will now point
to J, and J will point to I, as that is what α was pointing to. Our lists will now look like the following.

    a → I → b → c → J → d → e → f → g → h → k → m → n → p → a
    α → J → I → β → ℽ → δ → α

    I(i), J(o)

We will check (c,d) against the line formed by (β,ℽ). And see there isn't an intersection.
We will check (c,d) against the line formed by (ℽ,δ). And see there isn't an intersection.
We will check (c,d) against the line formed by (δ,α). We see there is an intersection, but it is outside of the clipping area.

Next we look at (d,e), notice they are both outside the clipping area, and don't cross through the clipping aread. Thus we can ignore
the points.

                  k——m
        β---------|--|-ℽ
        |         |  | |
 a——————I———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——J———c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----|------|-δ
 |   e————————f      |
 |                   |
 p———————————————————n


Next we look at (e,f), and just (d,e) we can ignore the points as they lie outside and don't cross the clipping area.

Now we look at (f,g), we notice that f is outside, and g is inside the clipping area. This means that The intersection is entering into the clipping area.

We will check (f,g) against the line formed by (α,β). And see there isn't an intersection.
We will check (f,g) against the line formed by (β,ℽ). And see there isn't an intersection.
We will check (f,g) against the line formed by (ℽ,δ). And see there isn't an intersection.

                  k——m
        β---------|--|-ℽ
        |         |  | |
 a——————I———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——J———c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----K------|-δ
 |   e————————f      |
 |                   |
 p———————————————————n

We will check (f,g) against the line formed by (δ,α). We see there is an intersection, and from the previous statement, we know it's an intersection point that is heading inwards. We adjust the link lists to include the point.

    a → I → b → c → J → d → e → f → K → g → h → k → m → n → p → a
    α → J → I → β → ℽ → δ → K → α

    I(i), J(o), K(i)

Looking at (g,h) we realize they are both in the clipping area, and can ignore them.
Next we look at (h,k), here we see that h is inside and k is outside. This means that the intersection point will be outbound.

We will check (h,k) against the line formed by (α,β). And see there isn't an intersection.

                  k——m
        β---------L--|-ℽ
        |         |  | |
 a——————I———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——J———c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----K------|-δ
 |   e————————f      |
 |                   |
 p———————————————————n

We will check (h,k) against the line formed by (β,ℽ). We see there is an intersection (L); also, note that we can stop look at the points, as we found the intersection. We adjust the link lists to include the point.

    a → I → b → c → J → d → e → f → K → g → h → L → k → m → n → p → a
    α → J → I → β → L → ℽ → δ → K → α

    I(i), J(o), K(i), L(o),


Next we look at (k,m) and notice they are not crossing the clipping area and are both outside. So, we know we can skip them.

Looking at (m,n); we notice they are both outside, but are crossing the clipping area, which means there will be two intersection points.

We will check (f,g) against the line formed by (α,β). And see there isn't an intersection.

                  k——m
        β---------L--M-ℽ
        |         |  | |
 a——————I———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——J———c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----K------|-δ
 |   e————————f      |
 |                   |
 p———————————————————n


We will check (f,g) against the line formed by (β,ℽ). We find our first intersection point. We go ahead and insert point (M), as we have done for the other points. We know it's bound as it's the first point in the crossing.
We, adjust, the point we are comparing against from (m,n) to (M,n). Also, note we need to place the point in the correct position between β and ℽ, after L.

    a → I → b → c → J → d → e → f → K → g → h → L → k → m → M → n → p → a
    α → J → I → β → L → M → ℽ → δ → K → α

    I(i), J(o), K(i), L(o), M(i)


We will check (M,n) against the line formed by (ℽ,δ). And see there isn't an intersection.

                  k——m
        β---------L--M-ℽ
        |         |  | |
 a——————I———b     |  | |
 |      |   |     |  | |
 |      |   |     |  | |
 |   d——J———c g———h  | |
 |   |  |     |      | |
 |   |  |     |      | |
 |   |  α-----K------N-δ
 |   e————————f      |
 |                   |
 p———————————————————n

We will check (M,n) against the line formed by (δ,α). We see there is an intersection, and from the previous statement, we know it's an intersection point that is heading outwards. We adjust the link lists to include the point.

    a → I → b → c → J → d → e → f → K → g → h → L → k → m → M → N → n → p → a
    α → J → I → β → L → M → ℽ → δ → N → K → α

    I(i), J(o), K(i), L(o), M(i), N(o)

Next we look at (n,p), and know we can skip the points as they are both outside, and not crossing the clipping area.
Finally we look at (p,a), and again because they are both outside, and not corssing the clipping area, we know we can skip them.

Finally we check to see if we have at least one external and one internal point. If we don't have any external points, we know the polygon is contained within the clipping area and can just return it.
If we don't have any internal points, and no Intersections points, we know the polygon is contained compleatly outside and we can return any empty array of polygons.

First thing we do is iterate our list of Intersection points looking for the first point that is an inward bound point. I is such a point. The rule is if an intersection point is inward, we following the subject links, if it's outward we follow, the clipping links.
Since I is inward, we write it down, and follow the subject link to b.

LineString1: I,b

Then we follow the links till we get to the next Intersection point.

LineString1: I,b,c,J

        •··············•
        ·              ·
        +———+          ·
        |   |          ·
        |   |          ·
        +———+          ·
        ·              ·
        ·              ·
        •··············•



Since, J is outward we follow the clipping links, which leads us to I. Since I is also the first point in this line string. We know we are done, with the first clipped polygon.

Next we iterate to the next inward Intersection point from J, to K.
LineString1: I,b,c,J
LineString2: K

And as before since K is inward point we follow the subject polygon points, till we get to an intersection point.

LineString1: I,b,c,J
LineString2: K,g,h,L

As L is an outward intersection point we follow the clipping polygon points, till we get to an intersection point.

LineString1: I,b,c,J
LineString2: K,g,h,L,M

As M is an inward intersection point we follow the subject.

LineString1: I,b,c,J
LineString2: K,g,h,L,M,N

As N is an outward intersection point we follow the clipping, and discover that the point is our starting Intersection point K. That ends is our second clipped polygon.

LineString1: I,b,c,J
LineString2: K,g,h,L,M,N


        •·········+--+·•
        ·         |  | ·
        +———+     |  | ·
        |   |     |  | ·
        |   |     |  | ·
        +———+ +———+  | ·
        ·     |      | ·
        ·     |      | ·
        •·····+------+·•


Since N(o), is the end of the array we, start at the beginning and notice, that we already accounted for I(i). And so, we are done.

*/

type point struct {
	// The X coord of the Point
	X float64
	// The Y coord of the Point
	Y float64
	// The Point type
	Type PointType
	// IsIn is overloaded depending on the PointType.
	// If the point type is a clipper point it means nothing.
	// if the point type is a Subject point it indicates weather the point is in the clipped region.
	// If the point type is an Intersect point it indicates weather the point indicates the the polygon is entering or leaving the clip reagion. True for entring, false for leaving. (Inwards v.s. Outwards.)
	IsIn bool
	// The Next for the subject.
	SubNext *point
	// the Prev point for the subject.
	SubPrev *point
	// The Next point for the clip
	ClipNext *point
	// The Prev point for the clip
	ClipPrev *point
	IntNext  *point
	IntPrev  *point
	seen     bool
}

func (pt *point) Next() *point {
	if pt == nil {
		return nil
	}
	switch pt.Type {
	case Clipper:
		return pt.ClipNext
	case Subject:
		return pt.SubNext
	case Intersect:
		if pt.IsIn {
			return pt.SubNext
		}
		return pt.ClipNext
	}
	panic("Should not get here.")
}
func (pt *point) Prev() *point {
	if pt == nil {
		return nil
	}
	switch pt.Type {
	case Clipper:
		return pt.ClipPrev
	case Subject:
		return pt.SubPrev
	case Intersect:
		if pt.IsIn {
			return pt.ClipPrev
		}
		return pt.SubPrev
	}
	panic("Should not get here.")
}

/*
func (pt *point) String() string {
	if pt == nil {
		return "Point( nil )"
	}
	switch pt.Type {
	case Clipper:
		return fmt.Sprintf("Clipper( %v %v )", pt.X, pt.Y)
	case Subject:
		return fmt.Sprintf("Subject( %v %v )", pt.X, pt.Y)
	case Intersect:
		return fmt.Sprintf("Intersect( %v %v )", pt.X, pt.Y)
	default:
		return fmt.Sprintf("Point( %v %v )", pt.X, pt.Y)
	}
}
*/
func (pt point) String() string {
	switch pt.Type {
	case Clipper:
		return fmt.Sprintf("Clipper( %v %v )", pt.X, pt.Y)
	case Subject:
		return fmt.Sprintf("Subject( %v %v )", pt.X, pt.Y)
	case Intersect:
		ms := "<-"
		if pt.IsIn {
			ms = "->"
		}
		return fmt.Sprintf("Intersect%v( %v %v )", ms, pt.X, pt.Y)
	default:
		return fmt.Sprintf("Point( %v %v )", pt.X, pt.Y)
	}
}

func (pt *point) NextInboundIntersect() *point {
	if pt == nil {
		return nil
	}
	nxt := pt.IntNext
	if nxt == nil {
		return nil
	}
	for nxt.seen || !nxt.IsIn {
		nxt = nxt.IntNext
		if nxt == nil {
			return nil
		}
		if nxt == pt {
			return nil
		}
	}
	return nxt
}

func (p *point) ResetSeen() {
	if p == nil {
		return
	}
	p.seen = false
	nxt := p.IntNext
	for nxt != p {
		if nxt == nil {
			return
		}
		nxt.seen = false
		nxt = nxt.IntNext
	}
}
func (p *point) WalkTree() {
	if p == nil {
		return
	}
	//log.Println(p)
	nxt := p.IntNext
	for nxt != p {
		if nxt == nil {
			return
		}
		//log.Println(nxt)
		nxt = nxt.IntNext
	}
}

type ClippingRegion [4]float64

func (c ClippingRegion) MinX() float64 {
	return c[0]
}
func (c ClippingRegion) MaxX() float64 {
	return c[2]
}
func (c ClippingRegion) MinY() float64 {
	return c[1]
}
func (c ClippingRegion) MaxY() float64 {
	return c[3]
}

func (c ClippingRegion) ContainsPt(x, y float64) bool {
	return c[0] < x && x < c[2] && c[1] < y && y < c[3]
	// return Contains([2]float64{x, y}, []float64{c[0], c[1], c[2], c[3]})
}

func (c ClippingRegion) Axises(idx int, winding WindingOrder) (a [4]float64) {
	cc := func(idx ...int) [4]float64 {
		for i := 0; i < 4; i++ {
			a[i] = c[idx[i]]
		}
		return a
	}
	switch idx {
	case 0:
		if winding == Clockwise {
			return cc(0, 3, 0, 1)
		}
		return cc(2, 1, 0, 1)
	case 1:
		if winding == Clockwise {
			return cc(0, 1, 2, 1)
		}
		return cc(0, 1, 0, 3)

	case 2:
		if winding == Clockwise {
			return cc(2, 1, 2, 3)
		}
		return cc(0, 3, 2, 3)
	case 3:
		if winding == Clockwise {
			return cc(2, 3, 0, 3)
		}
		return cc(2, 3, 2, 1)
	default:
		panic("Should not get here!")
	}

	return a
}

func (c ClippingRegion) Inward(idx int, winding WindingOrder, endpt *point) bool {
	i := idx % 4
	switch i {
	case 0:
		if winding == Clockwise {
			return endpt.X > c.MinX()
		}
		return endpt.Y > c.MinY()
	case 1:
		if winding == Clockwise {
			return endpt.Y < c.MinY()
		}
		return endpt.X > c.MinX()
	case 2:
		if winding == Clockwise {
			return endpt.X < c.MaxX()
		}
		return endpt.Y < c.MaxY()
	case 3:
		if winding == Clockwise {
			return endpt.Y < c.MaxY()
		}
		return endpt.X < c.MaxX()
	}
	return false
}

type WindingOrder bool

const (
	Clockwise        = WindingOrder(false)
	CounterClockwise = WindingOrder(true)
)

func (w WindingOrder) String() string {
	if w {
		return "Counter Clockwise"
	}
	return "Clockwise"
}

func inClip(c ClippingRegion, x, y float64) bool {
	var goodX, goodY bool
	if c.MaxX() >= c.MinX() {
		goodX = c.MinX() < x && x < c.MaxX()

	} else {
		goodX = c.MaxX() < x && x < c.MinX()
	}
	if c.MaxY() >= c.MinY() {
		goodY = c.MinY() < y && y < c.MaxY()

	} else {
		goodY = c.MaxY() < y && y < c.MinY()
	}

	return goodX && goodY

}

/*
func qudrantForGivenRegion(r ClippingRegion, x, y float64) int {
	switch {
	case x <= c.MinX() && y >= c.MinY():
		return 0
	case x > c.MinX() && x < x.MaxX() && y >= c.MinY():
		return 1
	case x >= c.MxxX() && y >= c.MinY():
		return 2
	case y > c.MaxY() && y < x.MinY() && x <= c.MinX():
		return 3
	case x > c.MinX() && x < x.MaxX() && y > c.MaxY() && y < x.MinY():
		return 4
	case y > c.MaxY() && y < x.MinY() && x >= c.MaxX():
		return 5
	}

}
*/

func doesCrossClip(c ClippingRegion, pts [4]float64) bool {
	x1, y1 := pts[0], pts[1]
	x2, y2 := pts[2], pts[3]

	x1out := x1 <= c.MinX() || x1 >= c.MaxX()
	x2out := x2 <= c.MinX() || x2 >= c.MaxX()
	y1out := y1 <= c.MinY() || y1 >= c.MaxY()
	y2out := y2 <= c.MinY() || y2 >= c.MaxY()

	// Check diag.
	if x1out && y1out && x2out && y2out {
		//log.Println("diag")
		return true
	}

	// Is pt1/pt2 outside the x min bound and pt2/pt1 outside the x max bound
	if ((x1 <= c.MinX() && x2 >= c.MaxX()) || (x2 <= c.MinX() && x1 >= c.MaxX())) &&
		((y1 > c.MinY() && y1 < c.MaxY()) || (y2 > c.MinY() && y2 < c.MaxY())) {
		//log.Println("X min")
		return true
	}
	// Check y
	if ((y1 <= c.MinY() && y2 >= c.MaxY()) || (y2 <= c.MinY() && y1 >= c.MaxY())) &&
		((x1 > c.MinX() && x1 < c.MaxX()) || (x2 > c.MinX() && x2 < c.MaxX())) {
		//log.Println("Y min")
		return true
	}
	return false
}

func slopeIntercept(l [4]float64) (m, b float64, slopeDefined bool) {
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
func intersect(line1 [4]float64, line2 [4]float64) (x, y float64, ok bool) {

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

	m1, b1, definedSlope1 := slopeIntercept(line1)
	m2, b2, definedSlope2 := slopeIntercept(line2)

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

func doesCrossAxisFloat(c ClippingRegion, idx int, winding WindingOrder, pts [4]float64) (x, y float64, ok bool) {

	axises := c.Axises(idx, winding)
	x1, y1 := int64(axises[0]), int64(axises[1])
	x2, y2 := int64(axises[2]), int64(axises[3])
	px1, py1 := int64(pts[0]), int64(pts[1])
	px2, py2 := int64(pts[2]), int64(pts[3])
	x, y, ok = intersect(axises, pts)
	ix, iy := int64(x), int64(y)
	x, y = float64(ix), float64(iy)

	// log.Printf("For Points ( %v %v , %v %v ) intersect with axis(%v -- %+v) x %v y %v ok: %v", px1, py1, px2, py2, idx, axises, x, y, ok)
	if !ok {
		return x, y, ok
	}
	if x1 <= x2 {
		ok = x1 <= ix && ix <= x2
	} else {
		//log.Println("x2 less then x1", x2, ix, x1, x2 <= ix, ix <= x1)
		ok = x2 <= ix && ix <= x1
	}
	if !ok {
		//log.Println("After x redo ok", idx, x, y, ok)
		return x, y, ok
	}

	if y1 <= y2 {
		ok = y1 <= iy && iy <= y2
	} else {
		//log.Println("y2 less then y1")
		ok = y2 <= iy && iy <= y1
	}
	if !ok {
		//log.Println("After y redo ok", idx, x, y, ok)
		return x, y, ok
	}

	// Now we need to make sure it's between the points as well.

	if px1 <= px2 {
		ok = px1 <= ix && ix <= px2
	} else {
		//log.Println("x2 less then x1", px2, ix, px1, px2 <= ix, ix <= px1)
		ok = px2 <= ix && ix <= px1
	}
	if !ok {
		//log.Println("After px redo ok", idx, x, y, ok)
		return x, y, ok
	}

	if py1 <= py2 {
		ok = py1 <= iy && iy <= py2
	} else {
		//log.Println("y2 less then y1")
		ok = py2 <= iy && iy <= py1
	}

	//log.Println("After py redo ok", idx, x, y, ok)
	return x, y, ok

}
func doesCrossAxis(c ClippingRegion, idx int, winding WindingOrder, pt1, pt2 *point) (x, y float64, ok bool) {
	return doesCrossAxisFloat(c, idx, winding, [4]float64{pt1.X, pt1.Y, pt2.X, pt2.Y})
}

func clipCoord(c ClippingRegion, winding WindingOrder, idx int) (float64, float64) {
	switch idx {
	case 0:
		if winding == Clockwise {
			return c.MinX(), c.MaxY()
		}
		return c.MaxX(), c.MinY()
	case 1:
		return c.MinX(), c.MinY()
	case 2:
		if winding == Clockwise {
			return c.MaxX(), c.MinY()
		}
		return c.MinX(), c.MaxY()
	case 3:
		return c.MaxX(), c.MaxY()
	default:
		panic("Expected idx to be between 0,3 inclusive.")
	}
}

//
func findListPos(start *point, cmp func(pt *point) bool, next func(pt *point) *point) (startpt, endpt *point) {
	startpt = start
	endpt = next(startpt)
	for {
		if endpt == start {
			panic("Went around the entire lists and did not find a pos.")
		}
		if cmp(endpt) {
			// log.Println("Found point", endpt)
			return startpt, endpt
		}
		startpt = endpt
		//endpt = startpt.ClipNext
		endpt = next(startpt)
	}
}
func insertIntoList(start, end, item *point, next func(pt *point) *point, finilize func(start, end, item *point)) {

	startpt := start
	endpt := end
	x, y := item.X, item.Y
	// First I need to figure out which is changing. We will favor change in X.
	if start.Y == end.Y {
		//log.Println("Increasing in x")
		// need to find insertion point in x.
		if start.X < end.X {
			// increasing.
			if item.X < start.X {
				panic(fmt.Sprintf("Increasing in X Item is before start! %t, %t", start, item))
			}
			if item.X > end.X {
				panic(fmt.Sprintf("Increasing in X Item is after end! %t, %t", start, item))
			}
			startpt, endpt = findListPos(start, func(pt *point) bool {
				return x <= pt.X
			}, next)
		} else {
			//log.Println("Increasing in x")
			if item.X > start.X {
				panic(fmt.Sprintf("Decreasing in x Item is before start! %t, %t", start, item))
			}
			if item.X < end.X {
				panic(fmt.Sprintf("Decreasing in x Item is after end! %t, %t", start, item))
			}
			startpt, endpt = findListPos(start, func(pt *point) bool {
				return x >= pt.X
			}, next)
		}
	} else {
		if start.Y < end.Y {
			// increasing.
			//log.Println("Increasing in y")
			if item.Y < start.Y {
				panic(fmt.Sprintf("Increasing in y Item is before start! %t, %t", start, item))
			}
			if item.Y > end.Y {
				panic(fmt.Sprintf("Increasing in y Item is after end! %t, %t", start, item))
			}
			startpt, endpt = findListPos(start, func(pt *point) bool {
				return y <= pt.Y
			}, next)
		} else {
			//log.Println("Decreasing in y")
			if item.Y > start.Y {
				panic(fmt.Sprintf("decreasing in y Item is before start! %t, %t", start, item))
			}
			if item.Y < end.Y {
				panic(fmt.Sprintf("decreasing in y Item is after end! %t, %t", start, item))
			}
			startpt, endpt = findListPos(start, func(pt *point) bool {
				return y >= pt.Y
			}, next)
		}
	}
	finilize(startpt, endpt, item)

}
func insertIntoClipList(start, end, item *point) {
	insertIntoList(start, end, item, func(pt *point) *point {
		return pt.ClipNext
	}, func(s, e, i *point) {
		i.ClipPrev = s
		i.ClipNext = e
		s.ClipNext = i
		e.ClipPrev = i
	})
}
func insertIntoSubList(start, end, item *point) {
	//log.Println("Inserting into sub", start, end, item)
	insertIntoList(start, end, item, func(pt *point) *point {
		//log.Println("Next", pt.SubNext)
		return pt.SubNext
	}, func(s, e, i *point) {
		i.SubPrev = s
		i.SubNext = e
		s.SubNext = i
		e.SubPrev = i
	})
}

// Contains will tell you if the point is contained inside of the region. It does this by drawing
// a line from the pt to the extent_x, and counting the number of intersection points. If the the
// number is odd it's contained, if it's even it's not contained.
func Contains(pt [2]float64, region []float64) bool {
	numpts := len(region) / 2
	s := numpts - 1
	ptx := region[s*2]
	pty := region[(s*2)+1]
	count := 0
	for i := 0; i < numpts; i++ {
		lptx := ptx
		lpty := pty
		ptx = region[i*2]
		pty = region[(i*2)+1]
		// If both x values are greater then our points x value, it's not going to intersect.
		if ptx > pt[0] && lptx > pt[0] {
			continue
		}
		// Skip if they don't cross the y-axsis
		if pty > pt[1] && lpty > pt[1] {
			continue
		}
		/*
			if (pty > pt[1] && lpty > pt[1]) ||
				(pty < pt[1] && lpty < pt[1]) {
				continue
			}
		*/
		count++
	}
	return count%2 != 0
}

func ClipPolygon(clipping ClippingRegion, winding WindingOrder, subject []float64) [][]float64 {
	clippedSubjects := make([][]float64, 0, 0)
	if len(subject) <= 6 {
		return clippedSubjects
	}
	clipper := make([]point, 4, 4)
	for i, _ := range clipper {
		clipper[i].X, clipper[i].Y = clipCoord(clipping, winding, i)
		clipper[i].Type = Clipper
		if i == (len(clipper) - 1) {
			clipper[i].ClipNext = &clipper[0]
		} else {
			clipper[i].ClipNext = &clipper[i+1]
		}
		if i == 0 {
			clipper[i].ClipPrev = &clipper[len(clipper)-1]
		} else {
			clipper[i].ClipPrev = &clipper[i-1]
		}
	}

	sublen := len(subject) / 2
	subStart := &point{
		X:    subject[0],
		Y:    subject[1],
		IsIn: inClip(clipping, subject[0], subject[1]),
		Type: Subject,
		SubNext: &point{
			X:    subject[len(subject)-2],
			Y:    subject[len(subject)-1],
			IsIn: inClip(clipping, subject[len(subject)-2], subject[len(subject)-1]),
			Type: Subject,
		},
	}
	subStart.SubPrev = subStart.SubNext
	subStart.SubNext.SubNext = subStart
	subStart.SubNext.SubPrev = subStart
	subHead := subStart
	var intHead *point
	var intCurrent *point
	var initialInboundIntersect *point

	allInside := true  // Assume all points are inside, to begin with.
	allOutside := true // Pardoxitly we will also assume all the points are outside. :)

SUBJECTFOR:
	for i := 1; i <= sublen; i++ {
		var nxtPt *point
		var inside bool
		if i != sublen {
			// log.Printf("Looking at subject vertex %v, %v:( %v %v , %v %v )", i-1, i, subStart.X, subStart.Y, subject[i*2], subject[(i*2)+1])
			inside = inClip(clipping, subject[i*2], subject[(i*2)+1])
			nxtPt = &point{
				X:       subject[i*2],
				Y:       subject[(i*2)+1],
				IsIn:    inside,
				Type:    Subject,
				SubNext: subStart.SubNext,
				SubPrev: subStart,
			}
		} else {
			// log.Printf("Looking at subject vertex %v, %v:( %v %v , %v %v )", i-1, 0, subStart.X, subStart.Y, subHead.X, subHead.Y)
			nxtPt = subHead
			inside = subHead.IsIn
		}
		if allInside && !inside {
			allInside = false
		}
		if allOutside && inside {
			allOutside = false
		}
		subStart.SubNext = nxtPt
		// If both points are in the clipping region, we move on.
		if subStart.IsIn && nxtPt.IsIn {
			// log.Println("Both In...moving on. ", winding)
			subStart = nxtPt
			continue SUBJECTFOR
		}
		// We are entering the clipping area. We need to find one Intersection point.
		if (!subStart.IsIn && nxtPt.IsIn) ||
			(subStart.IsIn && !nxtPt.IsIn) {
			// log.Println("Entering/Exiting clip, one in one out. ", winding)
			for j := 0; j < 4; j++ {
				x, y, ok := doesCrossAxis(clipping, j, winding, subStart, nxtPt)
				// we are done.
				if !ok {
					//	log.Println("Not ok", x, y, ok)
					continue
				}
				//log.Println("Not ok", x, y, ok)
				/*
					if (x == subStart.X && y == subStart.Y) ||
						(x == nxtPt.X && y == nxtPt.Y) {
						continue
					}
				*/
				// log.Println(clipping, "Crosses Axis", j, clipping.Axises(j, winding))
				intPt := &point{
					X:       x,
					Y:       y,
					Type:    Intersect,
					IsIn:    !subStart.IsIn && nxtPt.IsIn, // If this is true we are entering the clipping area, otherwise we are leaving it.
					SubPrev: subStart,
					SubNext: nxtPt,
				}
				insertIntoSubList(subStart, nxtPt, intPt)
				if j == 3 {
					insertIntoClipList(&clipper[j], &clipper[0], intPt)
				} else {
					insertIntoClipList(&clipper[j], &clipper[j+1], intPt)
				}
				// Now we have to update the interptlist.
				if intCurrent == nil {
					intCurrent = intPt
					intHead = intCurrent
				}
				intPt.IntPrev = intCurrent
				intCurrent.IntNext = intPt
				intPt.IntNext = intHead
				intCurrent = intPt
				if initialInboundIntersect == nil && intPt.IsIn {
					initialInboundIntersect = intPt
				}
				subStart = nxtPt
				// log.Println("To Next point")
				continue SUBJECTFOR
			}
			panic(fmt.Sprintf("Line: 1033 Did not find an intersection! This should not happen! %+v %v -- %+v , %+v %v , %+v %v", clipping, winding, clipper, subStart, subStart.IsIn, nxtPt, nxtPt.IsIn))
		}
		// All that's left is both points are outside. In this case we need to check if it crosses the clipping area.
		// log.Println("Crossing the clip. ", winding)
		// There are two points we need to find in this case; one will be entering the clipping area, the other will be
		// exiting the clipping area.
		// It's possible that all the points are outside but cross the clipping zone, as in this case, do we mark it
		// as there be points in the clipping zone as we are going to be adding them.
		//count := 0
		for j := 0; j < 4; j++ {
			x, y, ok := doesCrossAxis(clipping, j, winding, subStart, nxtPt)
			// we are done.
			if !ok {
				continue
			}
			allOutside = false
			inward := clipping.Inward(j, winding, nxtPt)
			intPt := &point{
				X:       x,
				Y:       y,
				Type:    Intersect,
				IsIn:    inward,
				SubPrev: subStart,
				SubNext: nxtPt,
			}
			// log.Println(count, clipping, "Crosses axis ", j, clipping.Axises(j, winding), subStart, nxtPt, intPt, inward)
			insertIntoSubList(subStart, nxtPt, intPt)
			if j == 3 {
				insertIntoClipList(&clipper[j], &clipper[0], intPt)
			} else {
				insertIntoClipList(&clipper[j], &clipper[j+1], intPt)
			}
			// Now we have to update the interptlist.
			if intCurrent == nil {
				intCurrent = intPt
				intHead = intCurrent
			}
			intPt.IntPrev = intCurrent
			intCurrent.IntNext = intPt
			intPt.IntNext = intHead
			intCurrent = intPt
			if initialInboundIntersect == nil && intPt.IsIn {
				initialInboundIntersect = intPt
			}
			//count++
		}

		/*
					if count > 0 {
			//			log.Println("Found additional intersetions points", count)
					}
		*/
		subStart = nxtPt
	}

	if allOutside {
		// All the points are outside, nothing to return.
		return clippedSubjects
	}

	if allInside {
		// All the points are inside, no clipping needs to be done.
		clippedSubjects = append(clippedSubjects, subject)
		return clippedSubjects
	}

	// Now we have the data structure in place we need to walk through the intersections point and get the clipped polygons.
	if initialInboundIntersect == nil {
		panic(fmt.Sprintf("We should have gotten an intersect point that was in bound."))
	}

	inboundInt := initialInboundIntersect
	inboundInt.seen = true
	count := 0
	for {
		nxtPt := inboundInt.Next()
		nxtPt.seen = true
		linestr := []float64{inboundInt.X, inboundInt.Y}
		// log.Println("Walking the points.", inboundInt)
		for nxtPt != inboundInt {
			count++
			// log.Println("NxtPt", nxtPt)
			linestr = append(linestr, nxtPt.X, nxtPt.Y)
			nxtPt = nxtPt.Next()
			nxtPt.seen = true
			if count > 10 {
				return clippedSubjects
			}
		}
		// log.Println("Found clipping", linestr)
		clippedSubjects = append(clippedSubjects, linestr)
		inboundInt = intHead.NextInboundIntersect()
		if inboundInt == nil {
			return clippedSubjects
		}
	}
	return clippedSubjects
}

func ClipLineString(clipping ClippingRegion, subject []float64) (clippedSubjects [][]float64) {
	clippedSubjects = make([][]float64, 0, 0)
	sublen := len(subject) / 2
	sub := make([]float64, 0, 0)
	lin, lx, ly := inClip(clipping, subject[0], subject[1]), subject[0], subject[1]
	for i := 1; i < sublen; i++ {
		//log.Printf("Looking at subject vertex %v, %v:( %v %v , %v %v )", i-1, i, subStart.X, subStart.Y, subject[i*2], subject[(i*2)+1])
		nin, nx, ny := inClip(clipping, subject[i*2], subject[(i*2)+1]), subject[i*2], subject[(i*2)+1]
		if !lin && !nin {
			if doesCrossClip(clipping, [4]float64{lx, ly, nx, ny}) {
				if len(sub) > 0 {
					clippedSubjects = append(clippedSubjects, sub)
					sub = make([]float64, 0, 0)
				}

				for j := 0; j < 4; j++ {
					x, y, ok := doesCrossAxisFloat(clipping, j, Clockwise, [4]float64{lx, ly, nx, ny})
					// we are done.
					if !ok {
						continue
					}

					//log.Println("Adding intersect point ", x, y)
					sub = append(sub, x, y)
				}
				clippedSubjects = append(clippedSubjects, sub)
				sub = make([]float64, 0, 0)

			}
			//log.Println("Skipping point ", lx, ly)
			lx, ly, lin = nx, ny, nin
			continue
		}
		if lin && nin {
			sub = append(sub, lx, ly)
			//log.Println("Adding point ", lx, ly)
			lx, ly, lin = nx, ny, nin
			continue
		}
		if !lin && nin {
			for j := 0; j < 4; j++ {
				x, y, ok := doesCrossAxisFloat(clipping, j, Clockwise, [4]float64{lx, ly, nx, ny})
				// we are done.
				if !ok {
					continue
				}

				//log.Println("Adding intersect point ", x, y)
				sub = append(sub, x, y)
				break
			}
			//log.Println("Skipping point ", lx, ly)
			lx, ly, lin = nx, ny, nin
			continue
		}
		if lin && !nin {
			//log.Println("Adding point ", lx, ly)
			sub = append(sub, lx, ly)
			for j := 0; j < 4; j++ {
				x, y, ok := doesCrossAxisFloat(clipping, j, Clockwise, [4]float64{lx, ly, nx, ny})
				// we are done.
				if !ok {
					continue
				}
				//log.Println("Adding intersect point ", x, y)
				sub = append(sub, x, y)
				clippedSubjects = append(clippedSubjects, sub)
				//log.Println("New line:", sub)
				sub = make([]float64, 0, 0)
				break
			}
			lx, ly, lin = nx, ny, nin
			continue
		}
	}

	if len(sub) > 0 {
		//log.Println("New line:", sub)
		clippedSubjects = append(clippedSubjects, sub)
	}
	return clippedSubjects

}
