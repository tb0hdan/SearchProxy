package server

import (
	"net/http"

	"searchproxy/mirrorsearch"
	"searchproxy/util/network"
)

// serveRoot - index page renderer (unexported)
func (ms *MirrorServer) serveRoot(w http.ResponseWriter, _ *http.Request) {
	network.WriteNormalResponse(w, "hello index")
}

// CatchAllHandler - catch all requested URLs and dispatch them accordingly
func (ms *MirrorServer) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	strippedURI := network.StripRequestURI(r.RequestURI, ms.Prefix)
	if strippedURI == "/" || strippedURI == "/index.htm" || strippedURI == "/index.html" {
		ms.serveRoot(w, r)
		return
	}
	// This is configured via NewMirrorServer
	ms.SearchMethod(r.RequestURI, w, r)
}

// NewMirrorServer - create mirror server instance and populate it properly
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
