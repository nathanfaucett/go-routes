package routes

import (
	"regexp"
)

var (
	parts_matcher = regexp.MustCompile(`\/+\w+|\/\:\w+(\[.+?\])?|\(.+?\)`)
	part_matcher  = regexp.MustCompile(`(\:?\w+)(\[.+?\])?`)
)

type param struct {
	Name     string
	Required bool
}

func new_param(name string, required bool) *param {
	this := new(param)
	this.Name = name
	this.Required = required

	return this
}

type Route struct {
	Method string
	Path   string
	Params []*param
	Stack  []*handler
	regex  *regexp.Regexp
}

func NewRoute(method, path string, stack []*handler) *Route {
	this := new(Route)
	this.Method = method
	this.Path = path
	this.Stack = stack
	this.compile_regex()

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
	if len(this.Params) == 0 {
		return params
	}

	for i := range this.Params {
		if i < length {
			params[this.Params[i].Name] = result[i+1]
		}
	}

	return params
}

func (this *Route) compile_regex() {
	pattern := "^"
	parts := parts_matcher.FindAllString(this.Path, -1)

	for i := range parts {
		part := parts[i]
		if len(part) <= 0 {
			continue
		}

		if string(part[0]) == "(" {
			pattern += "(?:\\" + string(part[1])
			partParts := part_matcher.FindAllStringSubmatch(part, -1)[0]
			part = partParts[1]

			if string(part[0]) == ":" {
				regex := partParts[2]
				if regex == "" {
					regex = "[a-zA-Z0-9-_]"
				}

				pattern += "(" + regex + "+?)"
				this.Params = append(this.Params, new_param(part[1:], false))
			} else {
				pattern += part
			}

			pattern += ")?"
		} else {
			pattern += "\\" + string(part[0]) + "+"
			partParts := part_matcher.FindAllStringSubmatch(part, -1)[0]
			part = partParts[1]

			if string(part[0]) == ":" {
				regex := partParts[2]
				if regex == "" {
					regex = "[a-zA-Z0-9-_]"
				}

				pattern += "(" + regex + "+)"
				this.Params = append(this.Params, new_param(part[1:], true))
			} else {
				pattern += part
			}
		}
	}

	pattern += "(\\/)?(\\?.+)?$"
	this.regex = regexp.MustCompile(pattern)
}
