package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"searchproxy/server"
)

var debug = flag.Bool("debug", false, "enable debug")

func main() {
	flag.Parse()
	if *debug {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}

	searchProxyServer := server.New("0.0.0.0:8000", 30, 30)
	searchProxyServer.ConfigFromFile("mirrors", ".")
	searchProxyServer.Run()
}
