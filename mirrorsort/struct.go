package mirrorsort

import "searchproxy/geoip"

type MirrorInfo struct {
	URL string
	PingMS int64
	GeoIPInfo *geoip.GeoIPInfo
}

