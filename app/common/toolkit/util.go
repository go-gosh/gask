package tk

import "reflect"

var errorCheck = func(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](t T, err error) T {
	errorCheck(err)
	return t
}

func Pointer[T any](t T) *T {
	return &t
}

func ZeroNilPointer[T any](t T) *T {
	val := reflect.ValueOf(t)
	if val.IsZero() {
		return nil
	}
	return &t
}

func SetErrorHandler(fn func(err error)) {
	errorCheck = fn
}
