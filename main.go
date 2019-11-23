package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"searchproxy/memcache"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var MIRRORS = []string{"https://ftp.fau.de/gentoo", "https://ftp-stud.hs-esslingen.de/pub/Mirrors/gentoo",
	"http://ftp.fi.muni.cz/pub/linux/gentoo", "http://gentoo.mirror.web4u.cz",
	"http://gentoo.mirror.web4u.cz/", "http://gentoo.modulix.net/gentoo",
	"http://ftp-stud.hs-esslingen.de/pub/Mirrors/gentoo",
	"https://mirror.eu.oneandone.net/linux/distributions/gentoo/gentoo",
	"https://mirror.netcologne.de/gentoo/",
	"https://ftp.halifax.rwth-aachen.de/gentoo/",
	"http://ftp.ntua.gr/pub/linux/gentoo/",
	"https://mirrors.evowise.com/gentoo/",
	"https://ftp.snt.utwente.nl/pub/os/linux/gentoo/",
	"https://mirror.leaseweb.com/gentoo/",
	"http://ftp.vectranet.pl/gentoo/",
	"http://ftp.dei.uc.pt/pub/linux/gentoo/",
	"https://gentoo.wheel.sk/",
	"http://tux.rainside.sk/gentoo/",
	"https://mirror.bytemark.co.uk/gentoo/",
	"http://mirror.isoc.org.il/pub/gentoo/",
	"https://gentoo.ussg.indiana.edu/",
}

type MirrorServer struct {
	Cache *memcache.MemCacheType
	Mirrors []string
}

func (ms *MirrorServer) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" || r.RequestURI == "/index.htm" || r.RequestURI == "/index.html" {
		ms.serveRoot(w, r)
		return
	}

	ms.findMirror(r.RequestURI, w, r)
}

func (ms *MirrorServer) findMirror(requestURI string, w http.ResponseWriter, r *http.Request) {
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
			http.Redirect(w, r, mirrorURL+requestURI, http.StatusTemporaryRedirect)
			ms.Cache.SetEx(requestURI, mirrorURL+requestURI, 86400)
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



func main() {
	cache := memcache.New()
	ms := &MirrorServer{Cache: cache, Mirrors: MIRRORS}

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(ms.CatchAllHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
