package routes

import (
	"github.com/nathanfaucett/ptor"
	"reflect"
	"regexp"
)

var (
	methods = []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"HEAD",
		"OPTIONS",
		"DELETE",
	}
)

type Route struct {
	path    string
	regex   *regexp.Regexp
	params  []*ptor.Param
	stack   map[string][]*Handler
	methods map[string]bool
}

func NewRoute(path string, sensitive bool) *Route {
	this := new(Route)

	this.path = path
	this.regex, this.params = ptor.PathToRegexp(path, sensitive, true)

	this.stack = make(map[string][]*Handler)
	this.methods = make(map[string]bool)

	return this
}

func (this *Route) Match(path string) map[string]string {
	test := this.regex.FindAllStringSubmatch(path, -1)
	if len(test) == 0 {
		return nil
	}
	result := test[0]
	length := len(result)
	if length == 0 {
		return nil
	}

	params := make(map[string]string)
	if len(this.params) == 0 {
		return params
	}

	for i := range this.params {
		if i < length {
			params[this.params[i].Name] = result[i+1]
		}
	}

	return params
}

func (this *Route) Mount(method string, arguments []interface{}) *Route {
	has := false
	for _, m := range methods {
		if method == m {
			has = true
		}
	}
	if !has {
		panic("Router does not support method " + method)
	}

	for _, handler := range arguments {
		if !this.methods[method] {
			this.methods[method] = true
		}
		this.stack[method] = append(this.stack[method], NewHandler(handler))
	}

	return this
}

func (this *Route) Unmount(method string, arguments ...interface{}) *Route {
	has := false
	for _, m := range methods {
		if method == m {
			has = true
		}
	}
	if !has {
		panic("Router does not support method " + method)
	}

	for i, handler := range this.stack[method] {
		for _, function := range arguments {
			fn := reflect.ValueOf(function)

			if fn == handler.Func() {
				this.stack[method] = append(this.stack[method][:i], this.stack[method][i+1:]...)
			}
		}
	}
	if len(this.stack[method]) == 0 {
		this.methods[method] = false
	}

	return this
}

func (this *Route) All(arguments ...interface{}) *Route {
	for _, handler := range arguments {
		fn := NewHandler(handler)

		this.stack[GET] = append(this.stack[GET], fn)
		this.stack[POST] = append(this.stack[POST], fn)
		this.stack[PUT] = append(this.stack[PUT], fn)
		this.stack[PATCH] = append(this.stack[PATCH], fn)
		this.stack[HEAD] = append(this.stack[HEAD], fn)
		this.stack[OPTIONS] = append(this.stack[OPTIONS], fn)
		this.stack[DELETE] = append(this.stack[DELETE], fn)
	}

	this.methods[GET] = true
	this.methods[POST] = true
	this.methods[PUT] = true
	this.methods[PATCH] = true
	this.methods[HEAD] = true
	this.methods[OPTIONS] = true
	this.methods[DELETE] = true

	return this
}
func (this *Route) Get(arguments ...interface{}) *Route {
	return this.Mount(GET, arguments)
}
func (this *Route) Post(arguments ...interface{}) *Route {
	return this.Mount(POST, arguments)
}
func (this *Route) Put(arguments ...interface{}) *Route {
	return this.Mount(PUT, arguments)
}
func (this *Route) Patch(arguments ...interface{}) *Route {
	return this.Mount(PATCH, arguments)
}
func (this *Route) Update(arguments ...interface{}) *Route {
	this.Mount(PUT, arguments)
	return this.Mount(PATCH, arguments)
}
func (this *Route) Head(arguments ...interface{}) *Route {
	return this.Mount(HEAD, arguments)
}
func (this *Route) Options(arguments ...interface{}) *Route {
	return this.Mount(OPTIONS, arguments)
}
func (this *Route) Delete(arguments ...interface{}) *Route {
	return this.Mount(DELETE, arguments)
}
