Routes.go
=====

a simple router

##Example
```
package main

import (
	"github.com/nathanfaucett/routes"
	"fmt"
)

func main() {
	// creates new Router
	router := routes.NewRouter()
	
	// use some middleware
	// for all routes "/"
	router.Use(
		func(a, b int) {
			// do work
		},
	)
	
	// attach middleware for only the "/carts/:cart_id[0-9]/items" urls
	router.Use(
		"/carts/:cart_id[0-9]/items",
		func(a, b int) {
			// do work
		},
	)
	
	// add some routes
	// restiful routes support GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE
	router.Get("/carts/:cart_id[0-9]/items/:id[0-9](.:format)",
		func(a, b int) {
			fmt.Println(a + b)
		},
		func(a, b int) {
			fmt.Println(a + b)
		},
		func(a, b int) {
			fmt.Println(a + b)
		},
	);
	
	// other ways of adding routes
	
	route := router.Route("/some/path/:param")
	route.Get(func() {})
	route.Post(func() {})
	router.Route("/some/path/:param").Get(func() {})
	
	// find stack that matches method and path
	stack := router.Find("GET", "/carts/1/items/1.json")
	
	// now that we have a stack, loop though stack
	for _, middleware := range route.Stack {
		params := middleware.Params // map[string][string]
		handler := middleware.Handler // function to be called
		
		// handler has 4 methods,
		//   Call(...interface{}) passes arguments to function
		//   Func() returns function as reflect.Value
		//   NumIn() returns number of arguments
		//   NumOut() returns number of return arguments
		handler.Call(1, 2)
		
		// also can grab return values as []reflect.Value
		// values := handler.Call(1, 2)
	}
}
```