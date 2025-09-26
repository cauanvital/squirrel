package squirrel

type ValIf[T any] struct {
	Value   T
	Include bool
}
