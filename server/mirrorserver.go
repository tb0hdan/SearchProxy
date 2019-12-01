package server

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"searchproxy/memcache"
)

type MirrorServer struct {
	Cache   *memcache.MemCacheType
	Mirrors []string
	Prefix  string
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

func (ms *MirrorServer) findMirror(requestURI string, w http.ResponseWriter, r *http.Request) {
	requestURI = ms.StripRequestURI(requestURI)

	for _, mirrorURL := range ms.Mirrors {
		url := strings.TrimRight(mirrorURL, "/") + requestURI
		if value, ok := ms.Cache.Get(requestURI); ok {
			log.Printf("Cached URL for %s found at %s", requestURI, url)
			http.Redirect(w, r, value, http.StatusTemporaryRedirect)
			return
		}
		res, err := http.Head(url)
		//defer res.Body.Close()

		if err != nil {
			log.Println(err)
			continue
		}
		if res.StatusCode == http.StatusOK {
			log.Printf("Requested URL for %s found at %s", requestURI, url)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			ms.Cache.SetEx(requestURI, url, 86400)
			return
		}

	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 page not found")
}

func (ms *MirrorServer) serveRoot(w http.ResponseWriter, r *http.Request) {
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

