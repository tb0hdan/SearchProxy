package network

import (
	"fmt"
	"net/http"
	"time"

	"searchproxy/util/miscellaneous"

	log "github.com/sirupsen/logrus"
)

type HTTPUtilities struct {
	BuildInfo *miscellaneous.BuildInfo
}

func (hu *HTTPUtilities) HTTPHEAD(url string) (res *http.Response, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent",
		fmt.Sprintf("Mozilla/5.0 (compatible; SearchProxy/%s; %s; +https://github.com/tb0hdan/SearchProxy)",
			hu.BuildInfo.Version, hu.BuildInfo.GoVersion))

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (hu *HTTPUtilities) PingHTTP(url string) (elapsed int64) {
	start := time.Now().UnixNano()
	res, err := hu.HTTPHEAD(url)
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

func NewHTTPUtilities(buildInfo *miscellaneous.BuildInfo) *HTTPUtilities {
	return &HTTPUtilities{BuildInfo: buildInfo}
}
