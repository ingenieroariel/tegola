package spatialite_test

import (
	"os"
	"testing"

	"context"

	"github.com/terranodo/tegola"
	"github.com/terranodo/tegola/provider/spatialite"
)

func TestNewProvider(t *testing.T) {
	if os.Getenv("RUN_spatialite_TEST") != "yes" {
		return
	}

	testcases := []struct {
		config map[string]interface{}
	}{
		{
			config: map[string]interface{}{
				spatialite.ConfigKeyHost:     "localhost",
				spatialite.ConfigKeyPort:     int64(5432),
				spatialite.ConfigKeyDB:       "tegola",
				spatialite.ConfigKeyUser:     "postgres",
				spatialite.ConfigKeyPassword: "",
				spatialite.ConfigKeyLayers: []map[string]interface{}{
					{
						spatialite.ConfigKeyLayerName: "land",
						spatialite.ConfigKeyTablename: "ne_10m_land_scale_rank",
					},
				},
			},
		},
	}

	for i, tc := range testcases {
		_, err := spatialite.NewProvider(tc.config)
		if err != nil {
			t.Errorf("Failed test %v. Unable to create a new provider. err: %v", i, err)
			return
		}
	}
}

func TestMVTLayer(t *testing.T) {
	if os.Getenv("RUN_spatialite_TEST") != "yes" {
		return
	}

	testcases := []struct {
		config               map[string]interface{}
		tile                 tegola.Tile
		expectedFeatureCount int
	}{
		{
			config: map[string]interface{}{
				spatialite.ConfigKeyHost:     "localhost",
				spatialite.ConfigKeyPort:     int64(5432),
				spatialite.ConfigKeyDB:       "tegola",
				spatialite.ConfigKeyUser:     "postgres",
				spatialite.ConfigKeyPassword: "",
				spatialite.ConfigKeyLayers: []map[string]interface{}{
					{
						spatialite.ConfigKeyLayerName: "land",
						spatialite.ConfigKeyTablename: "ne_10m_land_scale_rank",
					},
				},
			},
			tile: tegola.Tile{
				Z: 1,
				X: 1,
				Y: 1,
			},
			expectedFeatureCount: 614,
		},
		//	scalerank test
		{
			config: map[string]interface{}{
				spatialite.ConfigKeyHost:     "localhost",
				spatialite.ConfigKeyPort:     int64(5432),
				spatialite.ConfigKeyDB:       "tegola",
				spatialite.ConfigKeyUser:     "postgres",
				spatialite.ConfigKeyPassword: "",
				spatialite.ConfigKeyLayers: []map[string]interface{}{
					{
						spatialite.ConfigKeyLayerName: "land",
						spatialite.ConfigKeySQL:       "SELECT gid, ST_AsBinary(geom) AS geom FROM ne_10m_land_scale_rank WHERE scalerank=!ZOOM! AND geom && !BBOX!",
					},
				},
			},
			tile: tegola.Tile{
				Z: 1,
				X: 1,
				Y: 1,
			},
			expectedFeatureCount: 23,
		},
		//	decode numeric(x,x) types
		{
			config: map[string]interface{}{
				spatialite.ConfigKeyHost:     "localhost",
				spatialite.ConfigKeyPort:     int64(5432),
				spatialite.ConfigKeyDB:       "tegola",
				spatialite.ConfigKeyUser:     "postgres",
				spatialite.ConfigKeyPassword: "",
				spatialite.ConfigKeyLayers: []map[string]interface{}{
					{
						spatialite.ConfigKeyLayerName:   "buildings",
						spatialite.ConfigKeyGeomIDField: "osm_id",
						spatialite.ConfigKeyGeomField:   "geometry",
						spatialite.ConfigKeySQL:         "SELECT ST_AsBinary(geometry) AS geometry, osm_id, name, nullif(as_numeric(height),-1) AS height, type FROM osm_buildings_test WHERE geometry && !BBOX!",
					},
				},
			},
			tile: tegola.Tile{
				Z: 16,
				X: 11241,
				Y: 26168,
			},
			expectedFeatureCount: 101,
		},
	}

	for i, tc := range testcases {
		p, err := spatialite.NewProvider(tc.config)
		if err != nil {
			t.Errorf("test (%v) failed. Unable to create a new provider. err: %v", i, err)
			return
		}

		//	iterate our configured layers
		for _, tcLayer := range tc.config[spatialite.ConfigKeyLayers].([]map[string]interface{}) {
			layerName := tcLayer[spatialite.ConfigKeyLayerName].(string)

			l, err := p.MVTLayer(context.Background(), layerName, &tc.tile, map[string]interface{}{})
			if err != nil {
				t.Errorf("test (%v) failed to create mvt layer err: %v", i, err)
				return
			}

			if len(l.Features()) != tc.expectedFeatureCount {
				t.Errorf("test (%v) failed.. expected feature count (%v), got (%v)", i, tc.expectedFeatureCount, len(l.Features()))
				return
			}
		}
	}
}
