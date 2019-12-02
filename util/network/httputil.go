package network

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
		log.Debugf("An error %v occurred while running ping on %s", err, url)
		// failed servers should be marked as slow, with negative values
		elapsed = MirrorUnreachable * elapsed
	} else {
		res.Body.Close()
	}

	return
}
