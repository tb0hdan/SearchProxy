package mirrorsearch

import (
	"net"
	"net/http"
	"strings"

	"searchproxy/geoip"
	"searchproxy/mirrorsort"
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

func (ms *MirrorSearch) Redirect(mirror *mirrorsort.MirrorInfo, url string, w http.ResponseWriter, r *http.Request) {
	mirror.PlusConnection()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

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

func (ms *MirrorSearch) GetDistanceRemoteMirror(r *http.Request, mirror *mirrorsort.MirrorInfo) (distance float64) {
	var (
		err error
	)

	hostIP := r.Header.Get("X-Real-IP")

	if hostIP == "" {
		hostIP = r.Header.Get("X-Forwarded-For")
	}

	if hostIP == "" {
		hostIP, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// Something's very wrong with the request
			return -1
		}
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
	}

	log.Println(distance, hostIP, mirror.IP)

	return distance
}

func (ms *MirrorSearch) CheckMirror(mirrorURL string) (res *http.Response, err error) {
	// This method will be extended with rate limiting a little bit later
	myHTTP := network.NewHTTPUtilities(ms.BuildInfo)
	return myHTTP.HTTPHEAD(mirrorURL)
}

func (ms *MirrorSearch) SetMirrorSearchAlgorithm(algorithm string) (result func(requestURI string,
	w http.ResponseWriter, r *http.Request)) {
	switch algorithm {
	case "first":
		result = ms.FindMirrorFirst
	default:
		log.Fatalf("Unknown mirror search algorithm: %s\n", algorithm)
	}

	return result
}
