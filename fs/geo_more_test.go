package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeoDistOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_geo_dist_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mygeo")
	_, _ = dotfs.GeoAdd(key, []dotpip.GeoLocation{{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"}, {Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"}})

	// Dist in m
	distM, err := dotfs.GeoDist(key, "Palermo", "Catania", dotpip.GeoUnitM)
	assert.NoError(t, err)
	assert.Greater(t, distM, 0.0)

	// Dist in km
	distKM, err := dotfs.GeoDist(key, "Palermo", "Catania", dotpip.GeoUnitKM)
	assert.NoError(t, err)
	assert.Greater(t, distKM, 0.0)

	// Dist in ft
	distFT, err := dotfs.GeoDist(key, "Palermo", "Catania", dotpip.GeoUnitFT)
	assert.NoError(t, err)
	assert.Greater(t, distFT, 0.0)

	// Dist in mi
	distMI, err := dotfs.GeoDist(key, "Palermo", "Catania", dotpip.GeoUnitMI)
	assert.NoError(t, err)
	assert.Greater(t, distMI, 0.0)

	// Missing member
	_, err = dotfs.GeoDist(key, "Palermo", "Rome", dotpip.GeoUnitM)
	assert.Error(t, err)
}

func TestGeoSearchMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_geo_search_more_test_")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mygeo")
	_, _ = dotfs.GeoAdd(key, []dotpip.GeoLocation{
		{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
		{Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"},
	})

	// BYRADIUS with KM
	resRadius, err := dotfs.GeoSearch(key,
		dotpip.WithGeoSearchFromMember("Palermo"),
		dotpip.WithGeoSearchByRadius(200, dotpip.GeoUnitKM),
		dotpip.WithGeoSearchWithDist(),
	)
	assert.NoError(t, err)
	assert.Len(t, resRadius, 2)

	// BYRADIUS with FT
	resRadius2, err := dotfs.GeoSearch(key,
		dotpip.WithGeoSearchFromMember("Palermo"),
		dotpip.WithGeoSearchByRadius(2000000, dotpip.GeoUnitFT),
		dotpip.WithGeoSearchWithDist(),
	)
	assert.NoError(t, err)
	assert.Len(t, resRadius2, 2)

	// BYRADIUS with MI
	resRadius3, err := dotfs.GeoSearch(key,
		dotpip.WithGeoSearchFromMember("Palermo"),
		dotpip.WithGeoSearchByRadius(200, dotpip.GeoUnitMI),
		dotpip.WithGeoSearchWithDist(),
	)
	assert.NoError(t, err)
	assert.Len(t, resRadius3, 2)

	// BYBOX with M
	resBox, err := dotfs.GeoSearch(key,
		dotpip.WithGeoSearchFromMember("Palermo"),
		dotpip.WithGeoSearchByBox(400000, 400000, dotpip.GeoUnitM),
	)
	assert.NoError(t, err)
	assert.Len(t, resBox, 2)
}
