package mirrorsort

import (
	"searchproxy/util/miscellaneous"
	"searchproxy/util/network"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"searchproxy/geoip"
)

func (mi *MirrorInfo) UpdateMS() {
	myHTTP := network.NewHTTPUtilities(mi.BuildInfo)
	mi.PingMS = myHTTP.PingHTTP(mi.URL)
}

func (mi *MirrorInfo) UpdateIP() {
	ips, err := network.LookupIPByURL(mi.URL)
	if err != nil {
		log.Printf("Could not update IP: %v", err)
		return
	}

	mi.IP = ips[0].String()
}

func (mi *MirrorInfo) UpdateGeo() {
	db := geoip.New(mi.GeoIPDBFile)
	geoIPInfo, err := db.LookupURL(mi.URL)

	if err != nil {
		log.Printf("UpdateGeo %v\n", err)
	}

	mi.GeoIPInfo = geoIPInfo
}

func (mi *MirrorInfo) UpdateUUID() {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("could not generate UUID: %v", err)
	}

	mi.UUID = uuid4.String()
}

func (mi *MirrorInfo) Update() {
	mi.UpdateMS()
	mi.UpdateGeo()
	mi.UpdateIP()
	mi.UpdateUUID()
}

func (mi *MirrorInfo) PlusConnection() {
	mi.Stats.ConnectionsSinceStart++
}

func NewMirror(url, geoIPDBFile string, buildInfo *miscellaneous.BuildInfo) *MirrorInfo {
	return &MirrorInfo{
		URL:         url,
		Stats:       &MirrorStats{},
		GeoIPDBFile: geoIPDBFile,
		BuildInfo:   buildInfo,
	}
}
