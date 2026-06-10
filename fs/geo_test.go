package fs

import (
	"dotpip"
	"math"
	"os"
	"testing"
)

func assertGeoEqual(t *testing.T, expected, actual float64, tolerance float64, msg string) {
	if math.Abs(expected-actual) > tolerance {
		t.Errorf("%s: expected %f, got %f", msg, expected, actual)
	}
}

func TestGeoJSON(t *testing.T) {
	runGeoTests(t, JSON)
}

func TestGeoYAML(t *testing.T) {
	runGeoTests(t, YAML)
}

func TestGeoTOML(t *testing.T) {
	runGeoTests(t, TOML)
}

func TestGeoRAW(t *testing.T) {
	runGeoTests(t, RAW)
}

func runGeoTests(t *testing.T, encodeType FileEncodeType) {
	t.Helper()

	testDir := "test_geo_data_" + string(encodeType)
	os.RemoveAll(testDir)
	_ = os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	fs := NewFileSystem(testDir + "/")
	fs.EncodeType(encodeType)
	defer fs.Close()

	key := dotpip.NewKey("Sicily")

	// GeoSearch on empty key should return empty array
	emptyResults, err := fs.GeoSearch(dotpip.NewKey("EmptyGeo"), dotpip.WithGeoSearchFromLonLat(15.0, 37.5), dotpip.WithGeoSearchByBox(50, 50, dotpip.GeoUnitKM))
	if err != nil {
		t.Fatalf("GeoSearch empty key failed: %v", err)
	}
	if len(emptyResults) != 0 || emptyResults == nil {
		t.Fatalf("Expected empty array for missing key, got %v", emptyResults)
	}

	// GeoAdd
	added, err := fs.GeoAdd(key, []dotpip.GeoLocation{
		{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
		{Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"},
	})
	if err != nil {
		t.Fatalf("GeoAdd failed: %v", err)
	}
	if added != 2 {
		t.Errorf("Expected to add 2 elements, added %d", added)
	}

	// GeoAdd XX
	added, err = fs.GeoAdd(key, []dotpip.GeoLocation{
		{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
		{Longitude: 12.0, Latitude: 40.0, Name: "Rome"},
	}, dotpip.WithGeoAddXX())
	if err != nil {
		t.Fatalf("GeoAdd XX failed: %v", err)
	}
	if added != 0 {
		t.Errorf("Expected to add 0 elements with XX, added %d", added)
	}

	// GeoAdd NX
	added, err = fs.GeoAdd(key, []dotpip.GeoLocation{
		{Longitude: 13.0, Latitude: 38.0, Name: "Palermo"},
		{Longitude: 12.0, Latitude: 40.0, Name: "Rome"},
	}, dotpip.WithGeoAddNX())
	if err != nil {
		t.Fatalf("GeoAdd NX failed: %v", err)
	}
	if added != 1 {
		t.Errorf("Expected to add 1 element with NX, added %d", added)
	}

	// GeoDist
	dist, err := fs.GeoDist(key, "Palermo", "Catania", dotpip.GeoUnitKM)
	if err != nil {
		t.Fatalf("GeoDist failed: %v", err)
	}
	assertGeoEqual(t, 166.274, dist, 0.1, "GeoDist km")

	// GeoHash
	hashes, err := fs.GeoHash(key, "Palermo", "Catania", "Unknown")
	if err != nil {
		t.Fatalf("GeoHash failed: %v", err)
	}
	if len(hashes) != 3 {
		t.Fatalf("Expected 3 hashes, got %d", len(hashes))
	}
	if hashes[0] != "sqc8b49rnyt" {
		t.Errorf("Expected hash for Palermo to be sqc8b49rnyt, got %s", hashes[0])
	}
	if hashes[2] != "" {
		t.Errorf("Expected empty hash for Unknown, got %s", hashes[2])
	}

	// GeoPos
	positions, err := fs.GeoPos(key, "Palermo", "Unknown")
	if err != nil {
		t.Fatalf("GeoPos failed: %v", err)
	}
	if len(positions) != 2 {
		t.Fatalf("Expected 2 positions, got %d", len(positions))
	}
	if positions[0] == nil {
		t.Fatalf("Expected position for Palermo, got nil")
	}
	assertGeoEqual(t, 13.361389, positions[0].Longitude, 0.0001, "Palermo Longitude")
	assertGeoEqual(t, 38.115556, positions[0].Latitude, 0.0001, "Palermo Latitude")
	if positions[1] != nil {
		t.Errorf("Expected nil position for Unknown, got %+v", positions[1])
	}
}

func TestGeoSearchAndStore(t *testing.T) {
	runGeoSearchTests(t, JSON)
}

func runGeoSearchTests(t *testing.T, encodeType FileEncodeType) {
	t.Helper()

	testDir := "test_geosearch_data_" + string(encodeType)
	os.RemoveAll(testDir)
	_ = os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	fs := NewFileSystem(testDir + "/")
	fs.EncodeType(encodeType)
	defer fs.Close()

	key := dotpip.NewKey("Sicily")

	_, err := fs.GeoAdd(key, []dotpip.GeoLocation{
		{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
		{Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"},
		{Longitude: 13.583333, Latitude: 37.316667, Name: "Agrigento"}, // Agrigento is closer to Palermo than Catania
	})
	if err != nil {
		t.Fatalf("GeoAdd failed: %v", err)
	}

	// 1. GeoSearch FromMember ByRadius
	results, err := fs.GeoSearch(key, dotpip.WithGeoSearchFromMember("Palermo"), dotpip.WithGeoSearchByRadius(200, dotpip.GeoUnitKM), dotpip.WithGeoSearchWithDist(), dotpip.WithGeoSearchAsc())
	if err != nil {
		t.Fatalf("GeoSearch failed: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}
	if results[0].Name != "Palermo" || results[0].Distance != 0 {
		t.Errorf("Expected Palermo first with 0 distance, got %s %f", results[0].Name, results[0].Distance)
	}
	if results[1].Name != "Agrigento" {
		t.Errorf("Expected Agrigento second, got %s", results[1].Name)
	}
	if results[2].Name != "Catania" {
		t.Errorf("Expected Catania third, got %s", results[2].Name)
	}

	// 2. GeoSearch FromLonLat ByBox
	results, err = fs.GeoSearch(key, dotpip.WithGeoSearchFromLonLat(15.0, 37.5), dotpip.WithGeoSearchByBox(50, 50, dotpip.GeoUnitKM))
	if err != nil {
		t.Fatalf("GeoSearch Box failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result in box, got %d", len(results))
	}
	if results[0].Name != "Catania" {
		t.Errorf("Expected Catania in box, got %s", results[0].Name)
	}

	// 3. GeoSearchStore
	destKey := dotpip.NewKey("SicilySubset")
	added, err := fs.GeoSearchStore(destKey, key, []dotpip.GeoSearchOption{
		dotpip.WithGeoSearchFromMember("Palermo"),
		dotpip.WithGeoSearchByRadius(100, dotpip.GeoUnitKM),
	}, dotpip.WithGeoSearchStoreDist())
	if err != nil {
		t.Fatalf("GeoSearchStore failed: %v", err)
	}

	// Should find Palermo and Agrigento (Palermo dist=0, Agrigento dist~90km)
	if added != 2 {
		t.Fatalf("Expected 2 elements stored, got %d", added)
	}

	// Check if ZSET is properly populated
	count, err := fs.ZCard(destKey)
	if err != nil {
		t.Fatalf("ZCard failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected ZSET size 2, got %d", count)
	}
}
