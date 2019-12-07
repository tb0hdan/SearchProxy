package mirrorsort

import (
	"searchproxy/geoip"
	"searchproxy/util/miscellaneous"
)

type MirrorStats struct {
	// Timestamp
	LastChecked           int64
	ConnectionsSinceStart int64
}

type MirrorInfo struct {
	URL         string
	IP          string
	PingMS      int64
	GeoIPInfo   *geoip.Info
	GeoIPDBFile string
	Stats       *MirrorStats
	BuildInfo   *miscellaneous.BuildInfo
	UUID        string
	// Used for closes mirror search, calculated only then
	Distance float64
}
