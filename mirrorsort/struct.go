package mirrorsort

import "searchproxy/geoip"

type MirrorStats struct {
	// Timestamp
	LastChecked           int64
	ConnectionsSinceStart int64
}

type MirrorInfo struct {
	URL       string
	IP        string
	PingMS    int64
	GeoIPInfo *geoip.Info
	Stats     *MirrorStats
}
