# SearchProxy
Search for files on mirrors and redirect to matching one


## Running

`go run *.go`

`wget http://localhost:8000/gentoo/distfiles/01-iosevka-1.14.1.zip`

or

`wget http://localhost:8000/debian/ls-lR.gz`


## GeoIP notice
This project uses ![GeoIP2-Golang](https://github.com/oschwald/geoip2-golang) which in turn
relies on ![MaxMind's GeoLite database](https://dev.maxmind.com/geoip/geoip2/geolite2/)
