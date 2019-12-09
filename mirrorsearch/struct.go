package mirrorsearch

import (
	"searchproxy/memcache"
	"searchproxy/mirrorsort"
	"searchproxy/util/miscellaneous"
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
