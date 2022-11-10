package streams

// IThen represents a post conditional set of actions that can be used with stream operations.
type IThen[T comparable] interface {
	// Then executes the provided handler if a certain condition is met
	Then(func(s IStream[T])) IThen[T]

	// Else executes the provided handler if a certain condition is not met
	Else(func(s IStream[T])) IThen[T]
}

type thenWrapper[T comparable] struct {
	conditionMet bool
	stream       IStream[T]
}

func (t *thenWrapper[T]) Then(f func(s IStream[T])) IThen[T] {
	if t.conditionMet {
		f(t.stream)
	}
	return t
}

func (t *thenWrapper[T]) Else(f func(s IStream[T])) IThen[T] {
	if !t.conditionMet {
		f(t.stream)
	}
	return t
}
