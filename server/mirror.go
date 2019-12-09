package server

import (
	"fmt"
	"net/http"
	"searchproxy/mirrorsearch"
	"searchproxy/util/network"
)

func (ms *MirrorServer) serveRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello index")
}

func (ms *MirrorServer) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	strippedURI := network.StripRequestURI(r.RequestURI, ms.Prefix)
	if strippedURI == "/" || strippedURI == "/index.htm" || strippedURI == "/index.html" {
		ms.serveRoot(w, r)
		return
	}
	// This is configured via NewMirrorServer
	ms.SearchMethod(r.RequestURI, w, r)
}

func NewMirrorServer(config *MirrorServerConfig) *MirrorServer {
	ms := &MirrorServer{
		Prefix: config.Prefix,
	}
	Search := &mirrorsearch.MirrorSearch{
		Cache:       config.Cache,
		Mirrors:     config.Mirrors,
		Prefix:      config.Prefix,
		GeoIPDBFile: config.GeoIPDBFile,
		BuildInfo:   config.BuildInfo,
	}
	// Should be set via *MirrorServerConfig / yml
	ms.SearchMethod = Search.SetMirrorSearchAlgorithm(config.SearchAlgorithm)

	return ms
}
