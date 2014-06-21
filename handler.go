package routes

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

type handler struct {
	function reflect.Value
	num_in   int
	num_out  int
}

func new_handler(fn interface{}) *handler {
	defer func() {
		if r := recover(); nil != r {
			err := errors.New(fmt.Sprintf("%v", r))
			debug.Error(err)
			os.Exit(1)
		}
	}()
	function := reflect.ValueOf(fn)
	if function.Kind() != reflect.Func {
		panic("Route's handlers must be funtions")
	}
	typeof := function.Type()

	this := new(handler)
	this.function = function
	this.num_in = typeof.NumIn()
	this.num_out = typeof.NumOut()

	return this
}

func (this *handler) NumIn() int {
	return this.num_in
}

func (this *handler) NumOut() int {
	return this.num_out
}

// handler.Call(arguments ...interface{}) []reflect.Value,
//    calls handlers function passing in arguments and returns
//    functions return values as []reflect.Value
func (this *handler) Call(arguments ...interface{}) []reflect.Value {
	defer func() {
		if r := recover(); nil != r {
			err := errors.New(fmt.Sprintf("Route handler %v", r))
			debug.Error(err)
			os.Exit(1)
		}
	}()
	var values []reflect.Value

	length := len(arguments)
	for i := 0; i < length; i++ {
		values = append(values, reflect.ValueOf(arguments[i]))
	}

	return this.function.Call(values)
}
