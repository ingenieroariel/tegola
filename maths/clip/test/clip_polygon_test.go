package clip_polygon_test

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/gdey/tbl"
	"github.com/terranodo/tegola/maths/clip"
)

type TestCase struct {
	region  clip.Region
	winding clip.WindingOrder
	subject []float64
	e       [][]float64
}

type TestCases []TestCase

const PixelWidth = 10

var showPng = flag.Bool("drawPNG", false, "Draw the PNG for the test cases even if the testcase passes.")

var testRegion = []clip.Region{
	clip.Region{0, 0, 10, 10},
	clip.Region{2, 2, 8, 8},
	clip.Region{-1, -1, 11, 11},
	clip.Region{-2, -2, 12, 12},
	clip.Region{-3, -3, 13, 13},
	clip.Region{-4, -4, 14, 14},
	clip.Region{5, 1, 7, 3},
	clip.Region{0, 5, 2, 7},
	clip.Region{-1, 4, 5, 8},
	clip.Region{5, 2, 11, 9},
}

// Drawing routings.
func minmax(s []float64, mix, miy, mx, my int) (minx, miny, maxx, maxy int) {
	minx = mix
	miny = miy
	maxx = mx
	maxy = my
	for i := 0; i < len(s); i += 2 {
		if int(s[i]) > maxx {
			maxx = int(s[i])
		}
		if int(s[i]) < minx {
			minx = int(s[i])
		}
		if int(s[i+1]) > maxy {
			maxy = int(s[i+1])
		}
		if int(s[i+1]) < miny {
			miny = int(s[i+1])
		}
	}
	return minx - 1, miny - 1, maxx + 1, maxy + 1
}

func drawLine(img *image.RGBA, pt1, pt2 image.Point, c color.RGBA) {

	if pt1.X == pt2.X && pt1.Y == pt2.Y {
		img.Set(pt1.X, pt1.Y, c)
		return
	}

	sx := pt1.X
	mx := pt2.X
	if pt2.X < pt1.X {
		sx = pt2.X
		mx = pt1.X
	}

	sy := pt1.Y
	my := pt2.Y
	if pt2.Y < pt1.Y {
		sy = pt2.Y
		my = pt1.Y
	}

	img.Set(sx, sy, c)
	img.Set(mx, my, c)
	xdelta := mx - sx

	// We have a veritcal line.
	if xdelta == 0 {
		for y := sy; y < my; y++ {
			img.Set(sx, y, c)
		}
		return
	}
	ydelta := my - sy
	if ydelta == 0 {
		for x := sx; x < mx; x++ {
			img.Set(x, sy, c)
		}
		return
	}
	m := int(ydelta / xdelta)
	b := int(sy - (m * sx))
	//y = mx+b
	for x := sx; x < mx; x++ {
		y := (m * x) + b
		img.Set(x, y, c)
	}
}

func scaleToPoint(minx, miny int, x, y float64) image.Point {
	sx, sy := (int(x)-minx)*PixelWidth, (int(y)-miny)*PixelWidth
	return image.Pt(sx, sy)
}

func drawRegion(img *image.RGBA, minx, miny int, r clip.Region, c color.RGBA) {

	drawLine(img, scaleToPoint(minx, miny, r[0], r[1]), scaleToPoint(minx, miny, r[2], r[1]), c)
	drawLine(img, scaleToPoint(minx, miny, r[2], r[1]), scaleToPoint(minx, miny, r[2], r[3]), c)
	drawLine(img, scaleToPoint(minx, miny, r[2], r[3]), scaleToPoint(minx, miny, r[0], r[3]), c)
	drawLine(img, scaleToPoint(minx, miny, r[0], r[3]), scaleToPoint(minx, miny, r[0], r[1]), c)
}

func drawSegment(img *image.RGBA, minx, miny int, s []float64, c color.RGBA) {
	if len(s) == 0 {
		return
	}
	pt := scaleToPoint(minx, miny, s[len(s)-2], s[len(s)-1])

	for i := 0; i < len(s); i += 2 {
		npt := scaleToPoint(minx, miny, s[i], s[i+1])
		drawLine(img, pt, npt, c)
		pt = npt
	}

}

func drawTestCase(tc *TestCase, got [][]float64, filename string) {
	log.Println("Creating png: ", filename)

	s := tc.subject
	r := tc.region
	minx, miny, maxx, maxy := minmax(s, int(r[0]), int(r[1]), int(r[2]), int(r[3]))

	for _, i := range got {
		minx, miny, maxx, maxy = minmax(i, minx, miny, maxx, maxy)
	}

	xdelta := int(maxx - minx)
	ydelta := int(maxy - miny)
	m := image.NewRGBA(image.Rect(0, 0, (xdelta * PixelWidth), (ydelta * PixelWidth)))
	for x := 0; x < (xdelta * PixelWidth); x++ {
		for y := 0; y < (ydelta * PixelWidth); y++ {
			m.Set(x, y, color.RGBA{255, 255, 255, 255})
			if y == 0 || x == 0 || y == ((ydelta*PixelWidth)-1) || x == ((xdelta*PixelWidth)-1) {
				continue
			}
			if y%PixelWidth == 0 || x%PixelWidth == 0 {
				m.Set(x, y, color.RGBA{0xF0, 0xF0, 0xF0, 0xFF})
			}
		}
	}
	prl := color.RGBA{255, 0, 255, 100}
	orn := scaleToPoint(minx, miny, 0, 0)
	/*
		drawLine(m, orn, scaleToPoint(minx, miny, 0, 1), prl)
		drawLine(m, orn, scaleToPoint(minx, miny, 0, -1), prl)
		drawLine(m, orn, scaleToPoint(minx, miny, 1, 0), prl)
		drawLine(m, orn, scaleToPoint(minx, miny, -1, 0), prl)
	*/
	//drawRegion(m, minx, miny, r, color.RGBA{0xBD, 0x15, 0x50, 0xAF})
	drawRegion(m, minx, miny, r, color.RGBA{0xD0, 0xD0, 0xD0, 0xFF})
	drawSegment(m, minx, miny, s, color.RGBA{0xA0, 0xA0, 0xA0, 0xFF})

	for _, i := range tc.e {
		drawSegment(m, minx, miny, i, color.RGBA{0xE9, 0x7F, 0x02, 0xFF})
	}
	for _, i := range got {
		drawSegment(m, minx, miny, i, color.RGBA{0x8A, 0x9B, 0x0F, 0xFF})
	}
	prl.A = 255
	m.Set(orn.X, orn.Y, prl)
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating file %v: %v\n", filename, err)
		return
	}
	png.Encode(f, m)

}

func TestClipPolygon(t *testing.T) {

	test := tbl.Cases(
		TestCase{
			region:  testRegion[0],
			winding: clip.Clockwise,
			subject: []float64{-2, 1, 2, 1, 2, 2, -1, 2, -1, 11, 2, 11, 2, 4, 4, 4, 4, 13, -2, 13},
			e: [][]float64{
				{0, 1, 2, 1, 2, 2, 0, 2},
				{2, 10, 2, 4, 4, 4, 4, 10},
			},
		},
		TestCase{
			region:  testRegion[0],
			winding: clip.Clockwise,
			subject: []float64{-2, 1, 12, 1, 12, 2, -1, 2, -1, 11, 2, 11, 2, 4, 4, 4, 4, 13, -2, 13},
			e: [][]float64{
				{0, 1, 10, 1, 10, 2, 0, 2},
				{2, 10, 2, 4, 4, 4, 4, 10},
			},
		},
		TestCase{
			region:  testRegion[0],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{0, 9, 10, 9, 10, 2, 5, 2, 5, 8, 0, 8},
				{0, 4, 3, 4, 3, 1, 0, 1},
			},
		},
		TestCase{
			region:  testRegion[1],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{5, 2, 5, 8, 8, 8, 8, 2},
				{2, 4, 3, 4, 3, 2, 2, 2},
			},
		},
		TestCase{
			region:  testRegion[2],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{-1, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8},
				{-1, 4, 3, 4, 3, 1, -1, 1},
			},
		},
		TestCase{
			region:  testRegion[3],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				{-2, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1, -2, 1},
			},
		},
		TestCase{
			region:  testRegion[4],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{-3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1, -3, 1},
			},
		},
		TestCase{
			region:  testRegion[5],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			},
		},
		TestCase{
			region:  testRegion[6],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{7, 2, 5, 2, 5, 3, 7, 3},
			},
		},
		TestCase{
			region:  testRegion[7],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e:       [][]float64{},
		},
		TestCase{
			region:  testRegion[8],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e:       [][]float64{},
		},
		TestCase{
			region:  testRegion[9],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{5, 2, 5, 9, 11, 9, 11, 2},
			},
		},
		TestCase{
			region:  testRegion[9],
			winding: clip.CounterClockwise,
			subject: []float64{-3, 1, -3, 10, 12, 10, 12, 1, 4, 1, 4, 8, -1, 8, -1, 4, 3, 4, 3, 1},
			e: [][]float64{
				[]float64{5, 2, 5, 9, 11, 9, 11, 2},
			},
		},
		TestCase{
			region:  testRegion[0],
			winding: clip.CounterClockwise,
			subject: []float64{-3, -3, -3, 10, 12, 10, 12, 1, 4, 1, 4, 8, -1, 8, -1, 4, 3, 4, 3, 3},
			e: [][]float64{
				[]float64{0, 0, 0, 4, 3, 4, 3, 3},
				[]float64{4, 1, 4, 8, 0, 8, 0, 10, 10, 10, 10, 1},
			},
		},
	)

	test.Run(func(i int, tc TestCase) {
		var drawPng bool
		t.Log("Starting test ", i)
		log.Println("Starting test ", i)
		got := clip.Polygon(tc.region, tc.winding, tc.subject)
		if len(tc.e) != len(got) {
			t.Errorf("Test %v: Expected number of slices to be %v got: %v -- %+v", i, len(tc.e), len(got), got)
			drawTestCase(&tc, got, fmt.Sprintf("testcase%v.png", i))
			return
		}
		for j := range tc.e {

			if len(tc.e[j]) != len(got[j]) {
				drawPng = true
				t.Errorf("Test %v: Expected slice %v to have %v items got: %v -- %+v", i, i, len(tc.e[j]), len(got[j]), got[j])
				continue
			}
			for k := 0; k < len(tc.e[j])/2; k++ {
				k1 := k * 2
				k2 := k1 + 1
				if (tc.e[j][k1] != got[j][k1]) || (tc.e[j][k2] != got[j][k2]) {
					drawPng = true
					t.Errorf("Test %v: Expected Sice: %v  item: %v to be ( %v %v ) got: ( %v %v)", i, j, k, tc.e[j][k1], tc.e[j][k2], got[j][k1], got[j][k2])
				}
			}
		}
		if drawPng || *showPng {
			drawTestCase(&tc, got, fmt.Sprintf("testcase%v.png", i))
		}

	})
}
