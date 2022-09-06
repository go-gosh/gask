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

func Map[T, U any](data []T, fn func(*T) (*U, error)) ([]U, error) {
	res := make([]U, 0, len(data))
	for i := 0; i < len(data); i++ {
		u, err := fn(&data[i])
		if err != nil {
			return nil, err
		}
		res = append(res, *u)
	}
	return res, nil
}
