package routes

import (
	"github.com/nathanfaucett/ptor"
	"regexp"
)

type Layer struct {
	stack  []*Handler
	route  *Route
	params []*ptor.Param
	path   string
	regex  *regexp.Regexp
}

func NewLayer(path string, sensitive, end bool) *Layer {
	this := new(Layer)
	this.path = path
	this.regex, this.params = ptor.PathToRegexp(path, sensitive, end)

	return this
}

func (this *Layer) Match(path string) (bool, map[string]string) {
	test := this.regex.FindAllStringSubmatch(path, -1)
	if len(test) == 0 {
		return false, nil
	}
	result := test[0]
	length := len(result)
	if length == 0 {
		return false, nil
	}

	pass := false
	for _, value := range result {
		if len(value) != 0 {
			pass = true
		}
	}
	if !pass {
		return false, nil
	}

	params := make(map[string]string)
	if len(this.params) == 0 {
		return true, params
	}

	for i := range this.params {
		if i < length {
			params[this.params[i].Name] = result[i+1]
		}
	}

	return true, params
}
