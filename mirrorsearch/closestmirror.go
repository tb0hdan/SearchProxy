package mirrorsearch

import (
	"net/http"
	"sort"
	"strings"

	"searchproxy/mirrorsort"
	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

func (ms *MirrorSearch) FindClosestMirror(requestURI string, w http.ResponseWriter, r *http.Request) {
	requestURI = network.StripRequestURI(requestURI, ms.Prefix)

	repackedMirrors := make([]*mirrorsort.MirrorInfo, 0)

	for _, mirror := range ms.Mirrors {
		distance := ms.GetDistanceRemoteMirror(r, mirror)
		// No distance or local IP
		if distance <= 0 {
			continue
		}
		// This is unacceptably slow atm, use cache later
		url := strings.TrimRight(mirror.URL, "/") + requestURI
		res, err := ms.CheckMirror(url)

		// Not found
		if err != nil {
			log.Println(err)
			continue
		} else {
			res.Body.Close()
		}

		mirror.Distance = distance
		repackedMirrors = append(repackedMirrors, mirror)
	}

	sort.Sort(mirrorsort.ByDistance(repackedMirrors))

	if len(repackedMirrors) > 0 {
		mirror := repackedMirrors[0]
		url := strings.TrimRight(mirror.URL, "/") + requestURI
		log.Printf("Requested URL for %s found at %s with distance %d", requestURI, url, int(mirror.Distance))
		ms.Redirect(mirror, url, w, r)

		return
	}

	// No applicable mirrors were found, fall back to FindMirrorFirst
	ms.FindMirrorFirst(requestURI, w, r)
}
