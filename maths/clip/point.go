package clip

import (
	"fmt"
	"log"
)

type PointType uint8

const (
	Clipper = PointType(iota)
	Subject
	Intersect
)

func clipnext(pt *point) *point {
	return pt.ClipNext
}
func subnext(pt *point) *point {
	return pt.SubNext
}
func intnext(pt *point) *point {
	return pt.IntNext
}

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
	// The Next point for the clip
	ClipNext *point
	IntNext  *point
	seen     bool
}

func (pt *point) AppendSub(i *point) {
	if pt == nil {
		return
	}
	if i == nil {
		pt.SubNext = nil
		return
	}
	pt.SubNext, i.SubNext = i, pt.SubNext
}
func (pt *point) AppendClip(i *point) {
	if pt == nil {
		return
	}
	if i == nil {
		pt.ClipNext = nil
		return
	}
	pt.ClipNext, i.ClipNext = i, pt.ClipNext
}
func (pt *point) AppendInt(i *point) {
	if pt == nil {
		return
	}
	if i == nil {
		pt.IntNext = nil
		return
	}
	pt.IntNext, i.IntNext = i, pt.IntNext
}
func (pt *point) Next() *point {
	log.Println("Next point called, this will moved the point depending on the type.", pt)
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
			log.Println("Inbound intersect should be Subject next", pt.SubNext)
			return pt.SubNext
		}
		log.Println("Outbound intersect should be clip next", pt.ClipNext)
		return pt.ClipNext
	}
	panic("Should not get here.")
}

func (pt point) String() string {
	switch pt.Type {
	case Clipper:
		return fmt.Sprintf("Clipper( %v %v )", pt.X, pt.Y)
	case Subject:
		return fmt.Sprintf("Subject( %v %v )", pt.X, pt.Y)
	case Intersect:
		ms := "<[o]"
		if pt.IsIn {
			ms = "[i]>"
		}
		return fmt.Sprintf("Intersect%v( %v %v )", ms, pt.X, pt.Y)
	default:
		return fmt.Sprintf("Point( %v %v )", pt.X, pt.Y)
	}
}

func (pt *point) IsInboundIntersect() bool {
	if pt == nil {
		log.Println("pt is nil")
		return false
	}
	return pt.Type == Intersect && pt.IsIn
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

type PointWalker struct {
	Pt  *point
	Err error
}

func (p *PointWalker) Walk(fn func(p *point)) bool {

	if p.Pt == nil {
		return false
	}

	// First we need to move to an Inbound point.
	log.Println("Finding an inbound point.")
	spt := p.Pt
	if !spt.IsInboundIntersect() {
		log.Println("Get an inbound Intersect point.")
		spt = spt.NextInboundIntersect()
	}
	if spt == nil {
		p.Err = fmt.Errorf("Did not find an inbound intersect point. Something is wrong.")
		return false
	}

	head := spt
	fn(spt)
	spt.seen = true
	spt = spt.Next()
	for spt != head {
		if spt == nil {
			p.Err = fmt.Errorf("Got a nil pointer!")
			return false
		}
		if spt.seen {
			p.Err = fmt.Errorf("Revisting a node....")
			return false
		}
		fn(spt)
		spt.seen = true
		spt = spt.Next()
	}
	p.Pt = p.Pt.NextInboundIntersect()
	return !(p.Pt == nil)
}

type nextfn func(pt *point) *point

// findListPos will find the next point, to statify the cmp function.
func findListPos(start *point, cmp func(pt *point) bool, next nextfn) (*point, error) {
	//	log.Println("Looking for a position")
	startpt := start
	for endpt := next(start); endpt != start; startpt, endpt = endpt, next(startpt) {
		if cmp(endpt) {
			return startpt, nil
		}
	}
	return start, fmt.Errorf("Went around the entire list and id not find a pos")
}

// inTheMiddle will return -1, 0, or 1 depending on where m is between s and e.
func inTheMiddle(s, m, e float64) int8 {
	switch {
	case m < s:
		return -1
	case m > e:
		return 1
	default:
		return 0
	}
}

func findInsertionPos(start, end, item *point, next nextfn) (pt *point, err error) {
	yMatch := start.Y == end.Y
	xGreater := start.X < end.X
	yGreater := start.Y < end.Y
	switch {
	case yMatch && xGreater:
		switch inTheMiddle(start.X, item.X, end.X) {
		default:
			err = fmt.Errorf("Increase in X item is after end!")
		case -1:
			err = fmt.Errorf("Increase in X item is before start!")
		case 0:
			pt, err = findListPos(start, func(pt *point) bool { return item.X <= pt.X }, next)
		}
	case yMatch && !xGreater:
		switch inTheMiddle(end.X, item.X, start.X) {
		default:
			err = fmt.Errorf("Decreasing in X item is before start!")
		case -1:
			err = fmt.Errorf("Decreasing in X item is after end!")
		case 0:
			pt, err = findListPos(start, func(pt *point) bool { return item.X >= pt.X }, next)
		}
	case !yMatch && yGreater:
		switch inTheMiddle(start.Y, item.Y, end.Y) {
		default:
			err = fmt.Errorf("Increase in Y item is after end!")
		case -1:
			err = fmt.Errorf("Increase in Y item is before start!")
		case 0:
			pt, err = findListPos(start, func(pt *point) bool { return item.Y <= pt.Y }, next)
		}
	case !yMatch && !yGreater:
		switch inTheMiddle(end.Y, item.Y, start.Y) {
		default:
			err = fmt.Errorf("Decreasing in Y item is before start!")
		case -1:
			err = fmt.Errorf("Decreasing in Y item is after end!")
		case 0:
			pt, err = findListPos(start, func(pt *point) bool { return item.Y >= pt.Y }, next)
		}
	}

	if err != nil {
		log.Printf("Did not find point because of error: %v", err)
	}
	log.Printf("Found the following pt(%v) between %v and %v for item %v.", pt, start, end, item)
	return pt, err

}

func insertIntoClipList(s, e, i *point) error {
	pt, err := findInsertionPos(s, e, i, clipnext)
	if err != nil {
		return err
	}
	pt.AppendClip(i)
	return nil
}
func insertIntoSubList(s, e, i *point) error {
	pt, err := findInsertionPos(s, e, i, subnext)
	if err != nil {
		return err
	}
	pt.AppendSub(i)
	return nil
}
func insertIntoIntList(s, e, i *point) error {
	pt, err := findInsertionPos(s, e, i, intnext)
	if err != nil {
		return err
	}
	pt.AppendInt(i)
	return nil
}
