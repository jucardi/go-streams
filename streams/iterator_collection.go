package streams

var (
	_ IIterator[string] = (*collectionIterator[string])(nil)
)

type collectionIterator[T comparable] struct {
	*IndexBasedIterator[T]
	col IList[T]
}

func (iter *collectionIterator[T]) Current() (ret T) {
	if iter.currentIndex >= iter.col.Len() {
		return
	}

	ret, _ = iter.col.Index(iter.currentIndex)
	return
}

func (iter *collectionIterator[T]) Len() int {
	return iter.col.Len()
}

func newCollectionIterator[T comparable](col ...IList[T]) IIterator[T] {
	ret := &collectionIterator[T]{}
	base := &IndexBasedIterator[T]{
		IAbstractIndexBasedIterator: ret,
	}
	ret.IndexBasedIterator = base
	if len(col) > 0 {
		ret.col = col[0]
	} else {
		ret.col = NewList[T]()
	}
	return ret
}
