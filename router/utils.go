package router

import (
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



// path joining helper (for Grouping)
func Join(p1, p2 string) string {
	switch {
	case p1 == "/":
		return "/" + strings.TrimLeft(p2, "/")
	case p2 == "/":
		return strings.TrimRight(p1, "/") + "/"
	default:
		return strings.TrimRight(p1, "/") + "/" + strings.TrimLeft(p2, "/")
	}
}
