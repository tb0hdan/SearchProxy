package main

import (
	"flag"
	"path"
	"strings"

	"searchproxy/server"
)

func main() {
	var (
		debug        = flag.Bool("debug", false, "enable debug")
		geoIPDBFile  = flag.String("geoipdb", "GeoLite2-City.mmdb", "Path to Maxmind's DB file")
		bind         = flag.String("bind", "0.0.0.0:8000", "Address to bind to, host:port")
		readTimeout  = flag.Int("readt", 30, "Read timeout, seconds")
		writeTimeout = flag.Int("writet", 30, "Write timeout, seconds")
		mirrorsPath  = flag.String("mirrors", "./mirrors.yml", "Path to mirrors.yml file")
	)

	flag.Parse()

	searchProxyServer := server.New(*bind, *readTimeout, *writeTimeout)
	searchProxyServer.SetDebug(*debug)
	searchProxyServer.SetGeoIPDBFile(*geoIPDBFile)

	directory, file := path.Split(*mirrorsPath)
	filePattern := strings.TrimRight(file, ".yml")

	searchProxyServer.ConfigFromFile(filePattern, directory)
	searchProxyServer.Run()
}
