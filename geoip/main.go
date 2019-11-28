package geoip

import (
	log "github.com/sirupsen/logrus"

	"net"
	"net/url"
)

type GeoIPDB struct {

}

func (gdb *GeoIPDB) LookupIP(ip string) (geo *GeoIPInfo, err error){
	geo, err = GeoIPLookupIP(ip)
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
