package geoip

import (
	net2 "searchproxy/util/network"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
	"github.com/umahmood/haversine"

	"net"
)

// GeoIPLookupIP - lookup IP in GeoIP database and return information
func (gdb *DB) GeoIPLookupIP(ip string) (info *Info, err error) {
	var (
		ok bool
		countryName,
		cityName string
	)

	db, err := geoip2.Open(gdb.file)

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

	info = &Info{
		CountryName: countryName,
		CountryCode: record.Country.IsoCode,
		CityName:    cityName,
		Latitude:    record.Location.Latitude,
		Longitude:   record.Location.Longitude,
	}

	return info, nil
}

// LookupIP - helper function that looks up IP and returns basic info
func (gdb *DB) LookupIP(ip string) (geo *Info, err error) {
	geo, err = gdb.GeoIPLookupIP(ip)
	if err != nil {
		log.Printf("Lookup failed with: %v\n", err)
		return nil, err
	}

	return geo, nil
}

// LookupDomain - helper function that translates domain to IP and returns basic info
func (gdb *DB) LookupDomain(domain string) (*Info, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Println("DNS resolve error")
		return nil, err
	}

	return gdb.LookupIP(ips[0].String())
}

// LookupURL - helper function that translates domain URL to domain to IP and returns basic info
func (gdb *DB) LookupURL(rurl string) (*Info, error) {
	host, err := net2.LookupHostByURL(rurl)
	if err != nil {
		return nil, err
	}

	return gdb.LookupDomain(host)
}

// DistanceIPLatLon - measure distance between IP and supplied coordinates
func (gdb *DB) DistanceIPLatLon(ip string, lat, lon float64) (distance float64, err error) {
	info, err := gdb.GeoIPLookupIP(ip)
	if err != nil {
		return -1, err
	}

	return gdb.DistanceLatLon(info.Latitude, info.Longitude, lat, lon), nil
}

// DistanceIP - measure distance between two IPs, in km
func (gdb *DB) DistanceIP(ip1, ip2 string) (distance float64, err error) {
	info1, err := gdb.GeoIPLookupIP(ip1)
	if err != nil {
		return -1, err
	}

	info2, err := gdb.GeoIPLookupIP(ip2)
	if err != nil {
		return -1, err
	}

	return gdb.DistanceLatLon(info1.Latitude, info1.Longitude,
		info2.Latitude, info2.Longitude), nil
}

// DistanceLatLon - measure distance between two geographical points, in km
func (gdb *DB) DistanceLatLon(lat1, lon1, lat2, lon2 float64) (distance float64) {
	point1 := haversine.Coord{Lat: lat1, Lon: lon1}
	point2 := haversine.Coord{Lat: lat2, Lon: lon2}
	_, distance = haversine.Distance(point1, point2)

	return distance
}

// New - return populated instance of GeoIP DB
func New(file string) *DB {
	return &DB{file: file}
}
