package mirrorsort

import (
	"searchproxy/geoip"
	"searchproxy/util/miscellaneous"
)

// MirrorStats - individual mirror statistics
type MirrorStats struct {
	// Timestamp
	LastChecked           int64
	ConnectionsSinceStart int64
}

// MirrorInfo - individual mirror information
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

// ByPing - used for sorting by mirror ping latency
type ByPing []*MirrorInfo

// ByDistance - used for sorting by distance to mirror
type ByDistance []*MirrorInfo

// Sorter - sorter with bound methods
type Sorter struct {
	GeoIPDBFile string
	BuildInfo   *miscellaneous.BuildInfo
}
