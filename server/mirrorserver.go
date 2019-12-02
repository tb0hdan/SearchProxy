package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"searchproxy/geoip"
	"searchproxy/memcache"
	"searchproxy/mirrorsort"

	log "github.com/sirupsen/logrus"
)

type MirrorServer struct {
	Cache       *memcache.CacheType
	Mirrors     []*mirrorsort.MirrorInfo
	Prefix      string
	GeoIPDBFile string
}

func (ms *MirrorServer) StripRequestURI(requestURI string) (result string) {
	result = strings.TrimLeft(requestURI, ms.Prefix)
	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}

	return
}

func (ms *MirrorServer) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	strippedURI := ms.StripRequestURI(r.RequestURI)
	if strippedURI == "/" || strippedURI == "/index.htm" || strippedURI == "/index.html" {
		ms.serveRoot(w, r)
		return
	}

	ms.findMirror(r.RequestURI, w, r)
}

func (ms *MirrorServer) Redirect(mirror *mirrorsort.MirrorInfo, url string, w http.ResponseWriter, r *http.Request) {
	mirror.PlusConnection()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ms *MirrorServer) FindMirrorByURL(url string) (match *mirrorsort.MirrorInfo) {
	for _, mirror := range ms.Mirrors {
		if strings.HasPrefix(url, mirror.URL) {
			match = mirror
			log.Debugf("Mirror for URL %s found at %s", url, mirror.URL)

			break
		}
	}

	return
}

func (ms *MirrorServer) GetDistanceRemoteMirror(r *http.Request, mirror *mirrorsort.MirrorInfo) (distance float64) {
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

	// localhost
	if hostIP == "127.0.0.1" {
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

func (ms *MirrorServer) findMirror(requestURI string, w http.ResponseWriter, r *http.Request) {
	requestURI = ms.StripRequestURI(requestURI)

	if value, ok := ms.Cache.Get(requestURI); ok {
		log.Printf("Cached URL for %s found at %s", requestURI, value)
		mirror := ms.FindMirrorByURL(value)

		if mirror != nil {
			ms.Redirect(mirror, value, w, r)
			return
		}

		log.Debugf("Could not find mirror for %s, proceeding with full search", requestURI)
	}

	for _, mirror := range ms.Mirrors {
		log.Println(ms.GetDistanceRemoteMirror(r, mirror))
		url := strings.TrimRight(mirror.URL, "/") + requestURI

		res, err := http.Head(url)

		if err != nil {
			log.Println(err)
			continue
		} else {
			res.Body.Close()
		}

		if res.StatusCode == http.StatusOK {
			log.Printf("Requested URL for %s found at %s", requestURI, url)
			ms.Redirect(mirror, url, w, r)
			ms.Cache.SetEx(requestURI, url, 86400)

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 page not found")
}

func (ms *MirrorServer) serveRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello index")
}

type MirrorConfig struct {
	Name   string   `mapstructure:"name"`
	Prefix string   `mapstructure:"prefix"`
	URLs   []string `mapstructure:"urls"`
}

type MirrorsConfig struct {
	Mirrors []MirrorConfig `mapstructure:"mirrors"`
}
