package tk

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
