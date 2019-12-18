package mirrorsearch

import (
	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"

	"github.com/tb0hdan/memcache"
)

// MirrorSearch - mirror search with bound methods
type MirrorSearch struct {
	Cache       *memcache.CacheType
	Mirrors     []*mirrorsort.MirrorInfo
	Prefix      string
	GeoIPDBFile string
	BuildInfo   *miscellaneous.BuildInfo
}

// MirrorCache - mirror info cache
type MirrorCache struct {
	KnownURLs map[string]bool
}
