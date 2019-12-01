package geoip

type GeoIPDB struct {
	DatabaseFile string
}

type GeoIPInfo struct {
	CountryName string
	CountryCode string
	CityName string
	Latitude float64
	Longitude float64
}

