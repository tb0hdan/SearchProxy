package httputil

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func PingHTTP(url string) (elapsed int64) {
	start := time.Now().UnixNano()
	res, err := http.Head(url)
	elapsed = (time.Now().UnixNano() - start) / time.Millisecond.Nanoseconds()
	if err != nil {
		log.Debug("An error %v occured while running ping on %s", err, url)
		// failed servers should be marked as slow, with negative values
		elapsed = -1 * elapsed
	} else {
		res.Body.Close()
	}
	return
}
