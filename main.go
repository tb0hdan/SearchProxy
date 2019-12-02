package main

import (
	"flag"

	"searchproxy/server"
)

var debug = flag.Bool("debug", false, "enable debug")

func main() {
	flag.Parse()

	searchProxyServer := server.New("0.0.0.0:8000", 30, 30)
	searchProxyServer.SetDebug(*debug)
	searchProxyServer.ConfigFromFile("mirrors", ".")
	searchProxyServer.Run()
}
