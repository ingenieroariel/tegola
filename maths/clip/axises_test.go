package clip_test

import (
	"testing"

	"github.com/gdey/tbl"
	"github.com/terranodo/tegola/maths/clip"
)

var testRegions = []clip.Region{
	{0, 0, 10, 10},
	{-10, -10, 0, 0},
	{2, 2, 8, 8},
}

func TestAxisesDoesCross(t *testing.T) {
	type ExpectedVal struct {
		ok bool
		x  int64
		y  int64
	}
	type testcase struct {
		r   *clip.Region
		pts [4]int64
		e   [8]ExpectedVal
	}
	test := tbl.Cases(
		testcase{
			r:   &testRegions[0],
			pts: [4]int64{-3, 1, -3, 9},
			e: [8]ExpectedVal{
				{ok: false},
				{ok: false, x: -3, y: 0},
				{ok: false},
				{ok: false, x: -3, y: 10},
				// CounterClockwise
				{ok: false},
				{ok: false, x: -3, y: 10},
				{ok: false},
				{ok: false, x: -3, y: 0},
			},
		},
		testcase{
			r:   &testRegions[0],
			pts: [4]int64{-3, -3, 1, 1},
			e: [8]ExpectedVal{
				{ok: true},
				{ok: true},
				{ok: false, x: 10, y: 10},
				{ok: false, x: 10, y: 10},
				// CounterClockwise
				{ok: true},
				{ok: false, x: 10, y: 10},
				{ok: false, x: 10, y: 10},
				{ok: true},
			},
		},
		testcase{
			r:   &testRegions[0],
			pts: [4]int64{11, 2, 5, 2},
			e: [8]ExpectedVal{
				{ok: false, x: 0, y: 2},
				{ok: false},
				{ok: true, x: 10, y: 2},
				{ok: false},
				// CounterClockwise
				{ok: false, x: 0, y: 2},
				{ok: false},
				{ok: true, x: 10, y: 2},
				{ok: false},
			},
		},
		testcase{
			r:   &testRegions[0],
			pts: [4]int64{2, 11, 2, 4},
			e: [8]ExpectedVal{
				{ok: false},
				{ok: false, x: 2, y: 0},
				{ok: false},
				{ok: true, x: 2, y: 10},
				// CounterClockwise
				{ok: false},
				{ok: true, x: 2, y: 10},
				{ok: false},
				{ok: false, x: 2, y: 0},
			},
		},
		testcase{
			r:   &testRegions[2],
			pts: [4]int64{5, 2, 5, 8},
			e: [8]ExpectedVal{
				{ok: false},
				{ok: false, x: 5, y: 2},
				{ok: false},
				{ok: false, x: 5, y: 8},
				// CounterClockwise
				{ok: false},
				{ok: false, x: 5, y: 8},
				{ok: false},
				{ok: false, x: 5, y: 2},
			},
		},
	)
	runTest := func(idx int, tc testcase, w clip.WindingOrder) {
		// Run though each of the axises.
		// The offset to start with.
		var offset int
		if w == clip.CounterClockwise {
			offset = 4
		}
		for i := 0; i < 4; i++ {
			a := tc.r.Axises(uint8(i), w)
			t.Logf("Running test(%3v:%1v):", idx, offset+i)
			e := tc.e[offset+i]
			x, y, ok := a.DoesCross(tc.pts)
			if e.ok != ok || e.x != x || e.y != y {
				t.Errorf("Failed  test(%3v:%1v): For Axises %v and Pts( %v ); Expected: %v Got: { ok: %v x: %v y: %v }", idx, offset+i, a, tc.pts, e, ok, x, y)
			}
		}
	}
	test.Run(func(idx int, tc testcase) {
		runTest(idx, tc, clip.Clockwise)
		runTest(idx, tc, clip.CounterClockwise)
	})

}

func TestAxisesIsInwardBound(t *testing.T) {

	//TODO(gdey): Move the points out of the test cases, and into their own structure; then reference then in the testcases. This should make it easer to construct test cases. Just like for the regions.
	var Pts = [][4]int64{
		{-3, -4, 1, 2},
		{-3, -3, 1, 1},
		{1, 1, 20, 20},
	}

	type testcase struct {
		r   *clip.Region
		pts [4]int64
		e   [8]bool
	}
	test := tbl.Cases(
		testcase{
			r:   &testRegions[0],
			pts: Pts[0],
			e:   [8]bool{true, true, false, false, true, false, false, true},
		},
		testcase{
			r:   &testRegions[0],
			pts: Pts[1],
			e:   [8]bool{true, true, false, false, true, false, false, true},
		},
		testcase{
			r:   &testRegions[0],
			pts: Pts[2],
		},

		testcase{
			r:   &testRegions[1],
			pts: Pts[0],
		},
		testcase{
			r:   &testRegions[1],
			pts: Pts[1],
		},
		testcase{
			r:   &testRegions[1],
			pts: Pts[2],
		},
	)
	runTest := func(idx int, tc testcase, w clip.WindingOrder) {
		// Run though each of the axises.
		// The offset to start with.
		var offset int
		if w == clip.CounterClockwise {
			offset = 4
		}
		for i := 0; i < 4; i++ {
			a := tc.r.Axises(uint8(i), w)
			t.Logf("Running test(%3v:%1v):", idx, offset+i)
			e := tc.e[offset+i]
			g := a.IsInwardBound(tc.pts)
			if e != g {
				t.Errorf("Failed test(%3v:%1v): Expected: %v Got: %v", idx, offset+i, e, g)
			}
		}
	}
	test.Run(func(idx int, tc testcase) {
		runTest(idx, tc, clip.Clockwise)
		runTest(idx, tc, clip.CounterClockwise)
	})
}
