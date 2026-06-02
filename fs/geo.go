package fs

import (
	"dotpip"
	"encoding/json"
	"fmt"
	yaml "github.com/goccy/go-yaml"
	"github.com/mmcloughlin/geohash"
	toml "github.com/pelletier/go-toml/v2"
	"math"
	"sort"
)

func (f *fileSystem) geospatialEncode(value map[string]dotpip.GeoLocation) (any, error) {
	switch f.encodeType {
	case JSON:
		return json.Marshal(value)
	case YAML:
		return yaml.Marshal(value)
	case TOML:
		// TOML needs a top-level table, and map[string]dotpip.GeoLocation is a top-level table
		return toml.Marshal(value)
	case RAW:
		return json.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *fileSystem) geospatialDecode(value any) (map[string]dotpip.GeoLocation, error) {
	finalValue := make(map[string]dotpip.GeoLocation)
	switch f.encodeType {
	case JSON:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case YAML:
		err := yaml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case TOML:
		err := toml.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	case RAW:
		err := json.Unmarshal(value.([]byte), &finalValue)
		return finalValue, err
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", f.encodeType)
	}
}

func (f *fileSystem) readGeo(key dotpip.Key) (map[string]dotpip.GeoLocation, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		return nil, err
	}
	if f.formatter.GeospatialDecode != nil {
		return f.formatter.GeospatialDecode(content)
	}
	return nil, fmt.Errorf("geospatial decoder not configured")
}

func (f *fileSystem) writeGeo(key dotpip.Key, geo map[string]dotpip.GeoLocation) error {
	if f.formatter.GeospatialEncode != nil {
		content, err := f.formatter.GeospatialEncode(geo)
		if err != nil {
			return err
		}
		return f.writeFileByKey(key, content.([]byte))
	}
	return fmt.Errorf("geospatial encoder not configured")
}

func (f *fileSystem) GeoAdd(key dotpip.Key, members []dotpip.GeoLocation, options ...dotpip.GeoAddOption) (int, error) {
	cmd := &dotpip.GeoAddCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	geoMap, err := f.readGeo(key)
	if err != nil {
		geoMap = make(map[string]dotpip.GeoLocation)
	}

	changed := 0
	added := 0

	for _, member := range members {
		if member.Longitude < -180 || member.Longitude > 180 || member.Latitude < -85.05112878 || member.Latitude > 85.05112878 {
			return 0, fmt.Errorf("ERR invalid longitude,latitude pair %f,%f", member.Longitude, member.Latitude)
		}

		existing, ok := geoMap[member.Name]
		if ok {
			if cmd.NX {
				continue
			}
			if existing.Longitude != member.Longitude || existing.Latitude != member.Latitude {
				geoMap[member.Name] = member
				changed++
			}
		} else {
			if cmd.XX {
				continue
			}
			geoMap[member.Name] = member
			added++
			changed++
		}
	}

	if changed > 0 || added > 0 {
		err = f.writeGeo(key, geoMap)
		if err != nil {
			return 0, err
		}
	}

	if cmd.CH {
		return changed, nil
	}
	return added, nil
}

func (f *fileSystem) GeoDist(key dotpip.Key, member1 string, member2 string, unit dotpip.GeoUnit) (float64, error) {
	geoMap, err := f.readGeo(key)
	if err != nil {
		return 0, err
	}

	loc1, ok1 := geoMap[member1]
	loc2, ok2 := geoMap[member2]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("ERR one or both members do not exist") // Redis returns nil, we return an error for simplicity
	}

	distKm := haversineDistance(loc1.Latitude, loc1.Longitude, loc2.Latitude, loc2.Longitude)

	switch unit {
	case dotpip.GeoUnitM:
		return distKm * 1000, nil
	case dotpip.GeoUnitKM:
		return distKm, nil
	case dotpip.GeoUnitMI:
		return distKm * 0.621371, nil
	case dotpip.GeoUnitFT:
		return distKm * 3280.84, nil
	default:
		return 0, fmt.Errorf("ERR unsupported unit provided")
	}
}

func (f *fileSystem) GeoHash(key dotpip.Key, members ...string) ([]string, error) {
	geoMap, err := f.readGeo(key)
	if err != nil {
		geoMap = make(map[string]dotpip.GeoLocation)
	}

	var hashes []string
	for _, member := range members {
		loc, ok := geoMap[member]
		if !ok {
			hashes = append(hashes, "") // Redis returns nil, we use empty string
		} else {
			// Redis uses an 11-character geohash string
			hashes = append(hashes, geohash.EncodeWithPrecision(loc.Latitude, loc.Longitude, 11))
		}
	}

	return hashes, nil
}

func (f *fileSystem) GeoPos(key dotpip.Key, members ...string) ([]*dotpip.GeoLocation, error) {
	geoMap, err := f.readGeo(key)
	if err != nil {
		geoMap = make(map[string]dotpip.GeoLocation)
	}

	var positions []*dotpip.GeoLocation
	for _, member := range members {
		loc, ok := geoMap[member]
		if !ok {
			positions = append(positions, nil)
		} else {
			positions = append(positions, &loc)
		}
	}

	return positions, nil
}

func extractMeters(dist float64, unit dotpip.GeoUnit) float64 {
	switch unit {
	case dotpip.GeoUnitM:
		return dist
	case dotpip.GeoUnitKM:
		return dist * 1000
	case dotpip.GeoUnitMI:
		return dist * 1609.34
	case dotpip.GeoUnitFT:
		return dist * 0.3048
	default:
		return dist
	}
}

func convertFromKmToUnit(distKm float64, unit dotpip.GeoUnit) float64 {
	switch unit {
	case dotpip.GeoUnitM:
		return distKm * 1000
	case dotpip.GeoUnitKM:
		return distKm
	case dotpip.GeoUnitMI:
		return distKm * 0.621371
	case dotpip.GeoUnitFT:
		return distKm * 3280.84
	default:
		return distKm
	}
}

func (f *fileSystem) GeoSearch(key dotpip.Key, options ...dotpip.GeoSearchOption) ([]dotpip.GeoSearchResult, error) {
	cmd := &dotpip.GeoSearchCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	geoMap, err := f.readGeo(key)
	if err != nil {
		return nil, nil // Or error? Redis returns empty array if key does not exist
	}

	var centerLat, centerLon float64
	if cmd.UseLonLat {
		centerLat = cmd.FromLatitude
		centerLon = cmd.FromLongitude
	} else if cmd.FromMember != "" {
		loc, ok := geoMap[cmd.FromMember]
		if !ok {
			return nil, fmt.Errorf("ERR could not decode requested zset member")
		}
		centerLat = loc.Latitude
		centerLon = loc.Longitude
	} else {
		return nil, fmt.Errorf("ERR either FROMMEMBER or FROMLONLAT must be provided")
	}

	var results []dotpip.GeoSearchResult

	for name, loc := range geoMap {
		distKm := haversineDistance(centerLat, centerLon, loc.Latitude, loc.Longitude)

		// This is a naive bounding box implementation by just using the distance in km
		// A real implementation would project the bounding box correctly.
		// For the sake of simplicity, we approximate it here or use distance directly.

		match := false
		var distToReturn float64

		if cmd.UseRadius {
			radiusMeters := extractMeters(cmd.ByRadius, cmd.RadiusUnit)
			if distKm*1000 <= radiusMeters {
				match = true
				distToReturn = convertFromKmToUnit(distKm, cmd.RadiusUnit)
			}
		} else if cmd.UseBox {
			// A highly simplified bounding box check. True bounding box requires more complex math (like taking longitude scaling by latitude into account).
			// We approximate by converting width/height to km and checking Haversine distances to edges.

			// Let's do a simple check.
			widthMeters := extractMeters(cmd.ByBoxWidth, cmd.BoxUnit)
			heightMeters := extractMeters(cmd.ByBoxHeight, cmd.BoxUnit)

			// Approximation: 1 degree latitude is approx 111km.
			// Longitude degree distance varies: cos(lat) * 111km.

			latDiff := math.Abs(loc.Latitude - centerLat)
			lonDiff := math.Abs(loc.Longitude - centerLon)

			latDistMeters := latDiff * 111320.0
			lonDistMeters := lonDiff * 40075000.0 * math.Cos(degreesToRadians(centerLat)) / 360.0

			if latDistMeters <= heightMeters/2 && lonDistMeters <= widthMeters/2 {
				match = true
				distToReturn = convertFromKmToUnit(distKm, cmd.BoxUnit)
			}
		} else {
			return nil, fmt.Errorf("ERR either BYRADIUS or BYBOX must be provided")
		}

		if match {
			res := dotpip.GeoSearchResult{
				Name: name,
			}
			if cmd.WithDist {
				res.Distance = distToReturn
			} else {
				// We need distance for sorting, but use a negative value or specific logic if needed.
				// However, Redis sorting relies on distance. We'll store it here so we can sort,
				// but in a real system we'd strip this before returning to the client if WithDist is false.
				// For the sake of dotpip's API design, we can return it and the caller handles presentation.
				res.Distance = distToReturn
			}

			if cmd.WithCoord {
				res.Latitude = loc.Latitude
				res.Longitude = loc.Longitude
			}
			if cmd.WithHash {
				res.Hash = geohash.EncodeWithPrecision(loc.Latitude, loc.Longitude, 11)
			}

			results = append(results, res)
		}
	}

	if cmd.Asc {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Distance < results[j].Distance
		})
	} else if cmd.Desc {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Distance > results[j].Distance
		})
	}

	if cmd.Count > 0 && len(results) > cmd.Count {
		results = results[:cmd.Count]
	}

	return results, nil
}

func (f *fileSystem) GeoSearchStore(destination dotpip.Key, source dotpip.Key, searchOptions []dotpip.GeoSearchOption, storeOptions ...dotpip.GeoSearchStoreOption) (int, error) {
	// First, run the search but make sure we get the distance and coordinates
	// since we need them to store them as a ZSet.

	// Add required options to get what we need internally
	searchOptions = append(searchOptions, dotpip.WithGeoSearchWithDist(), dotpip.WithGeoSearchWithCoord())

	results, err := f.GeoSearch(source, searchOptions...)
	if err != nil {
		return 0, err
	}

	storeCmd := &dotpip.GeoSearchStoreCommand{}
	for _, opt := range storeOptions {
		opt(storeCmd)
	}

	// In Redis, GeoSearchStore saves results into a ZSET.
	var zmembers []dotpip.Z
	for _, res := range results {
		score := res.Distance
		if !storeCmd.StoreDist {
			// If not STOREDIST, it stores the internal 52-bit geohash as the score
			hashInt := geohash.EncodeIntWithPrecision(res.Latitude, res.Longitude, 52)
			score = float64(hashInt)
		}
		zmembers = append(zmembers, dotpip.Z{
			Member: res.Name,
			Score:  score,
		})
	}

	// Remove the old key entirely if we are overwriting it (Redis behavior)
	f.removeFileByKey(destination)

	if len(zmembers) > 0 {
		_, err = f.ZAdd(destination, zmembers)
		if err != nil {
			return 0, err
		}
	}

	return len(results), nil
}
