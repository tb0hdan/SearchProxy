package mirrorsearch

import (
	"net/http"
	"strings"

	"searchproxy/mirrorsort"

	log "github.com/sirupsen/logrus"
)

// GetMirrors - returns list of mirrors suitable for sending requests to
func (ms *MirrorSearch) GetMirrors(requestURI string, r *http.Request,
	checkDistance bool) (mirrors []*mirrorsort.MirrorInfo) {
	for _, mirror := range ms.Mirrors {
		distance := ms.GetDistanceRemoteMirror(r, mirror)
		// No distance or local IP
		if distance <= 0 && checkDistance {
			continue
		}
		// This is unacceptably slow atm, use cache later
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

		mirror.Distance = distance
		mirrors = append(mirrors, mirror)
	}

	return mirrors
}
