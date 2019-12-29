package server

import (
	"net/http"
	"time"

	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"

	"github.com/gorilla/mux"
	"github.com/tb0hdan/memcache"
)

// SearchProxyServer - proxy server with bound methods
type SearchProxyServer struct {
	Gorilla        *mux.Router
	Addr           string
	ReadTimeout    int
	WriteTimeout   int
	RequestTimeout int
	Proxies        []string
	GeoIPDBFile    string
	BuildInfo      *miscellaneous.BuildInfo
	Debug          bool
}

// MirrorServer - mirror server with bound methods
type MirrorServer struct {
	Prefix       string
	SearchMethod func(requestURI string, w http.ResponseWriter, r *http.Request)
}

// MirrorServerConfig - mirror server configuration
type MirrorServerConfig struct {
	Cache           *memcache.CacheType
	Mirrors         []*mirrorsort.MirrorInfo
	Prefix          string
	GeoIPDBFile     string
	BuildInfo       *miscellaneous.BuildInfo
	SearchAlgorithm string
	RequestTimeout  time.Duration
}

// MirrorConfig - individual mirror configuration
type MirrorConfig struct {
	Name      string   `mapstructure:"name"`
	Prefix    string   `mapstructure:"prefix"`
	Algorithm string   `mapstructure:"algorithm"`
	URLs      []string `mapstructure:"urls"`
}

// MirrorsConfig - configuration for all mirrors in a config section
type MirrorsConfig struct {
	Mirrors []MirrorConfig `mapstructure:"mirrors"`
}
