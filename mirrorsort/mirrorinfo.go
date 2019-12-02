package mirrorsort

import (
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"

	"searchproxy/geoip"
)

func (mi *MirrorInfo) UpdateMS() {
	mi.PingMS = network.PingHTTP(mi.URL)
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
	db := geoip.New("GeoLite2-City.mmdb")
	geoipinfo, err := db.LookupURL(mi.URL)

	if err != nil {
		log.Printf("UpdateGeo %v\n", err)
	}

	mi.GeoIPInfo = geoipinfo
}

func (mi *MirrorInfo) Update() {
	mi.UpdateMS()
	mi.UpdateGeo()
	mi.UpdateIP()
}

func (mi *MirrorInfo) PlusConnection() {
	if mi.Stats == nil {
		mi.Stats = &MirrorStats{}
	}

	mi.Stats.ConnectionsSinceStart++
}
