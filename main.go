package main

import (
	"flag"

	"searchproxy/server"
)

func main() {
	var (
		debug       = flag.Bool("debug", false, "enable debug")
		geoIPDBFile = flag.String("geoipdb", "GeoLite2-City.mmdb", "Path to Maxmind's DB file")
	)

	flag.Parse()

	searchProxyServer := server.New("0.0.0.0:8000", 30, 30)
	searchProxyServer.SetDebug(*debug)
	searchProxyServer.SetGeoIPDBFile(*geoIPDBFile)
	searchProxyServer.ConfigFromFile("mirrors", ".")
	searchProxyServer.Run()
}
