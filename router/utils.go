package router

import (
	"bytes"
	"strings"
)

// func to check whether route is dynamic or static
func isDynamic(path string) bool {
	return strings.Contains(path, ":") // strings.Contains() returns a boolean
}

// func to split route path
func splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

// func to match request parts and route parts
func match(routeP []string, reqP []string) (bool, map[string]string) {
	// check if lengths match (if it isnt then it cannot be compared)
	if len(routeP) != len(reqP) {
		return false, nil
	}

	params := make(map[string]string, 500) // preallocate number of params

	for v := range routeP { // can use both reqP and routeP to loop over
		rp := routeP[v]
		reqp := reqP[v]

		if len(rp) > 0 && rp[0] == ':' { // if includes : then its a param, and checking whether its greater than 0 prevents crashes for when if it is an empty string.
			key := rp[1:]      // removes the :
			params[key] = reqp //
			continue           // continue to next loop iteration
		}

		// otherwise, if not a param then reject
		if rp != reqp {
			return false, nil
		}
	}

	//

	// everything matched
	return true, params
}

// custom limited writer function for Logger() to limit body size reading.
type LimitedBuffer struct {
	buf   *bytes.Buffer
	limit uint32
}

// custom write func for LimitedBuffer, (allocates new slice based on truncated size on original slice)
func (l *LimitedBuffer) Write(p []byte) (int, error) {
	remaining := l.limit - uint32(l.buf.Len())

	if remaining <= 0 {
		return len(p), nil
	}

	if len(p) > int(remaining) { // check if []byte that is being written to l.buf exceeds remaining.
		p = p[:remaining] // truncate
	}

	return l.buf.Write(p)
}
