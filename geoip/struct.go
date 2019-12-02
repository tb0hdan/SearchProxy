package geoip

type DB struct {
	file string
}

type Info struct {
	CountryName string
	CountryCode string
	CityName    string
	Latitude    float64
	Longitude   float64
}
