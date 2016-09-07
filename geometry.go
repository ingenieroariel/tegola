// Package tegola describes the basic geometeries that can be used to convert to
// and from.
package tegola

// Geometry describes a geometry.
type Geometry interface{}

// Point is how a point should look like.
type Point interface {
	Geometry
	X() float64
	Y() float64
}

// Point3 is a point with three dimensions; at current is just converted and treated as a point.
type Point3 interface {
	Point
	Z() float64
}

// MultiPoint is a Geometry with multiple individual points.
type MultiPoint interface {
	Geometry
	Points() []Point
}

// LineString is a Geometry of a line.
type LineString interface {
	Geometry
	Subpoints() []Point
}

// MultiLine is a Geometry with multiple individual lines.
type MultiLine interface {
	Geometry
	Lines() []LineString
}

// Polygon is a multi-line Geometry  where the first linestring is the exterior ring.
// The following linestrings will be interior rings that substract the area from the exterior ring
type Polygon interface {
	Geometry
	Sublines() []LineString
}

// MultiPolygon describes a Geometry multiple polygons.
type MultiPolygon interface {
	Geometry
	Polygons() []Polygon
}

// Collection is a collections of different geometries.
type Collection interface {
	Geometry
	Geometries() []Geometry
}
