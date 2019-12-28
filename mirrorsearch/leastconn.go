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
	repackedMirrors := make([]*mirrorsort.MirrorInfo, 0)

	for _, mirror := range ms.Mirrors {
		url := strings.TrimRight(mirror.URL, "/") + requestURI

		if value, ok := ms.Cache.Get(mirror.UUID); !ok {
			res, err := ms.CheckMirror(url)
			// Not found
			if err != nil {
				log.Println(err)
				continue
			} else {
				res.Body.Close()
			}

			mc := &MirrorCache{KnownURLs: map[string]bool{
				url: true,
			}}
			ms.Cache.SetEx(mirror.UUID, mc, 86400)
		} else {
			mirrorCache := value.(*MirrorCache)
			if _, ok := mirrorCache.KnownURLs[url]; ok {
				// URL is known
				log.Debugf("Found matching URL in cache: %s", url)
			} else {
				// URL is unknown
				res, err := ms.CheckMirror(url)
				if err != nil {
					log.Println(err)
					continue
				} else {
					res.Body.Close()
					if res.StatusCode == http.StatusOK {
						mirrorCache.KnownURLs[url] = true
					}
				}
			}
		}

		repackedMirrors = append(repackedMirrors, mirror)
	}

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
