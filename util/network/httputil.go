package network

import (
	"fmt"
	"net"
	"net/http"
	"strings"
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

func StripRequestURI(requestURI, prefix string) (result string) {
	result = strings.TrimLeft(requestURI, prefix)
	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}

	return
}

func GetRemoteAddressFromRequest(r *http.Request) (addr string, err error) {
	var (
		remoteAddr string
	)

	remoteAddr, _, err = net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		// Something's very wrong with the request
		return "", err
	}

	addr = r.Header.Get("X-Real-IP")

	if len(addr) == 0 {
		addr = r.Header.Get("X-Forwarded-For")
	}

	// Could not get IP from headers
	if len(addr) == 0 {
		addr = remoteAddr
	} else if !IsLocalNetworkString(remoteAddr) { // IP is from headers, check whether we can trust it
		// Nope, use remote address instead
		addr = remoteAddr
	}

	return addr, nil
}
