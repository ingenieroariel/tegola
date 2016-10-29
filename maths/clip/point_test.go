package clip

import "testing"

func TestIsInboundIntersect(t *testing.T) {
	pt := &point{
		Type: Intersect,
		seen: true,
		IsIn: true,
	}

	if !pt.IsInboundIntersect() {
		t.Errorf("Failed with Intersect point(%v) and IsIn(%v).", pt, pt.IsIn)
	}
	pt.IsIn = false
	if pt.IsInboundIntersect() {
		t.Errorf("Failed with Intersect point(%v) and IsIn(%v).", pt, pt.IsIn)
	}
	pt.Type = Subject
	pt.IsIn = true
	if pt.IsInboundIntersect() {
		t.Errorf("Failed with Subject point(%v) and IsIn(%v).", pt, pt.IsIn)
	}
	pt.IsIn = false
	if pt.IsInboundIntersect() {
		t.Errorf("Failed with Subject point(%v) and IsIn(%v).", pt, pt.IsIn)
	}
	pt.Type = Clipper
	pt.IsIn = true
	if pt.IsInboundIntersect() {
		t.Errorf("Failed with Clipper point(%v) and IsIn(%v).", pt, pt.IsIn)
	}
	pt.IsIn = false
	if pt.IsInboundIntersect() {
		t.Errorf("Failed with Clipper point(%v) and IsIn(%v).", pt, pt.IsIn)
	}
}
