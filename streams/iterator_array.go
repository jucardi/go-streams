package streams

var (
	_ IIterator[any] = (*arrayIterator[any])(nil)
)

type arrayIterator[T any] struct {
	*IndexBasedIterator[T]
	arr []T
}

func (a *arrayIterator[T]) Current() (ret T) {
	if a.currentIndex >= len(a.arr) {
		return
	}

	return a.arr[a.currentIndex]
}

func (a *arrayIterator[T]) Len() int {
	return len(a.arr)
}

func newArrayIterator[T any](arr ...[]T) IIterator[T] {
	ret := &arrayIterator[T]{}
	base := &IndexBasedIterator[T]{
		IAbstractIndexBasedIterator: ret,
	}
	ret.IndexBasedIterator = base
	if len(arr) > 0 {
		ret.arr = arr[0]
	}
	return ret
}
