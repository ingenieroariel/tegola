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
	}

	for i, testcase := range testcases {
		for j, e := range testcase.expected {
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
			g := contains(pt, testcase.region)
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
		eerr    error
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
		log.Println("Starting TestClip test ", k)
		log.Printf("%+v", testcase)
		s, e := ClipPolygon(testcase.region, testcase.winding, testcase.subject)
		if e != testcase.eerr {
			t.Errorf("Test %v: Expected error to be %v got: %v", k, e, testcase.eerr)
		}
		if testcase.eerr != nil {
			continue
		}
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
