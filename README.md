<h1 align="center">Xrouter.</h1>
<h2 align="center">An http router built from scratch, ontop of the go stdlib (net/http).</h2>

# key features:

- Method handlers for every http method.
- Route grouping for nested routes.
- Built in rate limiting middleware.
- No external dependencies (excluding ones made for Xrouter).
- Supports middleware aswell as route specific middleware.
- Supports dynamic routes.

Its also important to note that when handling dynamic routes, the route parameters are stored in a custom Request struct map field (params) in a map, to be retrieved from a particular field in that map via r.params, example:

```
 package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/Xrouter"
)

func main() {
	r := router.NewRouter()

	r.Get("/user/:64", func(w http.ResponseWriter, r *router.Request) {
		HandleThisFunc()...
		param, ok := r.params["user"]
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
<br>
<h4 align="center">note: when writing handlers, do make sure your using the custom Request struct instead of http.Request, which is router.Request</h4>

<h1 align="center">General usage:</h1>

```
 package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/Xrouter"
)

func main() {
	r := router.NewRouter()

	r.Get("/resp", func(w http.ResponseWriter, r *router.Request) {
		w.Write([]byte("responded"))
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
<h2 align="center">And for route grouping:</h2>

```
package main

import (
	"net/http"

	"github.com/72sevenzy2/Xrouter"
)

func main() {
	r := router.NewRouter()

	api := r.Group("/parent")
	/* 
	for inline nesting (without having to create seperate child routes individually):
	
	api := r.Group("/parent", "child1", "child2")
	(in which the registered route would be /parent/child1/child2 )

	though this prevents configuring middleware for a specific route
	*/

	v1 := api.Group("/child1")
	v1.Use(router.Logger(0)) // you can use middleware with specific route inside nested routes.

	v2 := v1.Group("/child2")
	v2.Use(router.BasicAuth("username", "password")) // would indicate only this route having this middleware

	v2.Get("/test", func(w http.ResponseWriter, r *router.Request) {
		w.Write([]byte("pong"))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
```

<br>
<h2 align="center" style.Background="gray">Usage with all the middleware(s):</h2>

```

package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/Xrouter"
)

func main() {
	r := router.NewRouter()

	limiter := router.NewLimiter(100, 1) // NewLimiter() takes in a maximum number of requests, and how often each token is refilled for      // the client.

	r.Use(limiter.RateLimiter())
	
	r.Use(router.Recoverer()) // recoverer middleware always goes first, prevents server crashes when a bug has occured.
	r.Use(router.Logger(0)) // standard logging, in which param 0 indicates default body size logging.

	r.Use(router.BearerAuth("secretKey")) // can be any token (which has to be a string),
	r.Use(router.BasicAuth("user1", "password1234")) // parameters username and password need to be included when using.
	r.Use(router.Timeout(5)) // can be any time (its in seconds) depending on how long you want the time limit on every request.
	// the Timeout mw is used for prevent slow requests by setting a timeout in which the request should last.

	r.Get("/resp", func(w http.ResponseWriter, r *router.Request) {
		w.Write([]byte("responded"))
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
^ you can limit how much the logging middleware reads from the request body, the default size set is 1kb, but you can change it via passing in an int in the Logger() func, an example: "r.Use(router.Logger(1024 * 2)" (limit to 2 kilobytes.) but make sure the size your going to configure is of appropriate type (uint32).
<br>

<br>
<h2 align="center">Example usage with route-specific middleware:</h2>

```
package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/Xrouter"
)

func main() {
	r := router.NewRouter()

	r.Get("/greet", func(w http.ResponseWriter, r *router.Request) {
		w.Write([]byte("hello"))
	}, router.Recoverer(), router.Logger(0)) // you can do route-specific middleware(s) like this (can be 1 or many).

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
<br>

while using the timeout middleware it is also important to note that whilst using it as so, your handler would also need to explicitly coorporate with the middleware inorder for the request to cancel after a given time, so example:

```
package handler

import (
	"context"
)

func SlowHiHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(5 * time.Second):
		// continue with task before cancellation

	case <-r.Context().Done():
		return // cancel
	}
}
```
