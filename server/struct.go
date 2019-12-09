package server

import (
	"net/http"

	"searchproxy/memcache"
	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"

	"github.com/gorilla/mux"
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

type MirrorServer struct {
	Prefix       string
	SearchMethod func(requestURI string, w http.ResponseWriter, r *http.Request)
}

type MirrorServerConfig struct {
	Cache           *memcache.CacheType
	Mirrors         []*mirrorsort.MirrorInfo
	Prefix          string
	GeoIPDBFile     string
	BuildInfo       *miscellaneous.BuildInfo
	SearchAlgorithm string
}

type MirrorConfig struct {
	Name      string   `mapstructure:"name"`
	Prefix    string   `mapstructure:"prefix"`
	Algorithm string   `mapstructure:"algorithm"`
	URLs      []string `mapstructure:"urls"`
}

type MirrorsConfig struct {
	Mirrors []MirrorConfig `mapstructure:"mirrors"`
}
