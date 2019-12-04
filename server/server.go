package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"searchproxy/memcache"
	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"

	"github.com/didip/tollbooth"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/didip/tollbooth/limiter"

	log "github.com/sirupsen/logrus"
)

type SearchProxyServer struct {
	Gorilla      *mux.Router
	Addr         string
	ReadTimeout  int
	WriteTimeout int
	Proxies      []string
	GeoIPDBFile  string
	BuildInfo    *miscellaneous.BuildInfo
}

func (sps *SearchProxyServer) Run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Handler: sps.Gorilla,
		Addr:    sps.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(int64(sps.WriteTimeout) * time.Second.Nanoseconds()),
		ReadTimeout:  time.Duration(int64(sps.ReadTimeout) * time.Second.Nanoseconds()),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[E] SearchProxy listen failed with: %v\n", err)
		}
	}()
	log.Print("[.] SearchProxy Started")

	<-done
	log.Print("[X] SearchProxy Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)

	defer func() {
		sps.Stop()
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Couldn't shut SearchProxy down:%+v", err)
	}

	log.Print("[!] SearchProxy normal exit")
}

func (sps *SearchProxyServer) setupRateLimitMiddleWare() (middleWare *limiter.Limiter) {
	middleWare = tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	middleWare.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})

	return
}

func (sps *SearchProxyServer) RegisterMirrorsWithPrefix(mirrors []*mirrorsort.MirrorInfo, prefix string) {
	cache := memcache.New()
	ms := &MirrorServer{
		Cache:       cache,
		Mirrors:     mirrors,
		Prefix:      prefix,
		GeoIPDBFile: sps.GeoIPDBFile,
		BuildInfo:   sps.BuildInfo,
	}
	middleWare := sps.setupRateLimitMiddleWare()
	sps.Gorilla.PathPrefix(prefix).Handler(tollbooth.LimitFuncHandler(middleWare, ms.CatchAllHandler))
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

	sorter := mirrorsort.NewSorter(sps.GeoIPDBFile, sps.BuildInfo)

	for _, cfg := range Config.Mirrors {
		log.Printf("[i] Registering mirror `%s` with prefix `%s`\n", cfg.Name, cfg.Prefix)
		mirrors := sorter.MirrorSort(cfg.URLs)
		sps.RegisterMirrorsWithPrefix(mirrors, cfg.Prefix)
	}

	log.Println("[i] Mirror registration complete")
}

func (sps *SearchProxyServer) SetDebug(debug bool) {
	if debug {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetReportCaller(false)
		log.SetLevel(log.InfoLevel)
	}
}

func (sps *SearchProxyServer) SetGeoIPDBFile(dbFile string) {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Fatalf("[E] Cannot start with non-existing GeoIP DB file: %s", dbFile)
	}

	sps.GeoIPDBFile = dbFile
}

func (sps *SearchProxyServer) Stop() {
	// no code yet
}

func New(addr string, readTimeout, writeTimeout int, buildInfo *miscellaneous.BuildInfo) (sps *SearchProxyServer) {
	sps = &SearchProxyServer{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		BuildInfo:    buildInfo,
		Gorilla:      mux.NewRouter(),
	}

	sps.Gorilla.HandleFunc("/", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.htm", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.html", sps.serveRoot)

	return
}
