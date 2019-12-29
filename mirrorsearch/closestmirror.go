package mirrorsearch

import (
	"net/http"
	"sort"
	"strings"

	"searchproxy/mirrorsort"
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

// FindClosestMirror - this bound method searches for closest mirror based on client IP and returns redirect to it
func (ms *MirrorSearch) FindClosestMirror(requestURI string, w http.ResponseWriter, r *http.Request) { // nolint
	requestURI = network.StripRequestURI(requestURI, ms.Prefix)

	repackedMirrors := ms.GetMirrors(requestURI, r, true)
	sort.Sort(mirrorsort.ByDistance(repackedMirrors))

	if len(repackedMirrors) > 0 {
		mirror := repackedMirrors[0]
		url := strings.TrimRight(mirror.URL, "/") + requestURI
		log.Printf("Requested URL for %s found at %s with distance %d km", requestURI, url, int(mirror.Distance))
		ms.Redirect(mirror, url, w, r)

		return
	}

	// No applicable mirrors were found, fall back to FindMirrorFirst
	ms.FindMirrorFirst(requestURI, w, r)
}
