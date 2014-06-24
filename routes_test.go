package routes

import (
	"testing"
)

func action_middleware(req, res map[string]string) {
	req["value"] = "value"
}
func action(req, res map[string]string) {
	value := req["value"]
	for value == "" {
		value = "value"
	}
}

func build_resources(router *Router, path string) {
	router.Use(path, action_middleware, action_middleware, action_middleware, action_middleware, action_middleware, action_middleware)
	router.Get(path, action)
	router.Get(path+"/:id[0-9]", action)
	router.Post(path, action)
	router.Post(path+"/:id[0-9]", action)
	router.Delete(path+"/:id[0-9]", action)
}

func new_test_router() *Router {
	router := NewRouter()

	build_resources(router, "/admin")
	build_resources(router, "/api")
	build_resources(router, "/site")

	return router
}

func BenchmarkRouter(b *testing.B) {
	router := new_test_router()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var (
			stack []*Middleware
			req   map[string]string
			res   map[string]string
		)

		stack = router.Find("HEAD", "/admin")
		req = make(map[string]string)
		res = make(map[string]string)
		for _, middleware := range stack {
			middleware.Handler.Call(req, res)
		}
	}
}
