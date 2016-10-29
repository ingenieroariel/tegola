package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/terranodo/tegola/maths/clip"
)

type TestCase struct {
	region  clip.Region
	winding clip.WindingOrder
	subject []int64
	e       [][]int64
}

type TestCases []TestCase

const PixelWidth = 10

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

var testcases = TestCases{
	{
		region:  testRegion[0],
		winding: clip.CounterClockwise,
		subject: []int64{-3, 1, -3, 9, 11, 9, 11, 2, 5, 2, 5, 8, -1, 8, -1, 4, 3, 4, 3, 1},
		e: [][]int64{
			{0, 9, 10, 9, 10, 2, 5, 2, 5, 8, 0, 8},
			{0, 4, 3, 4, 3, 1, 0, 1},
		},
	},
}

func minmax(s []int64, r clip.Region) (minx, miny, maxx, maxy int) {
	minx = int(r[0])
	miny = int(r[1])
	maxx = int(r[2])
	maxy = int(r[3])
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

func scaleToPoint(minx, miny int, x, y int64) image.Point {
	sx, sy := (int(x)-minx)*PixelWidth, (int(y)-miny)*PixelWidth
	return image.Pt(sx, sy)
}

func drawRegion(img *image.RGBA, minx, miny int, r clip.Region, c color.RGBA) {

	drawLine(img, scaleToPoint(minx, miny, int64(r[0]), int64(r[1])), scaleToPoint(minx, miny, int64(r[2]), int64(r[1])), c)
	drawLine(img, scaleToPoint(minx, miny, int64(r[2]), int64(r[1])), scaleToPoint(minx, miny, int64(r[2]), int64(r[3])), c)
	drawLine(img, scaleToPoint(minx, miny, int64(r[2]), int64(r[3])), scaleToPoint(minx, miny, int64(r[0]), int64(r[3])), c)
	drawLine(img, scaleToPoint(minx, miny, int64(r[0]), int64(r[3])), scaleToPoint(minx, miny, int64(r[0]), int64(r[1])), c)
}

func drawSegment(img *image.RGBA, minx, miny int, s []int64, c color.RGBA) {
	pt := scaleToPoint(minx, miny, s[len(s)-2], s[len(s)-1])

	for i := 0; i < len(s); i += 2 {
		npt := scaleToPoint(minx, miny, s[i], s[i+1])
		drawLine(img, pt, npt, c)
		pt = npt
	}

}

func main() {

	s := testcases[0].subject
	r := testcases[0].region
	minx, miny, maxx, maxy := minmax(s, r)
	xdelta := int(maxx - minx)
	ydelta := int(maxy - miny)

	log.Println("minx", minx, "miny", miny, "maxx", maxx, "maxy", maxy)

	m := image.NewRGBA(image.Rect(0, 0, (xdelta * PixelWidth), (ydelta * PixelWidth)))
	for x := 0; x < (xdelta * PixelWidth); x++ {
		for y := 0; y < (ydelta * PixelWidth); y++ {
			m.Set(x, y, color.RGBA{255, 255, 255, 255})
			if y == 0 || x == 0 || y == ((ydelta*PixelWidth)-1) || x == ((xdelta*PixelWidth)-1) {
				continue
			}
			if y%PixelWidth == 0 || x%PixelWidth == 0 {
				m.Set(x, y, color.RGBA{16, 16, 16, 153})
			}
		}
	}
	prl := color.RGBA{255, 0, 255, 100}
	orn := scaleToPoint(minx, miny, 0, 0)
	drawLine(m, orn, scaleToPoint(minx, miny, 0, 1), prl)
	drawLine(m, orn, scaleToPoint(minx, miny, 0, -1), prl)
	drawLine(m, orn, scaleToPoint(minx, miny, 1, 0), prl)
	drawLine(m, orn, scaleToPoint(minx, miny, -1, 0), prl)
	drawSegment(m, minx, miny, s, color.RGBA{255, 0, 0, 255})
	drawRegion(m, minx, miny, r, color.RGBA{0, 255, 0, 255})
	drawSegment(m, minx, miny, testcases[0].e[0], color.RGBA{0, 0, 255, 155})
	drawSegment(m, minx, miny, testcases[0].e[1], color.RGBA{0, 0, 255, 155})
	prl.A = 255
	m.Set(orn.X, orn.Y, prl)

	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, m)
}
