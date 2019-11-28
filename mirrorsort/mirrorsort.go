package mirrorsort

import (
	"log"
	"sort"

	"searchproxy/geoip"
	"searchproxy/httputil"
	"searchproxy/workerpool"
)

type MirrorInfo struct {
	URL string
	PingMS int64
	GeoIPInfo *geoip.GeoIPInfo
}
type ByPing []MirrorInfo

func (a ByPing) Len() int           { return len(a) }
func (a ByPing) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPing) Less(i, j int) bool { return a[i].PingMS < a[j].PingMS }

func PingHTTPWrapper(item interface{}) interface{} {
	url := item.(string)
	db := &geoip.GeoIPDB{}
	geoipinfo, err := db.LookupURL(url)
	if err != nil {
		log.Printf("PingHTTPWrapper %v\n", err)
	}
	return MirrorInfo{URL: url, PingMS: httputil.PingHTTP(url), GeoIPInfo: geoipinfo}
}

func MirrorSort(urls []string) (result []string){
	var (
		repackURL []interface{}
		repackMirror []interface{}
		mirrors []MirrorInfo
	)
	for _, url := range urls {
		repackURL = append(repackURL, url)
	}
	// FIXME: Obey system limits
	wp := workerpool.NewWorkerPool(128, PingHTTPWrapper)
	repackMirror = wp.ProcessItems(repackURL)

	for _, mirror := range repackMirror {
		mirrorInfo := mirror.(MirrorInfo)
		// Add only working mirrors
		if mirrorInfo.PingMS < 0 {
			continue
		}
		mirrors = append(mirrors, mirrorInfo)
	}

	sort.Sort(ByPing(mirrors))

	for _, mirror := range mirrors {
		result = append(result, mirror.URL)
	}
	return
}
