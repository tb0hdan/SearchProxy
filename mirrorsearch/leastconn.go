package mirrorsearch

import (
	"net/http"
	"sort"
	"strings"

	"searchproxy/mirrorsort"
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

// LeastConn - this bound method searches for mirror based on least amount of connections returns redirect to it
func (ms *MirrorSearch) LeastConn(requestURI string, w http.ResponseWriter, r *http.Request) { // nolint
	requestURI = network.StripRequestURI(requestURI, ms.Prefix)

	repackedMirrors := ms.GetMirrors(requestURI, r, false)
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
