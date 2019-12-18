package mirrorsort

import (
	"searchproxy/util/miscellaneous"
	"searchproxy/util/network"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"searchproxy/geoip"
)

// UpdateMS - update mirror ping in milliseconds
func (mi *MirrorInfo) UpdateMS() {
	myHTTP := network.NewHTTPUtilities(mi.BuildInfo)
	mi.PingMS = myHTTP.PingHTTP(mi.URL)
}

// UpdateIP - update mirror IP address
func (mi *MirrorInfo) UpdateIP() {
	ips, err := network.LookupIPByURL(mi.URL)
	if err != nil {
		log.Printf("Could not update IP: %v", err)
		return
	}

	mi.IP = ips[0].String()
}

// UpdateGeo - update mirror geographical info
func (mi *MirrorInfo) UpdateGeo() {
	db := geoip.New(mi.GeoIPDBFile)
	geoIPInfo, err := db.LookupURL(mi.URL)

	if err != nil {
		log.Printf("UpdateGeo %v\n", err)
	}

	mi.GeoIPInfo = geoIPInfo
}

// UpdateUUID - generate random UUID for each mirror
func (mi *MirrorInfo) UpdateUUID() {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("could not generate UUID: %v", err)
	}

	mi.UUID = uuid4.String()
}

// Update -  shorthand method that runs all updates
func (mi *MirrorInfo) Update() {
	mi.UpdateMS()
	mi.UpdateGeo()
	mi.UpdateIP()
	mi.UpdateUUID()
}

// PlusConnection - increase internal connection counter
func (mi *MirrorInfo) PlusConnection() {
	mi.Stats.ConnectionsSinceStart++
}

// NewMirror - instantiate mirror and populate structure
func NewMirror(url, geoIPDBFile string, buildInfo *miscellaneous.BuildInfo) *MirrorInfo {
	return &MirrorInfo{
		URL:         url,
		Stats:       &MirrorStats{},
		GeoIPDBFile: geoIPDBFile,
		BuildInfo:   buildInfo,
	}
}
