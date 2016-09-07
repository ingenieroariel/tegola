package basic

import (
	"github.com/terranodo/tegola"
	"github.com/terranodo/tegola/basic/clipping"
)

type lineStringKeyMapping struct {
	Line []float64
	idx  int
}

func clipInteriorRegions(scr clipping.ClippingRegion, cps [][]float64, inline tegola.Line) []lineStringKeyMapping {
	var pts []float64
	// Convert the linestring into a slice of points.
	for _, pt := range inline.Subpoints() {
		pts = append(pts, pt.X(), pt.Y())
	}
	cpts := clipping.ClipPolygon(scr, clipping.CounterClockwise, pts)
	if len(cpts) == 0 {
		return nil
	}
	mapping = make([]lineStringKeyMapping, 0, len(cpts))
	// Now we need to map the internal cliped line to the appropriate external linestring
	for i, ipts := range cpts {
		for j, epts := range cps {
			if clipping.Contains(ipts[:2], epts) {
				mapping = append(mapping, lineStringKeyMapping{Line: cpts[i], idx: j})
				break
			}
		}
	}
	return mapping
}

func clipPolygon(cr, scr clipping.ClippingRegion, geo tegola.Polygon) MultiPolygon {
	lines := geo.Sublines()
	var mlines []Polygon

	pts = pts[:0]
	// First linestring is the main extrior region
	for _, pt := range lines[0].Subpoints() {
		pts = append(pts, pt.X(), pt.Y())
	}
	cpts := clipping.ClipPolygon(cr, clipping.Clockwise, pts)
	for i, cpt := range cpts {
		l := []Polygon{NewLine(cpt)}
		mlines = append(mlines, l)
	}
	for inLine := range lines[1:] {
		mapping := clipInteriorRegions(scr, cpts, inline)
		for _, mps := range mapping {
			mlines[mps.idx] = append(mlines[mps.idx], NewLine(mps.Line))
		}
	}
	return MultiPolygon(mlines)
}

func ClipGeometry(tile tegola.BoundingBox, extent float64, geometry tegola.Geometry) (tegola.Geometry, error) {
	halfExtent := extent / 2
	quarterExtent := extent / 4
	clipBox := clipping.ClippingRegion{
		tile.Minx - halfExtent,
		tile.Miny - halfExtent,
		tile.Maxx + halfExtent,
		tile.Maxy + halfExtent,
	}
	subClipBox := clipping.ClippingRegion{
		clipBox[0] - quarterExtent,
		clipBox[1] - quarterExtent,
		clipBox[2] - quarterExtnet,
		clipbox[3] - quarterExtent,
	}
	var pts []float64

	switch geo := geometry.(type) {
	default:
		return geometry, nil
	case tegola.LineString:
		for _, pt := range geo.Subpoints() {
			pts = append(pts, pt.X(), pt.Y())
		}
		cpts := clipping.ClipLineString(clipBox, pts)
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
			cpts := clipping.ClipLineString(clipBox, pts)
			for _, l := range cpts {
				mline = append(mline, NewLine(l...))
			}
		}
		return nline, nil
	case tegola.Polygon:
		return clipPolygon(clipBox, subClipBox, geo)

	}

}
