package geoip

import (
	log "github.com/sirupsen/logrus"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPInfo struct {
	CountryName string
	CountryCode string
	CityName string
	Latitude float64
	Longitude float64
}

func GeoIPLookupIP(ip string) (info *GeoIPInfo, err error){
	var (
		ok bool
		countryName,
		cityName string
	)
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Printf("GeoIP DB open failed with: %v\n", err)
		return nil, err
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ipS := net.ParseIP(ip)
	record, err := db.City(ipS)
	if err != nil {
		log.Printf("GeoIP IP lookup failed with: %v\n", err)
		return nil, err
	}

	if countryName, ok = record.Country.Names["en"]; !ok {
		countryName = "Unknown"
	}

	if cityName, ok = record.City.Names["en"]; !ok {
		cityName = "Unknown"
	}

	info = &GeoIPInfo{
		CountryName: countryName,
		CountryCode: record.Country.IsoCode,
		CityName: cityName,
		Latitude: record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}
	return
}

