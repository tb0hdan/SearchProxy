package geoip

import (
	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"

	"net"
	"net/url"
)

func (gdb *GeoIPDB) GeoIPLookupIP(ip string) (info *GeoIPInfo, err error){
	var (
		ok bool
		countryName,
		cityName string
	)
	db, err := geoip2.Open(gdb.DatabaseFile)
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

func (gdb *GeoIPDB) LookupIP(ip string) (geo *GeoIPInfo, err error){
	geo, err = gdb.GeoIPLookupIP(ip)
	if err != nil {
		log.Printf("Lookup failed with: %v\n", err)
		return nil, err
	}
	return geo, nil
}

func (gdb *GeoIPDB) LookupDomain(domain string) (*GeoIPInfo, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Println("DNS resolve error")
		return nil, err
	}
	return gdb.LookupIP(ips[0].String())
}

func (gdb *GeoIPDB) LookupURL(rurl string) (*GeoIPInfo, error) {
	parsed, err := url.Parse(rurl)
	if err != nil {
		log.Println("URL parse error")
		return nil, err
	}
	return gdb.LookupDomain(parsed.Host)
}

func New(DatabaseFile string) *GeoIPDB{
	return &GeoIPDB{DatabaseFile:DatabaseFile}
}
