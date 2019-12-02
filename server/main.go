package server

import "C"
import (
	"fmt"
	"net/http"
	"time"

	"searchproxy/memcache"
	"searchproxy/mirrorsort"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type SearchProxyServer struct {
	Gorilla      *mux.Router
	Addr         string
	ReadTimeout  int
	WriteTimeout int
	Proxies      []string
	GeoIPDBFile  string
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

func (sps *SearchProxyServer) RegisterMirrorsWithPrefix(mirrors []*mirrorsort.MirrorInfo, prefix string) {
	cache := memcache.New()
	ms := &MirrorServer{Cache: cache, Mirrors: mirrors, Prefix: prefix, GeoIPDBFile: sps.GeoIPDBFile}
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

func (sps *SearchProxyServer) ConfigFromFile(fpattern, fdir string) {
	var Config MirrorsConfig

	viper.SetConfigName(fpattern)
	viper.AddConfigPath(fdir)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("err!")
		} else {
			// Config file was found but another error was produced
			log.Println("err no parse!")
		}
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("Unable to decode")
	}

	sorter := mirrorsort.NewSorter(sps.GeoIPDBFile)

	for _, cfg := range Config.Mirrors {
		log.Printf("Registering mirror `%s` with prefix `%s`\n", cfg.Name, cfg.Prefix)
		mirrors := sorter.MirrorSort(cfg.URLs)
		sps.RegisterMirrorsWithPrefix(mirrors, cfg.Prefix)
	}

	log.Println("SearchProxy started")
}

func (sps *SearchProxyServer) SetDebug(debug bool) {
	if debug {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
}

func (sps *SearchProxyServer) SetGeoIPDBFile(dbFile string) {
	sps.GeoIPDBFile = dbFile
}

func (sps *SearchProxyServer) Stop() {
	// no code yet
}

func New(addr string, readTimeout, writeTimeout int) (sps *SearchProxyServer) {
	sps = &SearchProxyServer{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Gorilla:      mux.NewRouter(),
	}

	sps.Gorilla.HandleFunc("/", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.htm", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.html", sps.serveRoot)

	return
}
