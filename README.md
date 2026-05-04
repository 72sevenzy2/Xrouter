<h1 align="center">Xrouter</h1>
An http router built from scratch, built ontop of the go stdlib (net/http).

# key features:

- No external dependencies used, other than a tool (json-parser) i made to handle json responses which is used here.
- Includes middleware(s) > basicAuth (with username and password), BearerAuth (with a bearer token), logging middleware, a timeout middleware (to cut off slow requests), and a recoverer middleware to prevent server crashes.
- Route specific middleware(s) > supports middlewares for specific routes without applying globally to all the routes.
- Supports route parameters > routes can have dynamic parameters for example, "/:id", although this feature can be slower than http routers like chi/gin as the implementation involves looping over routes to check if its dynamic (O(n)). But for non-dynamic routes (without parameters), this router uses a simple map to store all routes without looping over them which is faster (O(1)).

Its also important to note that when handling dynamic routes, the route parameters are stored in context in a map, so you can retrieve a particular field in that map via the Param() func. An example will be shown below:

```
 package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Handle(http.MethodGet, "/user/:64", func(w http.ResponseWriter, r *http.Request) {
		HandleThisFunc()...
		param, ok := router.Param(w, "user") // returns "64"
		if !ok {
			handleErr()
		}
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```

# General usage:

```
 package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Handle(http.MethodGet, "/resp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("responded"))
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
> That of course is a example with no middlewares attached yet.

Example usage with all the middlewares:

```

package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Use(router.Recoverer()) // recoverer middleware always goes first, prevents server crashes when a bug has occured.
	r.Use(router.Logger()) // standard logging middleware (to view request details.)
	r.Use(router.BearerAuth("secretKey")) // can be any token (which has to be a string),
	r.Use(router.BasicAuth("user1", "password1234")) // parameters username and password need to be included when using.
	r.Use(router.Timeout(5)) // can be any time (its in seconds) depending on how long you want the time limit on every request.
	// the Timeout mw is used for prevent slow requests by setting a timeout in which the request should last.

	r.Handle(http.MethodGet, "/resp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("responded"))
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
^ it is also important to note that you can limit how much the logging middleware reads from the request body, the default is 1 kilobyte as of now, but you can change it via passing in a SetBody() func as the parameter in the Logger() func, an example would be: "r.Use(router.Logger(SetBody(1024 * 2)))" (limit to 2 kilobytes.) but make sure the size your going to configure is of appropriate type (int64).

Example usage with route-specific middleware:

```
package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Handle(http.MethodGet, "/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}, router.Recoverer(), router.Logger()) // you can do route-specific middleware(s) like this (can be 1 or many).

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
