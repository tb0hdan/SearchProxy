package mirrorsort

import (
	"sort"

	"searchproxy/workerpool"
)

type ByPing []*MirrorInfo

func (a ByPing) Len() int           { return len(a) }
func (a ByPing) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPing) Less(i, j int) bool { return a[i].PingMS < a[j].PingMS }

func PingHTTPWrapper(item interface{}) interface{} {
	url := item.(string)
	mirror := &MirrorInfo{URL: url}
	mirror.Update()
	return mirror
}

func MirrorSort(urls []string) (result []*MirrorInfo){
	var (
		repackURL []interface{}
		repackMirror []interface{}
		mirrors []*MirrorInfo
	)
	for _, url := range urls {
		repackURL = append(repackURL, url)
	}
	// FIXME: Obey system limits
	wp := workerpool.New(128, PingHTTPWrapper)
	repackMirror = wp.ProcessItems(repackURL)

	for _, mirror := range repackMirror {
		mirrorInfo := mirror.(*MirrorInfo)
		// Add only working mirrors
		if mirrorInfo.PingMS < 0 {
			continue
		}
		mirrors = append(mirrors, mirrorInfo)
	}

	sort.Sort(ByPing(mirrors))

	for _, mirror := range mirrors {
		result = append(result, mirror)
	}
	return
}
