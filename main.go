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
	Prefix string
}

func (ms *MirrorServer) StripRequestURI(requestURI string) (result string) {
	result = strings.TrimLeft(requestURI, ms.Prefix)
	if ! strings.HasPrefix(result, "/") {
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


type SearchProxyServer struct {
	Gorilla *mux.Router
	Addr string
	ReadTimeout int
	WriteTimeout int
	Proxies []string
}

func (sps *SearchProxyServer) Run() {
	srv := &http.Server{
		Handler: sps.Gorilla,
		Addr:    sps.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(int64(sps.WriteTimeout) * time.Second.Nanoseconds()),
		ReadTimeout:  time.Duration(int64(sps.ReadTimeout) * time.Second.Nanoseconds()),
	}

	log.Fatal(srv.ListenAndServe())
}

func (sps *SearchProxyServer) RegisterMirrorsWithPrefix(mirrors []string, prefix string) {
	cache := memcache.New()
	ms := &MirrorServer{Cache: cache, Mirrors: mirrors, Prefix: prefix}
	sps.Gorilla.PathPrefix(prefix).HandlerFunc(ms.CatchAllHandler)
	sps.Proxies = append(sps.Proxies, prefix)
}

func (sps *SearchProxyServer) serveRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello normal index\n")
	for _, proxy := range sps.Proxies {
		fmt.Fprintf(w, "Endpoint: %s\n", proxy)
	}
}

func NewSearchProxyServer(addr string, readTimeout, writeTimeout int) (sps *SearchProxyServer){
	sps = &SearchProxyServer{
		Addr: addr,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
	}
	sps.Gorilla = mux.NewRouter()
	sps.Gorilla.HandleFunc("/", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.htm", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.html", sps.serveRoot)
	return
}

func main() {
	searchProxyServer := NewSearchProxyServer("0.0.0.0:8000", 30, 30)
	searchProxyServer.RegisterMirrorsWithPrefix(MIRRORS, "/gentoo")
	searchProxyServer.Run()
}
