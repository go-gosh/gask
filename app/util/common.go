package util

func Point[T any](t T) *T {
	return &t
}

func NilToDefault[T any](t *T, def T) T {
	if t == nil {
		return def
	}
	return *t
}
