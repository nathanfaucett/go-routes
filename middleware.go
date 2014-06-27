package routes

import (
	"github.com/nathanfaucett/ptor"
	"regexp"
)

type Middleware struct {
	stack []*Handler
	path  string
	regex *regexp.Regexp
}

func NewMiddleware(path string, sensitive bool) *Middleware {
	this := new(Middleware)
	this.path = path
	this.regex, _ = ptor.PathToRegexp(path, sensitive, false)

	return this
}

func (this *Middleware) Match(path string) bool {

	return this.regex.MatchString(path)
}
