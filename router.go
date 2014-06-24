package routes

import (
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

type Middleware struct {
	Handler *Handler
	Params  map[string]string
}

func NewMiddleware(handler *Handler, params map[string]string) *Middleware {
	this := new(Middleware)
	this.Handler = handler
	this.Params = params

	return this
}

type Router struct {
	*events.EventEmitter

	layers []*Layer
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
		path      string
		layer     *Layer
		testLayer *Layer
		ok        bool
		i         int
	)

	if path, ok = arguments[0].(string); !ok {
		path = "/"
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

	for _, testLayer = range this.layers {
		if testLayer.path == path && testLayer.route == nil {
			layer = testLayer
			break
		}
	}
	if layer == nil {
		layer = NewLayer(path, false, false)
		this.layers = append(this.layers, layer)
	}

	length := len(arguments)
	for ; i < length; i++ {
		layer.stack = append(layer.stack, NewHandler(arguments[i]))
	}

	return this
}

func (this *Router) Route(path string) *Route {
	this.Lock()
	var (
		route *Route
		layer *Layer
	)

	for _, layer = range this.layers {
		if layer.path == path && layer.route != nil {
			route = layer.route
		}
	}
	if route == nil {
		layer = NewLayer(path, false, true)
		route = NewRoute()
		layer.route = route
		this.layers = append(this.layers, layer)
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
func (this *Router) Find(method, path string) []*Middleware {
	this.Lock()
	defer this.Unlock()
	var stack []*Middleware

	for _, layer := range this.layers {
		pass, params := layer.Match(path)

		if pass {
			if route := layer.route; route != nil {
				tmp := method
				if method == "HEAD" && !route.methods["HEAD"] {
					tmp = "GET"
				}

				if !route.methods[tmp] {
					continue
				}
				for _, handler := range route.stack[tmp] {
					stack = append(stack, NewMiddleware(handler, params))
				}
			} else if layer.stack != nil {
				for _, handler := range layer.stack {
					stack = append(stack, NewMiddleware(handler, params))
				}
			}
		}
	}

	return stack
}
