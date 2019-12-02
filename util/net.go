package util

import (
	"net"

	log "github.com/sirupsen/logrus"

	"net/url"
)

func LookupHostByURL(rurl string) (host string, err error){
	parsed, err := url.Parse(rurl)
	if err != nil {
		log.Println("URL parse error")
		return "", err
	}
	return parsed.Host, nil
}

func LookupIPByHost(host string) (ips []net.IP, err error){
	ips, err = net.LookupIP(host)
	if err != nil {
		log.Printf("Could not lookup IP: %v", err)
		return ips, err
	}
	return ips, nil
}

func LookupIPByURL(rurl string) (ips []net.IP, err error) {
	host, err := LookupHostByURL(rurl)
	if err != nil {
		return ips, err
	}
	ips, err = LookupIPByHost(host)
	if err != nil {
		return ips, err
	}
	return ips, nil
}
