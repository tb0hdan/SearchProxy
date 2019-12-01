package mirrorsort

import (
	log "github.com/sirupsen/logrus"
	"searchproxy/geoip"
	"searchproxy/httputil"
)

func (mi *MirrorInfo) UpdateMS() {
	mi.PingMS = httputil.PingHTTP(mi.URL)
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
}

