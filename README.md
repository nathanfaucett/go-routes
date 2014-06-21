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
	
	// add some routes
	router.Post("/:name[a-zA-z-_](.:format)")
	router.Put("/carts(.:format)")
	router.Delete("/carts/:cart_id[0-9](.:format)")
	router.Patch("/carts/:cart_id[0-9]/items(.:format)")
	router.Head("/carts/:cart_id[0-9]/items(.:format)")
	router.Options("/carts/:cart_id[0-9]/items(.:format)")
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
	
	// find route with method and path
	err, route, params := router.Find("GET", "/carts/1/items/1.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	// now that we have a route loop though stack calling them
	for _, handler := range route.Stack {
	
		// route handler has 3 methods,
		//   Call(...interface{}) passes arguments to function
		//   NumIn() returns number of arguments
		//   NumOut() returns number of return arguments
		handler.Call(1, 2)
		
		// also can grab return values as []reflect.Value
		// values := handler.Call(1, 2)
	}
	
	fmt.Println(out[0].Int(), params)
}
```