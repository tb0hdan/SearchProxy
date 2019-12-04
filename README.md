# SearchProxy
This project offers functionality similar to HAProxy/Nginx though it checks for file
presence prior to returning redirect to respective upstream. Mainly intended for
opensource mirrors but can be used (possibly) as a CDN frontend.

Returns fastest (lowest HTTP ping) upstream so far.

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
