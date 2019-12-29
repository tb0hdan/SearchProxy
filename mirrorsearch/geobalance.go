package mirrorsearch

import (
	"net/http"
	"sort"
	"strings"

	"searchproxy/mirrorsort"
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

// GeoBalance - bound method that provides both fast mirror search and load balance.
func (ms *MirrorSearch) GeoBalance(requestURI string, w http.ResponseWriter, r *http.Request) { // nolint
	requestURI = network.StripRequestURI(requestURI, ms.Prefix)

	repackedMirrors := ms.GetMirrors(requestURI, r, true)
	// Select few closest ones
	if len(repackedMirrors) >= 3 {
		repackedMirrors = repackedMirrors[:3]
	}
	//
	sort.Sort(mirrorsort.ByConnection(repackedMirrors))

	if len(repackedMirrors) > 0 {
		mirror := repackedMirrors[0]
		url := strings.TrimRight(mirror.URL, "/") + requestURI
		log.Printf("Requested URL for %s found at %s with connection count %d",
			requestURI, url, int(mirror.Stats.ConnectionsSinceStart))
		ms.Redirect(mirror, url, w, r)

		return
	}

	// No applicable mirrors were found, fall back to FindMirrorFirst
	ms.FindMirrorFirst(requestURI, w, r)
}
