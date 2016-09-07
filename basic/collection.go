package basic

import "github.com/terranodo/tegola"

// Geometry represents a geomentry of the basic type.
type Geometry interface {
	basicType() // does nothing, but here to make Geometry for basic types unique.
}

// Collection type can represent one or more other basic types.
type Collection []Geometry

//Geometeries return a set of geometeies that make that collection.
func (c Collection) Geometeries() (geometeries []tegola.Geometry) {
	geometeries = make([]tegola.Geometry, 0, len(c))
	for i := range c {
		geometeries = append(geometeries, c[i])
	}
	return geometeries
}

func (Collection) String() string {
	return "Collection"
}

// private this is for membership to basic types.
func (Collection) basicType() {}
