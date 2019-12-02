package main

import (
	"flag"

	"searchproxy/server"
)

func main() {
	var debug = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	searchProxyServer := server.New("0.0.0.0:8000", 30, 30)
	searchProxyServer.SetDebug(*debug)
	searchProxyServer.ConfigFromFile("mirrors", ".")
	searchProxyServer.Run()
}
