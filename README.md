# SearchProxy
This project offers functionality similar to HAProxy/Nginx though it checks for file
presence prior to returning redirect to respective upstream. Mainly intended for
opensource mirrors but can be used (possibly) as a CDN frontend.

Mirror selection algorithms so far:

- First available mirror (the're sorted by latency during app startup). YAML value: `first`
- The one closest to client (if none are matching, fallback to first available). YAML value: `closest`


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

`make`

`./searchproxy`

`wget http://localhost:8000/gentoo/distfiles/01-iosevka-1.14.1.zip`

or

`wget http://localhost:8000/debian/ls-lR.gz`

## Upstream proxy support
SearchProxy uses Go's built-in http client with HTTP proxy support. In order to send all
requests through a proxy, export environment variable like this:

`export HTTP_PROXY=http://example.com:8000`

## GeoIP notice
This project uses [GeoIP2-Golang](https://github.com/oschwald/geoip2-golang) which in turn
relies on [MaxMind's GeoLite database](https://dev.maxmind.com/geoip/geoip2/geolite2/)
