package echo

import (
	"net/http"
)

type (
	// Group is a set of sub-routes for a specified route. It can be used for inner
	// routes that share a common middleware or functionality that should be separate
	// from the parent echo instance while still inheriting from it.
	Group struct {
		common
		host       string
		prefix     string
		middleware []MiddlewareFunc
		echo       *Echo
	}
)

// Use implements `Echo#Use()` for sub-routes within the Group.
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
	if len(g.middleware) == 0 {
		return
	}
	// Allow all requests to reach the group as they might get dropped if router
	// doesn't find a match, making none of the group middleware process.
	g.Any("", "-", NotFoundHandler)
	g.Any("/*", "-", NotFoundHandler)
}

// CONNECT implements `Echo#CONNECT()` for sub-routes within the Group.
func (g *Group) CONNECT(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodConnect, path, desc, h, m...)
}

// DELETE implements `Echo#DELETE()` for sub-routes within the Group.
func (g *Group) DELETE(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodDelete, path, desc, h, m...)
}

// GET implements `Echo#GET()` for sub-routes within the Group.
func (g *Group) GET(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodGet, path, desc, h, m...)
}

// HEAD implements `Echo#HEAD()` for sub-routes within the Group.
func (g *Group) HEAD(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodHead, path, desc, h, m...)
}

// OPTIONS implements `Echo#OPTIONS()` for sub-routes within the Group.
func (g *Group) OPTIONS(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodOptions, path, desc, h, m...)
}

// PATCH implements `Echo#PATCH()` for sub-routes within the Group.
func (g *Group) PATCH(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPatch, path, desc, h, m...)
}

// POST implements `Echo#POST()` for sub-routes within the Group.
func (g *Group) POST(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPost, path, desc, h, m...)
}

// PUT implements `Echo#PUT()` for sub-routes within the Group.
func (g *Group) PUT(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPut, path, desc, h, m...)
}

// TRACE implements `Echo#TRACE()` for sub-routes within the Group.
func (g *Group) TRACE(path, desc string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodTrace, path, desc, h, m...)
}

// Any implements `Echo#Any()` for sub-routes within the Group.
func (g *Group) Any(path, desc string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, m := range methods {
		routes[i] = g.Add(m, path, desc, handler, middleware...)
	}
	return routes
}

// Match implements `Echo#Match()` for sub-routes within the Group.
func (g *Group) Match(methods []string, path, desc string, handler HandlerFunc, middleware ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, m := range methods {
		routes[i] = g.Add(m, path, desc, handler, middleware...)
	}
	return routes
}

// Group creates a new sub-group with prefix and optional sub-group-level middleware.
func (g *Group) Group(prefix string, middleware ...MiddlewareFunc) (sg *Group) {
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	sg = g.echo.Group(g.prefix+prefix, m...)
	sg.host = g.host
	return
}

// Static implements `Echo#Static()` for sub-routes within the Group.
func (g *Group) Static(prefix, root, desc string) {
	g.static(prefix, root, desc, g.GET)
}

// File implements `Echo#File()` for sub-routes within the Group.
func (g *Group) File(path, file, desc string) {
	g.file(path, file, desc, g.GET)
}

// Add implements `Echo#Add()` for sub-routes within the Group.
func (g *Group) Add(method, path, desc string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	return g.echo.add(g.host, method, g.prefix+path, handler, desc, m...)
}
