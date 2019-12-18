package mirrorsearch

import (
	"net/http"
	"strings"

	"searchproxy/geoip"
	"searchproxy/mirrorsort"
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

// Redirect - bound method that redirects to provided URL
func (ms *MirrorSearch) Redirect(mirror *mirrorsort.MirrorInfo, url string, w http.ResponseWriter, r *http.Request) {
	mirror.PlusConnection()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// FindMirrorByURL - bound method that returns mirror instance (if found) by mirror URL
func (ms *MirrorSearch) FindMirrorByURL(url string) (match *mirrorsort.MirrorInfo) {
	for _, mirror := range ms.Mirrors {
		if strings.HasPrefix(url, mirror.URL) {
			match = mirror
			log.Debugf("Mirror for URL %s found at %s", url, mirror.URL)

			break
		}
	}

	return
}

// GetDistanceRemoteMirror - get distance to remote mirror, used to find closest one
func (ms *MirrorSearch) GetDistanceRemoteMirror(r *http.Request, mirror *mirrorsort.MirrorInfo) (distance float64) {
	var (
		err error
	)

	hostIP, err := network.GetRemoteAddressFromRequest(r)

	if err != nil {
		return -1
	}

	if network.IsLocalNetworkString(hostIP) {
		return 0
	}

	geo := geoip.New(ms.GeoIPDBFile)

	if mirror.GeoIPInfo == nil {
		return -1
	}

	distance, err = geo.DistanceIPLatLon(hostIP, mirror.GeoIPInfo.Latitude, mirror.GeoIPInfo.Longitude)

	if err != nil {
		log.Printf("Distance err: %v", err)
		return -1
	}

	return distance
}

// CheckMirror - ping mirror and confirm that it's ok
func (ms *MirrorSearch) CheckMirror(mirrorURL string) (res *http.Response, err error) {
	// This method will be extended with rate limiting a little bit later
	myHTTP := network.NewHTTPUtilities(ms.BuildInfo)
	return myHTTP.HTTPHEAD(mirrorURL)
}

// SetMirrorSearchAlgorithm - used for configuration purposes. Sets mirror search algorithm during startup.
func (ms *MirrorSearch) SetMirrorSearchAlgorithm(algorithm string) (result func(requestURI string,
	w http.ResponseWriter, r *http.Request)) {
	switch algorithm {
	case "first":
		result = ms.FindMirrorFirst
	case "closest":
		result = ms.FindClosestMirror
	case "geobalance":
		result = ms.GeoBalance
	default:
		log.Fatalf("Unknown mirror search algorithm: %s\n", algorithm)
	}

	return result
}
