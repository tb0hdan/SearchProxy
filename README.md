# SearchProxy

[![Build Status](https://travis-ci.org/tb0hdan/SearchProxy.svg?branch=master)](https://travis-ci.org/tb0hdan/SearchProxy)
[![GoDoc](https://godoc.org/github.com/tb0hdan/SearchProxy?status.svg)](https://godoc.org/github.com/tb0hdan/SearchProxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/tb0hdan/SearchProxy)](https://goreportcard.com/report/github.com/tb0hdan/SearchProxy)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftb0hdan%2FSearchProxy.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftb0hdan%2FSearchProxy?ref=badge_shield)

Redirect to backend server(s) that has the file(s) - no more four-o-four!

This project offers functionality similar to HAProxy/Nginx though it checks for file
presence before returning redirect to respective backend server. Started as a
frontend for opensource mirrors but can be used for other things like CDN.

Backend server selection algorithms so far:

- First available server (the're sorted by latency during app startup). YAML value: `first`
- The one closest to client (if none are good, fallback to first available). YAML value: `closest`
- Least connections (load balanced). YAML value: `leastconn`
- GeoBalance - combined approach. Three closest mirrors are rotated based on connections amount. YAML value: `geobalance`


Can be configured via YAML (with default being first available):

```yaml
mirrors:
  - name: "gentoo"
    prefix: "/gentoo"
    algorithm: "closest"
    urls:
      - "http://gentoo.mirrors.tera-byte.com/"
```

## Running

`make dockerimage`

`docker run -p 8000:8000 tb0hdan/searchproxy`


To confirm that SearchProxy returns link to file:

`wget --spider http://localhost:8000/gentoo/distfiles/01-iosevka-1.14.1.zip`


## HTTP proxy support
SearchProxy uses Go's built-in http client with HTTP proxy support. In order to send all
requests through a proxy, export environment variable like this:

`export HTTP_PROXY=http://example.com:8000`

## GeoIP notice
This project uses [GeoIP2-Golang](https://github.com/oschwald/geoip2-golang) which in turn
relies on [MaxMind's GeoLite database](https://dev.maxmind.com/geoip/geoip2/geolite2/)

Expected GeoIP DB file: `GeoLite2-City.mmdb`


## Thanks

[Docker Golang](https://www.docker.com/blog/docker-golang/)
