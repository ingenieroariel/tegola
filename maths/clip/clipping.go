package clip

import (
	"log"

	"github.com/terranodo/tegola/maths"
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
// dedupSub will remove duplicated paired coordinates. This will only removed coordinates that repeat.
func dedupSub(s []float64) (sub []float64) {
	if len(s) < 4 {
		return s
	}
	for x, y := 0, 1; x < len(s)-2; x, y = x+2, y+2 {
		log.Printf("Dedup:[%2v %2v](%3v %3v)-(%3v %3v) %v", x, y, s[x], s[y], s[x+2], s[y+2], sub)
		if s[x] == s[x+2] && s[y] == s[y+2] {
			continue
		}
		sub = append(sub, s[x], s[y])
	}
	// Let's check the first two and the last two to make sure they are not the same.
	if !(s[0] == s[len(s)-2] && s[1] == s[len(s)-1]) {
		sub = append(sub, s[len(s)-2], s[len(s)-1])
	}
	log.Println("Dedup:", sub)
	return sub
}
func Polygon(r Region, w WindingOrder, s []float64) (clippedSubjects [][]float64) {
	log.Println("Starting to clip polygon")
	sublen := len(s) / 2
	if sublen < 3 { // if we don't have at least three points we don't have a polygon.
		return clippedSubjects
	}
	clipper := r.Clipper(w)

	var subStart *point
	{
		var x, y, lx, ly = s[0], s[1], s[len(s)-2], s[len(s)-1]
		nxt := &point{
			X:       lx, // Last x
			Y:       ly, // Last y
			IsIn:    r.Contains(lx, ly),
			Type:    Subject,
			SubNext: subStart,
		}
		subStart = &point{
			X:       x,
			Y:       y,
			IsIn:    r.Contains(x, y),
			Type:    Subject,
			SubNext: nxt,
		}
	}
	subHead := subStart
	var intHead, intCurrent *point
	allInside := true  // Assume all points are inside of the region, to being with.
	allOutside := true // Paradoxically, assume all points are also outside of the region.

	insertClip := func(subStart, nxtPt *point, idx int, x, y int64, inward bool) {
		intPt := &point{
			X:       float64(x),
			Y:       float64(y),
			Type:    Intersect,
			IsIn:    inward,
			SubNext: nxtPt,
		}

		log.Println("Going to insert into sublist")

		insertIntoSubList(subStart, nxtPt, intPt)
		if idx == 3 {
			log.Println("Going to insert into cliplist ", idx, &clipper[idx], &clipper[0])
			insertIntoClipList(&clipper[idx], &clipper[0], intPt)
		} else {
			log.Println("Going to insert into cliplist ", idx, &clipper[idx], &clipper[idx+1])
			insertIntoClipList(&clipper[idx], &clipper[idx+1], intPt)
		}

		// Now we have to update the inter-list
		if intCurrent == nil {
			intCurrent = intPt
			intHead = intCurrent
		}

		intCurrent.IntNext = intPt
		intPt.IntNext = intHead
		intCurrent = intPt

	}

	log.Println("Walking through the points. ", sublen)
	for i, j := 2, 3; i <= len(s); i, j = i+2, j+2 {
		var nxtPt *point
		if i == len(s) {
			nxtPt = subHead
		} else {
			var x, y = s[i], s[j]
			nxtPt = &point{
				X:       x,
				Y:       y,
				IsIn:    r.Contains(x, y),
				Type:    Subject,
				SubNext: subStart.SubNext,
			}
		}
		inside := nxtPt.IsIn
		subStart.SubNext = nxtPt

		if allInside && !inside {
			allInside = false
		}
		if allOutside && inside {
			allOutside = false
		}
		log.Println("(", i, j, ") Looking at point(", i/2, "/", sublen, ") ", subStart, "to", nxtPt)

		toContinue := true
		switch {
		case subStart.IsIn && nxtPt.IsIn: // Both in nothing to do.
		case !subStart.IsIn && nxtPt.IsIn: // We are entering the clipping area.
			fallthrough
		case subStart.IsIn && !nxtPt.IsIn: // We are leaving the clipping area.
			toContinue = false
			fallthrough
		default:
			//			log.Println("Trying to find Crossings. Continue:", toContinue)
			r.FindCrossingsFunc(
				w,
				[4]int64{int64(subStart.X), int64(subStart.Y), int64(nxtPt.X), int64(nxtPt.Y)},
				func(idx int, x, y int64, inward bool) bool {
					log.Printf("Going to Insert intersection (%v %v [inward:%v]) between %v and %v\n", x, y, inward, subStart, nxtPt)
					insertClip(subStart, nxtPt, idx, x, y, inward)
					return toContinue
				},
			)
		}
		subStart = nxtPt
	}

	/*
		if allOutside {
			log.Println("All points are outside.")
			return clippedSubjects
		}
	*/

	if allInside {
		log.Println("All points are inside.")
		clippedSubjects = append(clippedSubjects, dedupSub(s))
		return clippedSubjects
	}

	walker := &PointWalker{Pt: intHead}
	var sub []float64
	for walker.Walk(func(p *point) {
		sub = append(sub, p.X, p.Y)
	}) {
		clippedSubjects = append(clippedSubjects, dedupSub(sub))
		sub = make([]float64, 0, 0)
	}

	if len(sub) > 0 {
		clippedSubjects = append(clippedSubjects, dedupSub(sub))
	}
	// We need to check to see if the region was completly enclose by the subject.
	if len(clippedSubjects) == 0 {
		if maths.Contains(s, clipper[0].X, clipper[0].Y) &&
			maths.Contains(s, clipper[1].X, clipper[1].Y) &&
			maths.Contains(s, clipper[2].X, clipper[2].Y) &&
			maths.Contains(s, clipper[3].X, clipper[3].Y) {
			clippedSubjects = append(clippedSubjects, []float64{
				clipper[0].X, clipper[0].Y,
				clipper[1].X, clipper[1].Y,
				clipper[2].X, clipper[2].Y,
				clipper[3].X, clipper[3].Y,
			})
		}
	}
	return clippedSubjects
}

/*
LineString clips the given line string to the given region.
*/
func LineString(r Region, s []float64) (clippedSubjects [][]float64) {
	clippedSubjects = make([][]float64, 0, 0)
	sub := make([]float64, 0, 0)
	subLen := len(s) / 2
	lIn, lX, lY := r.Contains(s[0], s[1]), s[0], s[1]
	for i := 1; i < subLen; i++ {
		idxX, idxY := i*2, (i*2)+1
		iIn, iX, iY := r.Contains(s[idxX], s[idxY]), s[idxX], s[idxY]
		switch {
		case lIn && iIn: // both points are in the clipping region.
			sub = append(sub, lX, lY)
		case !lIn && iIn: // Going into the clipping region.
			for _, ipt := range r.FindFirstCrossing(Clockwise, [4]int64{int64(lX), int64(lY), int64(iX), int64(iY)}) {
				sub = append(sub, float64(ipt[0]), float64(ipt[1]))
			}
		case !lIn && !iIn && r.DoesCross([4]float64{lX, lY, iX, iY}): // Both points are outside and crosses the region.
			if len(sub) > 0 {
				clippedSubjects = append(clippedSubjects, dedupSub(sub))
				sub = make([]float64, 0, 0)
			}
			for _, ipt := range r.FindCrossings(Clockwise, [4]int64{int64(lX), int64(lY), int64(iX), int64(iY)}) {
				sub = append(sub, float64(ipt[0]), float64(ipt[1]))
			}
			if len(sub) > 1 {
				clippedSubjects = append(clippedSubjects, dedupSub(sub))
				sub = make([]float64, 0, 0)
			}

		case lIn && !iIn: // Going out of the clipping region.
			sub = append(sub, lX, lY)
			for _, ipt := range r.FindFirstCrossing(Clockwise, [4]int64{int64(lX), int64(lY), int64(iX), int64(iY)}) {
				sub = append(sub, float64(ipt[0]), float64(ipt[1]))
			}
			if len(sub) > 1 {
				clippedSubjects = append(clippedSubjects, dedupSub(sub))
				sub = make([]float64, 0, 0)
			}
		}
		lIn, lX, lY = iIn, iX, iY
	}
	if len(sub) > 1 {
		clippedSubjects = append(clippedSubjects, dedupSub(sub))
		sub = make([]float64, 0, 0)
	}
	return clippedSubjects
}
