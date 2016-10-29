package basic

import (
	"log"

	"github.com/terranodo/tegola"
	"github.com/terranodo/tegola/maths"
	"github.com/terranodo/tegola/maths/clip"
)

const (
	// Version is the version of the basic package
	// grover: lkjdfshkljfsda
	Version = "v0.0.1"
)

type lineStringKeyMapping struct {
	Line []float64
	idx  int
}

func clipInteriorRegions(scr clip.Region, cps [][]float64, inline tegola.LineString) []lineStringKeyMapping {
	var pts []float64
	// Convert the linestring into a slice of points.
	for _, pt := range inline.Subpoints() {
		pts = append(pts, pt.X(), pt.Y())
	}
	cpts := clip.Polygon(scr, clip.CounterClockwise, pts)
	if len(cpts) == 0 {
		return nil
	}
	mapping := make([]lineStringKeyMapping, 0, len(cpts))
	// Now we need to map the internal cliped line to the appropriate external linestring
	for i, ipts := range cpts {
		for j, epts := range cps {
			if maths.Contains(epts, ipts[0], ipts[1]) {
				mapping = append(mapping, lineStringKeyMapping{Line: cpts[i], idx: j})
				break
			}
		}
	}
	return mapping
}

func clipPolygon(cr, scr clip.Region, geo tegola.Polygon) MultiPolygon {
	lines := geo.Sublines()
	var mp []Polygon
	var pts []float64

	// First linestring is the main extrior region
	for _, pt := range lines[0].Subpoints() {
		pts = append(pts, pt.X(), pt.Y())
	}
	// slice of slice of clipped points
	log.Println("Clipping Exteror line")
	scpts := clip.Polygon(cr, clip.Clockwise, pts)
	log.Println("Done Clipping Exteror line")
	for _, cpts := range scpts {
		mp = append(mp, Polygon{NewLine(cpts...)})
	}
	for i, inline := range lines[1:] {
		log.Println("Clipping Interior line", i)
		mapping := clipInteriorRegions(scr, scpts, inline)
		log.Println("Done Clipping Interior line", i)
		for _, mps := range mapping {
			mp[mps.idx] = append(mp[mps.idx], NewLine(mps.Line...))
		}
	}
	return MultiPolygon(mp)
}

func ClipGeometry(tile tegola.BoundingBox, extent int, geometry tegola.Geometry) (Geometry, error) {
	// halfExtent := extent
	quarterExtent := extent / 4
	/*
		tWidth := float64(int64(tile.Maxx - tile.Minx))
		tHeight := float64(int64(tile.Maxy - tile.Miny))

		clipBox := clipping.ClippingRegion{
			float64(-halfExtent),
			float64(+halfExtent),
			tWidth + float64(halfExtent),
			tHeight - float64(halfExtent),
		}
	*/
	clipBox := clip.Region{
		float64(0 - quarterExtent),
		float64(0 - quarterExtent),
		float64(extent + quarterExtent),
		float64(extent + quarterExtent),
	}
	subClipBox := clip.Region{
		clipBox[0] + float64(quarterExtent) - 1,
		clipBox[1] + float64(quarterExtent) - 1,
		clipBox[2] - float64(quarterExtent) + 1,
		clipBox[3] - float64(quarterExtent) + 1,
	}
	var pts []float64

	switch geo := geometry.(type) {
	default:
		return CloneGeometry(geometry)
	case tegola.Point:
		if clipBox.Contains(geo.X(), geo.Y()) {
			return CloneGeometry(geo)
		}
		log.Println("Point lies outside of the tile.", geo, clipBox)
		return CloneGeometry(geo)
		// return nil, nil
	case tegola.MultiPoint:
		var mp MultiPoint
		for _, pt := range geo.Points() {
			if clipBox.Contains(pt.X(), pt.Y()) {
				mp = append(mp, Point{pt.X(), pt.Y()})
			}
		}
		return mp, nil
	case tegola.LineString:
		for _, pt := range geo.Subpoints() {
			pts = append(pts, pt.X(), pt.Y())
		}
		cpts := clip.LineString(clipBox, pts)
		var mline MultiLine
		for _, l := range cpts {
			mline = append(mline, NewLine(l...))
		}
		return mline, nil
	case tegola.MultiLine:

		var mline MultiLine
		for _, line := range geo.Lines() {
			pts = pts[:0]
			for _, pt := range line.Subpoints() {
				pts = append(pts, pt.X(), pt.Y())
			}
			cpts := clip.LineString(clipBox, pts)
			for _, l := range cpts {
				mline = append(mline, NewLine(l...))
			}
		}
		return mline, nil
	case tegola.Polygon:
		return clipPolygon(clipBox, subClipBox, geo), nil
	case tegola.MultiPolygon:
		var mp MultiPolygon
		for _, p := range geo.Polygons() {
			cmp := clipPolygon(clipBox, subClipBox, p)
			mp = append(mp, cmp...)
		}
		return mp, nil
	}
}
