package geoip

// DB - GeoIP database with bound methods
type DB struct {
	file string
}

// Info - Basic IP information
type Info struct {
	CountryName string
	CountryCode string
	CityName    string
	Latitude    float64
	Longitude   float64
}
