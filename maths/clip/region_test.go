package clip_test

/*
func TestAxisesDoesCross(t *testing.T) {

	type testcase struct {
		desc   string
		obj    *clip.Axises
		points [4]int64
		x      int64
		y      int64
		ok     bool
	}

	region1 := clip.Region{0, 0, 10, 10}
	test := tbl.Cases(
		testcase{
			obj:    region1.Axises(3, clip.Clockwise),
			points: [4]int64{-3, -3, 1, 1},
			x:      10,
			y:      10,
			ok:     true,
		},
		testcase{
			obj:    region1.Axises(0, clip.Clockwise),
			points: [4]int64{-3, -3, 1, 1},
			x:      0,
			y:      0,
			ok:     true,
		},
		testcase{
			obj:    region1.Axises(0, clip.Clockwise),
			points: [4]int64{-3, -3, 1, 1},
			x:      0,
			y:      0,
			ok:     true,
		},
	)
	test.Run(func(idx int, tc testcase) {
		t.Logf("Running test %v — %v\n", idx, tc.desc)
		x, y, ok := tc.obj.DoesCross(tc.points)
		if x != tc.x {
			t.Errorf("For test %v: Got %v Expected %v for x value.", idx, x, tc.x)
		}
		if y != tc.x {
			t.Errorf("For test %v: Got %v Expected %v for y value.", idx, y, tc.y)
		}
		if ok != tc.ok {
			t.Errorf("For test %v: Got %v Expected %v for ok value.", idx, ok, tc.ok)
		}

	})
}
*/
