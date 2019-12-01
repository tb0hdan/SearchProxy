package server

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

func (sps *SearchProxyServer) ConfigFromFile(fpattern, fdir string) {
	viper.SetConfigName("mirrors")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("err!")
			// Config file not found; ignore error if desired
		} else {
			log.Println("err no parse!")
			// Config file was found but another error was produced
		}
	}
	var C MirrorsConfig

	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Unable to decode")
	}
	for _, cfg := range C.Mirrors {
		log.Printf("Registering mirror `%s` with prefix `%s`\n", cfg.Name, cfg.Prefix)
		sortedURLs := mirrorsort.MirrorSort(cfg.URLs)
		sps.RegisterMirrorsWithPrefix(sortedURLs, cfg.Prefix)
	}
	log.Println("SearchProxy started")
}

func NewSearchProxyServer(addr string, readTimeout, writeTimeout int) (sps *SearchProxyServer) {
	sps = &SearchProxyServer{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	sps.Gorilla = mux.NewRouter()
	sps.Gorilla.HandleFunc("/", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.htm", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.html", sps.serveRoot)
	return
}
