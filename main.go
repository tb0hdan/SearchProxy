package main

import (
	"flag"
	"path"
	"strings"

	"searchproxy/server"
	"searchproxy/util/miscellaneous"
)

// Global vars for versioning
var (
	Build     = "unset" // nolint
	BuildDate = "unset" // nolint
	GoVersion = "unset" // nolint
	Version   = "unset" // nolint
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

	buildInfo := &miscellaneous.BuildInfo{
		Build:     Build,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		Version:   Version,
	}
	searchProxyServer := server.New(*bind, *readTimeout, *writeTimeout, buildInfo)
	searchProxyServer.SetDebug(*debug)
	searchProxyServer.SetGeoIPDBFile(*geoIPDBFile)

	directory, file := path.Split(*mirrorsPath)
	filePattern := strings.TrimRight(file, ".yml")

	searchProxyServer.ConfigFromFile(filePattern, directory)
	searchProxyServer.Run()
}
