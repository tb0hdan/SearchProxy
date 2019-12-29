package mirrorsearch

import (
	"net/http"
	"strings"

	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

// FindMirrorFirst - this bound method looks for first mirror in the list that has requestURI.
// Mirrors are sorted ascending by ping
func (ms *MirrorSearch) FindMirrorFirst(requestURI string, w http.ResponseWriter, r *http.Request) {
	requestURI = network.StripRequestURI(requestURI, ms.Prefix)

	repackedMirrors := ms.GetMirrors(requestURI, r, false)
	if len(repackedMirrors) > 0 {
		mirror := repackedMirrors[0]
		url := strings.TrimRight(mirror.URL, "/") + requestURI
		log.Printf("Requested URL for %s found at %s", requestURI, url)
		ms.Redirect(mirror, url, w, r)

		return
	}

	network.WriteNotFound(w)
}
