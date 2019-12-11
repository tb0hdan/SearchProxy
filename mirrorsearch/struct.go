package mirrorsearch

import (
	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"

	"github.com/tb0hdan/memcache"
)

type MirrorSearch struct {
	Cache       *memcache.CacheType
	Mirrors     []*mirrorsort.MirrorInfo
	Prefix      string
	GeoIPDBFile string
	BuildInfo   *miscellaneous.BuildInfo
}

type MirrorCache struct {
	KnownURLs map[string]bool
}
