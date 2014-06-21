package routes

import (
	"errors"
	"github.com/nathanfaucett/debugger"
	"github.com/nathanfaucett/events"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	DELETE  = "DELETE"
)

var debug = debugger.Debug("Router")

type cache_route struct {
	route  *Route
	params map[string]string
}

func new_cache_route(route *Route, params map[string]string) *cache_route {
	this := new(cache_route)

	this.route = route
	this.params = params

	return this
}

type Router struct {
	*events.EventEmitter

	routes     map[string][]*Route
	routeCache map[string]*cache_route
}

// creates new router
func NewRouter() *Router {
	this := new(Router)
	this.EventEmitter = events.NewEventEmitter()

	this.routes = make(map[string][]*Route)
	this.routeCache = make(map[string]*cache_route)

	return this
}

func (this *Router) mount(method, path string, stack []interface{}) *Route {
	this.Lock()

	var route_stack []*handler
	for _, fn := range stack {
		route_stack = append(route_stack, new_handler(fn))
	}

	route := NewRoute(method, path, route_stack)
	this.routes[method] = append(this.routes[method], route)
	this.Unlock()

	this.Emit("mount", route)
	return route
}

func (this *Router) Unmount(method, path string) *Route {
	this.Lock()

	for i, route := range this.routes[method] {
		if route.Path == path {
			this.routes[method] = append(this.routes[method][:i], this.routes[method][i+1:]...)
			this.Unlock()

			this.Emit("unmount", route)
			return route
		}
	}

	return nil
}

func (this *Router) Get(path string, stack ...interface{}) *Route {
	return this.mount(GET, path, stack)
}

func (this *Router) Post(path string, stack ...interface{}) *Route {
	return this.mount(POST, path, stack)
}

func (this *Router) Put(path string, stack ...interface{}) *Route {
	return this.mount(PUT, path, stack)
}

func (this *Router) Patch(path string, stack ...interface{}) *Route {
	return this.mount(PATCH, path, stack)
}

func (this *Router) Update(path string, stack ...interface{}) *Route {
	this.mount(PUT, path, stack)
	return this.mount(PATCH, path, stack)
}

func (this *Router) Head(path string, stack ...interface{}) *Route {
	return this.mount(HEAD, path, stack)
}

func (this *Router) Options(path string, stack ...interface{}) *Route {
	return this.mount(OPTIONS, path, stack)
}

func (this *Router) Delete(path string, stack ...interface{}) *Route {
	return this.mount(DELETE, path, stack)
}

// finds route that matches path, if none found returns error
func (this *Router) Find(method, path string) (error, *Route, map[string]string) {
	this.Lock()
	defer this.Unlock()

	cacheKey := method + path
	cacheRoute := this.routeCache[cacheKey]

	if cacheRoute != nil {
		debug.Log("Found Cached Route " + cacheKey)
		return nil, cacheRoute.route, cacheRoute.params
	}

	for _, route := range this.routes[method] {
		params := route.Match(path)

		if params != nil {
			debug.Log("Found New Route " + cacheKey)
			this.routeCache[cacheKey] = new_cache_route(route, params)

			return nil, route, params
		}
	}

	return errors.New("No Route Matches " + path), nil, nil
}
