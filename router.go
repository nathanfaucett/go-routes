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

type Router struct {
	*events.EventEmitter

	middleware []*Middleware
	routes     []*Route
}

// creates new router
func NewRouter() *Router {
	this := new(Router)
	this.EventEmitter = events.NewEventEmitter()

	return this
}

func (this *Router) Use(arguments ...interface{}) *Router {
	this.Lock()
	defer this.Unlock()
	var (
		path            string
		middleware      *Middleware
		test_middleware *Middleware
		ok              bool
		i               int
	)

	if path, ok = arguments[0].(string); !ok {
		path = "/"
		i = 0
	} else {
		i = 1
	}
	if len(path) == 0 {
		path = "/"
	}
	if string(path[0]) != "/" {
		path = "/" + path
	}
	if size := len(path); size > 1 && string(path[size-1]) == "/" {
		path = path[:size-1]
	}

	for _, test_middleware = range this.middleware {
		if test_middleware.path == path {
			middleware = test_middleware
			break
		}
	}
	if middleware == nil {
		middleware = NewMiddleware(path, false)
		this.middleware = append(this.middleware, middleware)
	}

	length := len(arguments)
	for ; i < length; i++ {
		middleware.stack = append(middleware.stack, NewHandler(arguments[i]))
	}

	return this
}

func (this *Router) Route(path string) *Route {
	this.Lock()
	var (
		test_route *Route
		route      *Route
	)

	if len(path) == 0 {
		path = "/"
	}
	if string(path[0]) != "/" {
		path = "/" + path
	}
	if size := len(path); size > 1 && string(path[size-1]) == "/" {
		path = path[:size-1]
	}

	for _, test_route = range this.routes {
		if test_route.path == path {
			route = test_route
			break
		}
	}
	if route == nil {
		route = NewRoute(path, false)
		this.routes = append(this.routes, route)
	}

	this.Unlock()
	this.Emit("mount", route)

	return route
}

func (this *Router) Mount(method, path string, arguments ...interface{}) *Route {
	return this.Route(path).Mount(method, arguments)
}
func (this *Router) Unmount(method, path string, arguments ...interface{}) *Route {
	return this.Route(path).Unmount(method, arguments...)
}
func (this *Router) All(path string, arguments ...interface{}) *Route {
	return this.Route(path).All(arguments...)
}
func (this *Router) Get(path string, arguments ...interface{}) *Route {
	return this.Route(path).Get(arguments...)
}
func (this *Router) Post(path string, arguments ...interface{}) *Route {
	return this.Route(path).Post(arguments...)
}
func (this *Router) Put(path string, arguments ...interface{}) *Route {
	return this.Route(path).Put(arguments...)
}
func (this *Router) Patch(path string, arguments ...interface{}) *Route {
	return this.Route(path).Patch(arguments...)
}
func (this *Router) Update(path string, arguments ...interface{}) *Route {
	return this.Route(path).Update(arguments...)
}
func (this *Router) Head(path string, arguments ...interface{}) *Route {
	return this.Route(path).Head(arguments...)
}
func (this *Router) Options(path string, arguments ...interface{}) *Route {
	return this.Route(path).Options(arguments...)
}
func (this *Router) Delete(path string, arguments ...interface{}) *Route {
	return this.Route(path).Delete(arguments...)
}

// finds handlers that matches path and method
func (this *Router) Find(method, path string) ([]*Handler, map[string]string, error) {
	this.Lock()
	defer this.Unlock()

	for _, route := range this.routes {
		if !route.methods[method] {
			continue
		}
		params := route.Match(path)

		if params != nil {
			debug.Log("Routing " + method + path)
			var stack []*Handler

			for _, middleware := range this.middleware {
				pass := middleware.Match(path)

				if pass {
					stack = append(stack, middleware.stack...)
				}
			}
			stack = append(stack, route.stack[method]...)

			return stack, params, nil
		}
	}

	return nil, nil, errors.New("No Route matches " + path + " " + method)
}
