package spatialite

import "github.com/terranodo/tegola"

// layer holds information about a query.
type Layer struct {
	// The Name of the layer
	name string
	// The SQL to use when querying spatialite for this layer
	sql string
	// The ID field name, this will default to 'gid' if not set to something other then empty string.
	idField string
	// The Geometery field name, this will default to 'geom' if not set to soemthing other then empty string.
	geomField string
	// GeomType is the the type of geometry returned from the SQL
	geomType tegola.Geometry
	// The SRID that the data in the table is stored in. This will default to WebMercator
	srid int
}

func (l Layer) Name() string {
	return l.name
}

func (l Layer) GeomType() tegola.Geometry {
	return l.geomType
}

func (l Layer) SRID() int {
	return l.srid
}

func (l Layer) GeomFieldName() string {
	return l.geomField
}

func (l Layer) IDFieldName() string {
	return l.idField
}
