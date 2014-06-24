package routes

import (
	"errors"
	"reflect"
)

var (
	ErrorInvalidArgument = errors.New("Invalid Argument handler is not a Function")
)

type Handler struct {
	function  reflect.Value
	num_in    int
	num_out   int
	arguments []reflect.Type
}

func NewHandler(function interface{}) *Handler {
	fn := reflect.ValueOf(function)

	if fn.Kind() != reflect.Func {
		panic(ErrorInvalidArgument)
	}

	var arguments []reflect.Type

	typeof := fn.Type()
	length := typeof.NumIn()
	for i := 0; i < length; i++ {
		arguments = append(arguments, typeof.In(i))
	}

	this := new(Handler)
	this.function = fn
	this.arguments = arguments
	this.num_in = typeof.NumIn()
	this.num_out = typeof.NumOut()

	return this
}

func (this *Handler) Func() reflect.Value {
	return this.function
}

func (this *Handler) NumIn() int {
	return this.num_in
}

func (this *Handler) NumOut() int {
	return this.num_out
}

// calls handler's function passing in arguments and returns
// functions return values as []reflect.Value
func (this *Handler) Call(arguments ...interface{}) []reflect.Value {
	var values []reflect.Value

	for i, argument := range arguments {
		if argument == nil {
			values = append(values, reflect.Zero(this.arguments[i]))
		} else {
			values = append(values, reflect.ValueOf(argument))
		}
	}

	return this.function.Call(values)
}
