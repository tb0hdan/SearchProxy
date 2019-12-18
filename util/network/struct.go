package network

import "searchproxy/util/miscellaneous"

// MirrorUnreachable - special value that signals limited mirror availability
const MirrorUnreachable = -1

// HTTPUtilities - HTTP utilities structure with bound methods
type HTTPUtilities struct {
	BuildInfo *miscellaneous.BuildInfo
}
