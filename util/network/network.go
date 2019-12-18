package network

import (
	"net"

	log "github.com/sirupsen/logrus"

	"net/url"
)

// LookupHostByURL - return hostname by URL
func LookupHostByURL(rurl string) (host string, err error) {
	parsed, err := url.Parse(rurl)
	if err != nil {
		log.Println("URL parse error")
		return "", err
	}

	return parsed.Host, nil
}

// LookupIPByHost - return list of IP addresses for a host
func LookupIPByHost(host string) (ips []net.IP, err error) {
	ips, err = net.LookupIP(host)
	if err != nil {
		log.Printf("Could not lookup IP: %v", err)
		return ips, err
	}

	return ips, nil
}

// LookupIPByURL - return list of IP addresses for an URL
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

// IsLocalNetwork - confirm that IP is in local network
func IsLocalNetwork(ip net.IP) (result bool) {
	LocalNetworks := []string{
		"10.0.0.1/8",
		"127.0.0.1/8",
		"172.16.0.1/12",
		"192.168.0.0/16",
	}

	for _, network := range LocalNetworks {
		_, ipnet, err := net.ParseCIDR(network)
		if err != nil {
			// cannot and should not happen, but still
			log.Fatalf("LocalNetworks is broken!!!! %v", err)
		}

		if ipnet.Contains(ip) {
			result = true
			break
		}
	}
	// IP didn't match any of local network definitions above, could be IPv6 loopback, go with built-in method
	if !result {
		result = ip.IsLoopback()
	}

	return result
}

// IsLocalNetworkString - confirm that IP is in local network - works with string
func IsLocalNetworkString(ipAddress string) (result bool) {
	// convenience method
	return IsLocalNetwork(net.ParseIP(ipAddress))
}
