package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"
	"searchproxy/util/network"

	"github.com/didip/tollbooth"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/tb0hdan/memcache"

	"github.com/didip/tollbooth/limiter"

	log "github.com/sirupsen/logrus"
)

// Run - runs search proxy server iself
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

// setupRateLimitMiddleWare - configures client-based restrictions
func (sps *SearchProxyServer) setupRateLimitMiddleWare() (middleWare *limiter.Limiter) {
	middleWare = tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	middleWare.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})

	return
}

// RegisterMirrorsWithPrefix - create necessary configuration for multiple HTTP prefixes
func (sps *SearchProxyServer) RegisterMirrorsWithPrefix(mirrors []*mirrorsort.MirrorInfo, prefix, algorithm string) {
	requestTimeout := time.Duration(int64(sps.RequestTimeout) * time.Second.Nanoseconds())
	msConfig := &MirrorServerConfig{
		Cache:           memcache.New(log.New()),
		Mirrors:         mirrors,
		Prefix:          prefix,
		GeoIPDBFile:     sps.GeoIPDBFile,
		BuildInfo:       sps.BuildInfo,
		SearchAlgorithm: algorithm,
		RequestTimeout:  requestTimeout,
	}
	ms := NewMirrorServer(msConfig)
	middleWare := sps.setupRateLimitMiddleWare()
	sps.Gorilla.PathPrefix(prefix).Handler(tollbooth.LimitFuncHandler(middleWare, ms.CatchAllHandler))
	sps.Proxies = append(sps.Proxies, prefix)
}

// serveRoot - render index page (unexported)
func (sps *SearchProxyServer) serveRoot(w http.ResponseWriter, r *http.Request) {
	network.WriteNormalResponse(w, "Hello normal index\n")

	for _, proxy := range sps.Proxies {
		fmt.Fprintf(w, "Endpoint: %s\n", proxy)
	}
}

// ConfigFromFile - apply configuration read from file
func (sps *SearchProxyServer) ConfigFromFile(fpattern, fdir string) {
	const DefaultAlgorithm = "first"

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

	requestTimeout := time.Duration(int64(sps.RequestTimeout) * time.Second.Nanoseconds())
	sorter := mirrorsort.NewSorter(sps.GeoIPDBFile, sps.BuildInfo, requestTimeout)

	for _, cfg := range Config.Mirrors {
		log.Printf("[i] Registering mirror `%s` with prefix `%s`\n", cfg.Name, cfg.Prefix)
		mirrors := sorter.MirrorSort(cfg.URLs)

		if cfg.Algorithm == "" {
			cfg.Algorithm = DefaultAlgorithm
		}

		sps.RegisterMirrorsWithPrefix(mirrors, cfg.Prefix, cfg.Algorithm)
	}

	log.Println("[i] Mirror registration complete")
}

// SetDebug - enable/disable debug based on flag
func (sps *SearchProxyServer) SetDebug(debug bool) {
	sps.Debug = debug
	if sps.Debug {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
		// Register debug routes
		sps.Gorilla.HandleFunc("/debug/pprof/", pprof.Index)
		sps.Gorilla.HandleFunc("/debug/pprof/{profile}", pprof.Index)
		sps.Gorilla.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		sps.Gorilla.HandleFunc("/debug/pprof/profile", pprof.Profile)
		sps.Gorilla.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		sps.Gorilla.HandleFunc("/debug/pprof/trace", pprof.Trace)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetReportCaller(false)
		log.SetLevel(log.InfoLevel)
	}
}

// SetGeoIPDBFile - set path to GeoIP DB file
func (sps *SearchProxyServer) SetGeoIPDBFile(dbFile string) {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Fatalf("[E] Cannot start with non-existing GeoIP DB file: %s", dbFile)
	}

	sps.GeoIPDBFile = dbFile
}

// Stop - run shutdown chores
func (sps *SearchProxyServer) Stop() {
	// no code yet
}

// New - create search proxy server instance and populate it with data
func New(
	addr string,
	readTimeout,
	writeTimeout,
	requestTimeout int,
	buildInfo *miscellaneous.BuildInfo) (sps *SearchProxyServer) {
	sps = &SearchProxyServer{
		Addr:           addr,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		RequestTimeout: requestTimeout,
		BuildInfo:      buildInfo,
		Gorilla:        mux.NewRouter(),
	}

	sps.Gorilla.HandleFunc("/", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.htm", sps.serveRoot)
	sps.Gorilla.HandleFunc("/index.html", sps.serveRoot)

	return
}
