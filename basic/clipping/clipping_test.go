package clipping

import (
	"fmt"
	"log"
	"testing"
)

func TestInClip(t *testing.T) {
	testcases := []struct {
		region   ClippingRegion
		pts      []float64
		expected []bool
	}{
		{
			region:   ClippingRegion([4]float64{0, 0, 10, 10}),
			pts:      []float64{5, 5, 0, 0, 10, 10, 0, 5, 4, 4, 11, 11, -1, -1},
			expected: []bool{true, false, false, false, true, false, false},
		},
		{
			region:   ClippingRegion([4]float64{2, 2, 8, 8}),
			pts:      []float64{5, 2},
			expected: []bool{false},
		},
		{
			region:   ClippingRegion([4]float64{-2048, 2048, 6939, -6939}),
			pts:      []float64{2797, -70, 2780, -62, 920, 2059, 927, 2027},
			expected: []bool{true, true, false, true},
		},
	}

	for i, testcase := range testcases {
		for j, e := range testcase.expected {
			ptx, pty := testcase.pts[j*2], testcase.pts[(j*2)+1]
			g := inClip(testcase.region, ptx, pty)
			if g != e {
				t.Errorf("For Test (%v, %v) Got: %v Expected: %v, Clipping Region: %v, Point (%v,%v).", i, j, g, e, testcase.region, ptx, pty)
			}
		}
	}

}

func TestDoesCrossClip(t *testing.T) {
	testcases := []struct {
		region   ClippingRegion
		pts      []point
		expected []bool
	}{
		{
			region: ClippingRegion([4]float64{0, 0, 10, 10}),
			pts: []point{
				point{X: 5, Y: -1},
				point{X: 5, Y: 11},
				point{X: 5, Y: 11},
				point{X: 5, Y: -1},
				point{X: -1, Y: 5},
				point{X: 11, Y: 5},
				point{X: 11, Y: 5},
				point{X: -1, Y: 5},
				point{X: -1, Y: -1},
				point{X: 11, Y: 11},
				point{X: 11, Y: 11},
				point{X: -1, Y: -1},
				point{X: -1, Y: -1},
				point{X: -1, Y: 11},
				point{X: -1, Y: 11},
				point{X: -1, Y: -1},
				point{X: -1, Y: -1},
				point{X: 11, Y: -1},
				point{X: 11, Y: -1},
				point{X: -1, Y: -1},
				point{X: -2, Y: 1},
				point{X: 12, Y: 1},
			},
			expected: []bool{true, true, true, true, true, true, false, false, false, false, true},
		},
		{
			region: ClippingRegion([4]float64{5, 2, 11, 9}),
			pts: []point{
				point{X: -3, Y: 9},
				point{X: 11, Y: 9},
			},
			expected: []bool{true},
		},
		{
			region: ClippingRegion{-2048, 2048, 4493, -4493},
			pts: []point{
				point{X: 3723, Y: 2048},
				point{X: 3710, Y: 2056},
			},
			expected: []bool{false},
		},
	}

	for i, testcase := range testcases {
		for j, e := range testcase.expected {
			log.Println("Test", i, j)
			pt1, pt2 := testcase.pts[j*2], testcase.pts[(j*2)+1]
			g := doesCrossClip(testcase.region, [4]float64{pt1.X, pt1.Y, pt2.X, pt2.Y})
			if g != e {
				t.Errorf("For Test (%v, %v) Got: %v Expected: %v, Clipping Region: %v, Point (%v,%v).", i, j, g, e, testcase.region, pt1, pt2)
			}
		}
	}

}

func doesPanic(f func()) (msg string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprintf("%v", r)
			ok = true
		}
	}()
	f()
	return
}

func TestIntersect(t *testing.T) {
	testcases := []struct {
		line1    [4]float64
		line2    [4]float64
		expected [2]float64
		ok       bool
	}{
		{
			line1:    [4]float64{-10, 0, 10, 0},
			line2:    [4]float64{0, 10, 0, -10},
			expected: [2]float64{0, 0},
			ok:       true,
		},
		{
			line2:    [4]float64{-10, 0, 10, 0},
			line1:    [4]float64{0, 10, 0, -10},
			expected: [2]float64{0, 0},
			ok:       true,
		},
		{
			line2: [4]float64{-10, 0, 10, 0},
			line1: [4]float64{-10, 10, 10, 10},
			ok:    false,
		},
		{
			line2: [4]float64{10, -10, 10, 10},
			line1: [4]float64{-1, -10, -1, 10},
			ok:    false,
		},
		{
			line1: [4]float64{-1, -1, -11, -11},
			line2: [4]float64{-11, -11, -1, -1},
			ok:    false,
		},
		{
			line1:    [4]float64{0, 0, 10, 10},
			line2:    [4]float64{10, 0, 0, 10},
			expected: [2]float64{5, 5},
			ok:       true,
		},
		{
			line1:    [4]float64{-2048, 2038, 4493, 2038},
			line2:    [4]float64{927, 2027, 920, 2059},
			expected: [2]float64{924.59375, 2038},
			ok:       true,
		},
	}

	for i, testcase := range testcases {
		x, y, ok := intersect(testcase.line1, testcase.line2)
		if ok != testcase.ok {
			if testcase.ok {
				t.Errorf("For test case %v : Expected to be ok, but it wasn't: %v %v %v", i, x, y, ok)
			} else {
				t.Errorf("For test case %v: Did not expect it to be ok.", i)
			}
			continue
		}
		if !testcase.ok {
			continue
		}
		if x != testcase.expected[0] {
			t.Errorf("For test case %v: Expected x to be %v got %v", i, testcase.expected[0], x)
		}
		if y != testcase.expected[1] {
			t.Errorf("For test case %v: Expected y to be %v got %v", i, testcase.expected[1], y)
		}
	}
}

func TestContains(t *testing.T) {
	testcases := []struct {
		region   []float64
		pts      []float64
		expected []bool
	}{
		{
			region:   []float64{0, 0, 0, 10, 10, 10, 10, 0},
			pts:      []float64{5, 5, 0, 0, -1, -1, 9, 9, 9, 0, 9, 1},
			expected: []bool{true, false, false, true, false, true},
		},
	}
	for i, testcase := range testcases {
		for j, e := range testcase.expected {
			pt := [2]float64{testcase.pts[j*2], testcase.pts[(j*2)+1]}
			g := Contains(pt, testcase.region)
			if g != e {
				t.Errorf("For test %v.%v: Expected %v got %v, region: %v pt: %v", i, j, e, g, testcase.region, pt)
			}
		}
	}
}

func TestClip(t *testing.T) {
	testcases := []struct {
		region  ClippingRegion
		winding WindingOrder
		subject []float64
		e       [][]float64
	}{
		{
			region:  ClippingRegion{0, 0, 10, 10},
			winding: Clockwise,
			subject: []float64{-2, 1, 2, 1, 2, 2, -1, 2, -1, 11, 2, 11, 2, 4, 4, 4, 4, 13, -2, 13},
			e: [][]float64{
				{0, 1, 2, 1, 2, 2, 0, 2},
				{2, 10, 2, 4, 4, 4, 4, 10},
			},
		},
		{
			region:  ClippingRegion{0, 0, 10, 10},
			winding: Clockwise,
			subject: []float64{-2, 1, 12, 1, 12, 2, -1, 2, -1, 11, 2, 11, 2, 4, 4, 4, 4, 13, -2, 13},
			e: [][]float64{
				{0, 1, 10, 1, 10, 2, 0, 2},
				{2, 10, 2, 4, 4, 4, 4, 10},
			},
		},
		{
			region:  ClippingRegion{0, 0, 10, 10},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{0, 9, 10, 9, 10, 2, 5, 2, 5, 8, 0, 8},
				{0, 4, 3, 4, 3, 1, 0, 1},
			},
		},
		{
			region:  ClippingRegion{2, 2, 8, 8},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{5, 2, 5, 8, 8, 8, 8, 2},
				{2, 4, 3, 4, 3, 2, 2, 2},
			},
		},
		{
			region:  ClippingRegion{-1, -1, 11, 11},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{-1, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8},
				{-1, 4, 3, 4, 3, 1, -1, 1},
			},
		},
		{
			region:  ClippingRegion{-2, -2, 12, 12},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{-2, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1, -2, 1},
			},
		},
		{
			region:  ClippingRegion{-3, -3, 13, 13},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{-3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1, -3, 1},
			},
		},
		{
			region:  ClippingRegion{-4, -4, 14, 14},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			},
		},
		{
			region:  ClippingRegion{5, 1, 7, 3},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{7, 2, 5, 2, 5, 3, 7, 3},
			},
		},
		{
			region:  ClippingRegion{0, 5, 2, 7},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e:       [][]float64{},
		},
		{
			region:  ClippingRegion{-1, 4, 5, 8},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e:       [][]float64{},
		},
		{
			region:  ClippingRegion{5, 2, 11, 9},
			winding: CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{5, 9, 11, 9, 11, 2, 5, 2},
			},
		},
	}
	for k, testcase := range testcases {
		/*
			if k != 8 {
				continue
			}
		*/
		s := ClipPolygon(testcase.region, testcase.winding, testcase.subject)
		if len(testcase.e) != len(s) {
			t.Errorf("Test %v: Expected number of slices to be %v got: %v -- %+v", k, len(testcase.e), len(s), s)
			continue
		}
		for i := range testcase.e {
			if len(testcase.e[i]) != len(s[i]) {
				t.Errorf("Test %v: Expected slice %v to have %v items got: %v -- %+v", k, i, len(testcase.e[i]), len(s[i]), s[i])
				continue
			}
			for j := 0; j < len(testcase.e[i])/2; j++ {
				j1 := j * 2
				j2 := j1 + 1
				if (testcase.e[i][j1] != s[i][j1]) || (testcase.e[i][j2] != s[i][j2]) {
					t.Errorf("Test %v: Expected Sice: %v  item: %v to be ( %v %v ) got: ( %v %v)", k, i, j, testcase.e[i][j1], testcase.e[i][j2], s[i][j1], s[i][j2])
				}
			}
		}
	}
}

func TestClipLineString(t *testing.T) {
	testcases := []struct {
		region  ClippingRegion
		subject []float64
		e       [][]float64
	}{
		{
			region:  ClippingRegion{0, 0, 10, 10},
			subject: []float64{-2, 1, 2, 1, 2, 2, -1, 2, -1, 11, 2, 11, 2, 4, 4, 4, 4, 13, -2, 13},
			e: [][]float64{
				{0, 1, 2, 1, 2, 2, 0, 2},
				{2, 10, 2, 4, 4, 4, 4, 10},
			},
		},
		{
			region:  ClippingRegion{0, 0, 10, 10},
			subject: []float64{-2, 1, 12, 1, 12, 2, -1, 2, -1, 11, 2, 11, 2, 4, 4, 4, 4, 13, -2, 13},
			e: [][]float64{
				{0, 1, 10, 1},
				{0, 2, 10, 2},
				{2, 10, 2, 4, 4, 4, 4, 10},
			},
		},
	}
	for k, testcase := range testcases {
		s := ClipLineString(testcase.region, testcase.subject)
		if len(testcase.e) != len(s) {
			t.Errorf("Test %v: Expected number of slices to be %v got: %v -- %+v", k, len(testcase.e), len(s), s)
			continue
		}
		for i := range testcase.e {
			if len(testcase.e[i]) != len(s[i]) {
				t.Errorf("Test %v: Expected slice %v to have %v items got: %v -- %+v", k, i, len(testcase.e[i]), len(s[i]), s[i])
				continue
			}

			for j := 0; j < len(testcase.e[i])/2; j++ {
				j1 := j * 2
				j2 := j1 + 1
				if (testcase.e[i][j1] != s[i][j1]) || (testcase.e[i][j2] != s[i][j2]) {
					t.Errorf("Test %v: Expected Sice: %v  item: %v to be ( %v %v ) got: ( %v %v)", k, i, j, testcase.e[i][j1], testcase.e[i][j2], s[i][j1], s[i][j2])
				}
			}
		}
	}
}

type doesCrossAxisFloatTestCase struct {
	region  ClippingRegion
	winding WindingOrder
	pts     [4]float64

	epts [8]float64
	eOk  [4]bool
}

func newDoesCrossAxisFloatTestCase(r ClippingRegion, w WindingOrder, pts [4]float64, eoks [4]bool, epts ...float64) doesCrossAxisFloatTestCase {
	var aepts [8]float64
	var o int
	for i, k := range eoks {
		if k {
			aepts[i*2] = epts[o]
			aepts[(i*2)+1] = epts[o+1]
			o += 2
		}
	}
	return doesCrossAxisFloatTestCase{
		region:  r,
		winding: w,
		pts:     pts,
		epts:    aepts,
		eOk:     eoks,
	}
}

func TestDoesCrossAxisFloat(t *testing.T) {

	testcases := []doesCrossAxisFloatTestCase{
		newDoesCrossAxisFloatTestCase(
			ClippingRegion{-2048, 2048, 4493, -4493},
			Clockwise,
			[4]float64{927, 2027, 920, 2059},
			[4]bool{false, true, false, false},
			924.59375, 2048,
		),
		newDoesCrossAxisFloatTestCase(
			ClippingRegion{-2048, 2048, 4493, -4493},
			Clockwise,
			[4]float64{3723, 2048, 3710, 2056},
			[4]bool{false, true, false, false},
			3723, 2048,
		),
	}
	for tIdx, testcase := range testcases {
		if tIdx != 1 {
			continue
		}
		for i, eok := range testcase.eOk {
			gx, gy, gok := doesCrossAxisFloat(
				testcase.region,
				i,
				testcase.winding,
				testcase.pts,
			)
			if gok != eok {
				t.Errorf("Failed test %v idx %v: Expected ok to be %v got %v", tIdx, i, eok, gok)
				continue
			}
			if !eok {
				continue
			}
			ex, ey := testcase.epts[i*2], testcase.epts[(i*2)+1]
			if gx != ex && gy != ey {
				t.Errorf("Failed test %v idx %v: Expected (x y) to be (%v %v) got (%v %v)", tIdx, i, ex, ey, gx, gy)
			}
		}

	}
}
