package streams

// IThen represents a post conditional set of actions that can be used with stream operations.
type IThen interface {
	// Then executes the provided handler if a certain condition is met
	Then(func())

	// Else executes the provided handler if a certain condition is not met
	Else(func())
}

type thenWrapper struct {
	conditionMet bool
}

func (t *thenWrapper) Then(f func()) {
	if t.conditionMet {
		f()
	}
}

func (t *thenWrapper) Else(f func()) {
	if !t.conditionMet {
		f()
	}
}
