package router

import (
	"fmt"
	"net/http"
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
func match(routeP []string, reqP []string) (map[string]string, error) {

	// verify if incoming requests path and if registered route match.
	if len(routeP) != len(reqP) {
		return nil, fmt.Errorf("route path and request path do not match: %d, %d", len(routeP), len(reqP))
	}

	// preallocated map to store params of size len(route.Parts).
	// as number of params would vary on number of params in routeP, though it is safe to keep allocated size of len(routeP) for flexibility depending on number of parameters in a single route.
	params := make(map[string]string, len(routeP))

	for v := range routeP { // can use both reqP and routeP to loop over
		rp := routeP[v]
		reqp := reqP[v]

		if len(rp) > 0 && rp[0] == ':' { // validate whether route part contains an : (indicating that it is dynamic)
			key := rp[1:]
			params[key] = reqp // store dynamic route indicator as param to reqp value.
			continue           // continue to next loop iteration
		}

		// verifies whether request path matches registered route path
		// for RFC-3986 compliancy, it states that query, path components must be treated case-sensitive, so we dont normalise them here.
		if rp != reqp {
			return nil, fmt.Errorf("Request does not match Route path, found: %s", rp)
		}
	}

	//

	// everything matched
	return params, nil
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

// IsValidHTTPMethod determines whether a given method via router.Handle() or router.Group.Handle() is an valid http method.
func IsValidHTTPMethod(m string) bool {
	switch m {
	case http.MethodDelete,
		http.MethodGet,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace:
		return true

	default:
		return false
	}
}
