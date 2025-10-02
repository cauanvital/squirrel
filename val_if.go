package squirrel2

type valIf[T any] struct {
	Value   T
	Include bool
}

func ValIf[T any](value T, include bool) valIf[T] {
	return valIf[T]{Value: value, Include: include}
}
