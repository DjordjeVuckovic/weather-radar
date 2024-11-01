package server

import "net/http"

type Group struct {
	prefix     string
	middleware []MiddlewareFunc
	server     *Server
}

func (s *Server) Group(prefix string, middleware ...MiddlewareFunc) *Group {
	g := &Group{
		prefix: prefix,
		server: s,
	}
	g.Use(middleware...)
	return g
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
	if len(g.middleware) == 0 {
		return
	}
	// group level middlewares are different from Echo `Pre` and `Use` middlewares (those are global). Group level middlewares
	// are only executed if they are added to the Router with route.
	// So we register catch all route (404 is a safe way to emulate route match) for this group and now during routing the
	// Router would find route to match our request path and therefore guarantee the middleware(s) will get executed.
	// TODO: Implement NotFoundHandler
	//g.RouteNotFound("", NotFoundHandler)
	//g.RouteNotFound("/*", NotFoundHandler)
}

func (g *Group) handle(method, path string, handler HandlerFunc) {
	fullPath := g.prefix + path // Combine group prefix with the path

	finalHandler := handler
	for i := len(g.middleware) - 1; i >= 0; i-- {
		finalHandler = g.middleware[i](finalHandler)
	}

	// Apply server-level middleware
	for i := len(g.server.middleware) - 1; i >= 0; i-- {
		finalHandler = g.server.middleware[i](finalHandler)
	}

	// Wrap final handler to check HTTP method
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		finalHandler(w, r)
	})

	g.server.handle(fullPath, httpHandler)
}
