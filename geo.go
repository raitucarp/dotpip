package dotpip

// GeoLocation represents a geo location.
type GeoLocation struct {
	Longitude float64
	Latitude  float64
	Name      string
}

// GeoAddCommand options for GEOADD.
type GeoAddCommand struct {
	NX bool
	XX bool
	CH bool
}

// GeoAddOption configures GeoAddCommand.
type GeoAddOption func(*GeoAddCommand)

// WithGeoAddNX sets NX for GEOADD.
func WithGeoAddNX() GeoAddOption {
	return func(cmd *GeoAddCommand) {
		cmd.NX = true
	}
}

// WithGeoAddXX sets XX for GEOADD.
func WithGeoAddXX() GeoAddOption {
	return func(cmd *GeoAddCommand) {
		cmd.XX = true
	}
}

// WithGeoAddCH sets CH for GEOADD.
func WithGeoAddCH() GeoAddOption {
	return func(cmd *GeoAddCommand) {
		cmd.CH = true
	}
}

// GeoUnit represents a geo unit.
type GeoUnit string

const (
	// GeoUnitM represents meters.
	GeoUnitM  GeoUnit = "m"
	// GeoUnitKM represents kilometers.
	GeoUnitKM GeoUnit = "km"
	// GeoUnitMI represents miles.
	GeoUnitMI GeoUnit = "mi"
	// GeoUnitFT represents feet.
	GeoUnitFT GeoUnit = "ft"
)

// GeoSearchResult represents a result from GEOSEARCH.
type GeoSearchResult struct {
	Name      string
	Distance  float64
	Longitude float64
	Latitude  float64
	Hash      string
}

// GeoSearchCommand options for GEOSEARCH.
type GeoSearchCommand struct {
	// FROMMEMBER
	FromMember string
	// FROMLONLAT
	FromLongitude float64
	FromLatitude  float64
	UseLonLat     bool

	// BYRADIUS
	ByRadius   float64
	RadiusUnit GeoUnit
	UseRadius  bool

	// BYBOX
	ByBoxWidth  float64
	ByBoxHeight float64
	BoxUnit     GeoUnit
	UseBox      bool

	// Order
	Asc  bool
	Desc bool

	// Count
	Count int
	Any   bool

	// Return options
	WithCoord bool
	WithDist  bool
	WithHash  bool
}

// GeoSearchOption configures GeoSearchCommand.
type GeoSearchOption func(*GeoSearchCommand)

// WithGeoSearchFromMember sets FROMMEMBER for GEOSEARCH.
func WithGeoSearchFromMember(member string) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.FromMember = member
	}
}

// WithGeoSearchFromLonLat sets FROMLONLAT for GEOSEARCH.
func WithGeoSearchFromLonLat(longitude, latitude float64) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.FromLongitude = longitude
		cmd.FromLatitude = latitude
		cmd.UseLonLat = true
	}
}

// WithGeoSearchByRadius sets BYRADIUS for GEOSEARCH.
func WithGeoSearchByRadius(radius float64, unit GeoUnit) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.ByRadius = radius
		cmd.RadiusUnit = unit
		cmd.UseRadius = true
	}
}

// WithGeoSearchByBox sets BYBOX for GEOSEARCH.
func WithGeoSearchByBox(width, height float64, unit GeoUnit) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.ByBoxWidth = width
		cmd.ByBoxHeight = height
		cmd.BoxUnit = unit
		cmd.UseBox = true
	}
}

// WithGeoSearchAsc sets ASC for GEOSEARCH.
func WithGeoSearchAsc() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.Asc = true
	}
}

// WithGeoSearchDesc sets DESC for GEOSEARCH.
func WithGeoSearchDesc() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.Desc = true
	}
}

// WithGeoSearchCount sets COUNT for GEOSEARCH.
func WithGeoSearchCount(count int, anyCount bool) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.Count = count
		cmd.Any = anyCount
	}
}

// WithGeoSearchWithCoord sets WITHCOORD for GEOSEARCH.
func WithGeoSearchWithCoord() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.WithCoord = true
	}
}

// WithGeoSearchWithDist sets WITHDIST for GEOSEARCH.
func WithGeoSearchWithDist() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.WithDist = true
	}
}

// WithGeoSearchWithHash sets WITHHASH for GEOSEARCH.
func WithGeoSearchWithHash() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.WithHash = true
	}
}

// GeoSearchStoreCommand options for GEOSEARCHSTORE.
type GeoSearchStoreCommand struct {
	StoreDist bool
}

// GeoSearchStoreOption configures GeoSearchStoreCommand.
type GeoSearchStoreOption func(*GeoSearchStoreCommand)

// WithGeoSearchStoreDist sets STOREDIST for GEOSEARCHSTORE.
func WithGeoSearchStoreDist() GeoSearchStoreOption {
	return func(cmd *GeoSearchStoreCommand) {
		cmd.StoreDist = true
	}
}
