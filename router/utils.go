package router

import (
	"net/http"
	"strings"

	"github.com/72sevenzy2/http-router/core"
)

// func to check whether route is dynamic or static
func isDynamic(path string) bool {
	return strings.Contains(path, ":") // strings.Contains() returns a boolean
}

// func to split route path
func splitPath(path string) []string {
	// manual string indexing to avoid strings.genSplit overhead.
	path = strings.Trim(path, "/") // normalise from any other outter slashes.
	if path == "" {
		return nil
	}

	n := 1
	for v := 0; v < len(path); v++ {
		if path[v] == '/' {
			n++
		}
	}

	// preallocate according to n (number of parts in the request)
	parts := make([]string, 0, n)
	start := 0
	for v := 0; v < len(path); v++{
		if path[v] == '/' {
			parts = append(parts, path[start:v])
			start = v + 1 // update start position after each append
		}
	}

	parts = append(parts, path[start:]) // ensures the last part URL (after final occurence of '/') is appended.
	return parts
}

// func to match request parts and route parts
func match(routeP []string, reqP []string) (core.Params, error) {

	// verify if incoming requests path and if registered route match.
	if len(routeP) != len(reqP) {
		return nil, InvalidPathError
	}

	// preallocated map to store params of size len(route.Parts).
	// as number of params would vary on number of params in routeP, though it is safe to keep allocated size of len(routeP) for flexibility depending on number of parameters in a single route.
	//params := make(map[string]string, len(routeP))

	var params core.Params

	for v := range routeP { // can use both reqP and routeP to loop over
		rp := routeP[v]
		reqp := reqP[v]

		if len(rp) > 0 && rp[0] == ':' { // validate whether route part contains an : (indicating that it is dynamic)
			//params[key] = reqp // store dynamic route indicator as param to reqp value.
			params = append(params, &core.Param{Key: rp[1:], Value: reqp})
			continue // continue to next loop iteration
		}

		// verifies whether request path matches registered route path
		// for RFC-3986 compliancy, it states that query, path components must be treated case-sensitive, so we dont normalise them here.
		if rp != reqp {
			return nil, InvalidPathError
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
