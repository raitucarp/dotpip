package dotpip

type GeoLocation struct {
	Longitude float64
	Latitude  float64
	Name      string
}

type GeoAddCommand struct {
	NX bool
	XX bool
	CH bool
}

type GeoAddOption func(*GeoAddCommand)

func WithGeoAddNX() GeoAddOption {
	return func(cmd *GeoAddCommand) {
		cmd.NX = true
	}
}

func WithGeoAddXX() GeoAddOption {
	return func(cmd *GeoAddCommand) {
		cmd.XX = true
	}
}

func WithGeoAddCH() GeoAddOption {
	return func(cmd *GeoAddCommand) {
		cmd.CH = true
	}
}

type GeoUnit string

const (
	GeoUnitM  GeoUnit = "m"
	GeoUnitKM GeoUnit = "km"
	GeoUnitMI GeoUnit = "mi"
	GeoUnitFT GeoUnit = "ft"
)

type GeoSearchResult struct {
	Name      string
	Distance  float64
	Longitude float64
	Latitude  float64
	Hash      string
}

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

type GeoSearchOption func(*GeoSearchCommand)

func WithGeoSearchFromMember(member string) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.FromMember = member
	}
}

func WithGeoSearchFromLonLat(longitude, latitude float64) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.FromLongitude = longitude
		cmd.FromLatitude = latitude
		cmd.UseLonLat = true
	}
}

func WithGeoSearchByRadius(radius float64, unit GeoUnit) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.ByRadius = radius
		cmd.RadiusUnit = unit
		cmd.UseRadius = true
	}
}

func WithGeoSearchByBox(width, height float64, unit GeoUnit) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.ByBoxWidth = width
		cmd.ByBoxHeight = height
		cmd.BoxUnit = unit
		cmd.UseBox = true
	}
}

func WithGeoSearchAsc() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.Asc = true
	}
}

func WithGeoSearchDesc() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.Desc = true
	}
}

func WithGeoSearchCount(count int, any bool) GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.Count = count
		cmd.Any = any
	}
}

func WithGeoSearchWithCoord() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.WithCoord = true
	}
}

func WithGeoSearchWithDist() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.WithDist = true
	}
}

func WithGeoSearchWithHash() GeoSearchOption {
	return func(cmd *GeoSearchCommand) {
		cmd.WithHash = true
	}
}

type GeoSearchStoreCommand struct {
	StoreDist bool
}

type GeoSearchStoreOption func(*GeoSearchStoreCommand)

func WithGeoSearchStoreDist() GeoSearchStoreOption {
	return func(cmd *GeoSearchStoreCommand) {
		cmd.StoreDist = true
	}
}
