package mirrorsort

import (
	"sort"

	"searchproxy/util/miscellaneous"
	"searchproxy/util/network"
	"searchproxy/util/system"
	"searchproxy/workerpool"
)

func (a ByPing) Len() int           { return len(a) }
func (a ByPing) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPing) Less(i, j int) bool { return a[i].PingMS < a[j].PingMS }

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

// PingHTTPWrapper - wrapper around mirror ping, used by worker pool
func (srt *Sorter) PingHTTPWrapper(item interface{}) interface{} {
	url := item.(string)
	mirror := NewMirror(url, srt.GeoIPDBFile, srt.BuildInfo)
	mirror.Update()

	return mirror
}

// MirrorSort - sort mirrors
func (srt *Sorter) MirrorSort(urls []string) (mirrors []*MirrorInfo) {
	var (
		repackMirror []interface{}
		repackURL    = make([]interface{}, 0, len(urls))
	)

	for _, url := range urls {
		repackURL = append(repackURL, url)
	}

	maxOpenFiles, _ := system.GetLimits()
	workerCount := int(maxOpenFiles / 8)

	if workerCount > 1024 {
		// Be sensible and don't overload system
		workerCount = 1024
	}

	wp := workerpool.New(workerCount, srt.PingHTTPWrapper)
	repackMirror = wp.ProcessItems(repackURL)

	for _, mirror := range repackMirror {
		mirrorInfo := mirror.(*MirrorInfo)
		// Add only working mirrors
		if mirrorInfo.PingMS <= network.MirrorUnreachable {
			continue
		}

		mirrors = append(mirrors, mirrorInfo)
	}

	sort.Sort(ByPing(mirrors))

	return mirrors
}

// NewSorter - create mirror sorter instance and populate it properly
func NewSorter(geoIPDBFile string, buildInfo *miscellaneous.BuildInfo) *Sorter {
	return &Sorter{
		GeoIPDBFile: geoIPDBFile,
		BuildInfo:   buildInfo,
	}
}
