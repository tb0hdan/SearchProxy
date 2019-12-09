package mirrorsearch

import (
	"net/http"
	"strings"

	"searchproxy/util/network"

	log "github.com/sirupsen/logrus"
)

func (ms *MirrorSearch) FindMirrorFirst(requestURI string, w http.ResponseWriter, r *http.Request) {
	requestURI = network.StripRequestURI(requestURI, ms.Prefix)

	if value, ok := ms.Cache.Get(requestURI); ok {
		log.Printf("Cached URL for %s found at %s", requestURI, value)
		mirrorURL := value.(string)
		mirror := ms.FindMirrorByURL(mirrorURL)

		if mirror != nil {
			ms.Redirect(mirror, mirrorURL, w, r)
			return
		}

		log.Debugf("Could not find mirror for %s, proceeding with full search", requestURI)
	}

	for _, mirror := range ms.Mirrors {
		url := strings.TrimRight(mirror.URL, "/") + requestURI
		res, err := ms.CheckMirror(url)

		if err != nil {
			log.Println(err)
			continue
		} else {
			res.Body.Close()
		}

		if res.StatusCode == http.StatusOK {
			log.Printf("Requested URL for %s found at %s", requestURI, url)
			ms.Redirect(mirror, url, w, r)
			ms.Cache.SetEx(requestURI, url, 86400)

			return
		}
	}

	network.WriteNotFound(w)
}
