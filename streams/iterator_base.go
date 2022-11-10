package streams

type IAbstractIndexBasedIterator[T any] interface {
	Current() T
	Len() int
}

type IndexBasedIterator[T any] struct {
	IAbstractIndexBasedIterator[T]
	currentIndex int
	stopAt       int
}

func (iter *IndexBasedIterator[T]) MoveNext() bool {
	if !iter.HasNext() {
		return false
	}

	iter.currentIndex++
	return true
}

func (iter *IndexBasedIterator[T]) HasNext() bool {
	return iter.Len()-1 >= iter.currentIndex
}

func (iter *IndexBasedIterator[T]) Next() (ret T) {
	if !iter.MoveNext() {
		return
	}

	return iter.Current()
}

func (iter *IndexBasedIterator[T]) Skip(n int) IIterator[T] {
	iter.currentIndex += n
	return iter
}

func (iter *IndexBasedIterator[T]) ForEachRemaining(f IterFunc[T]) {
	for val := iter.Current(); iter.HasNext() && (iter.stopAt <= 0 || iter.stopAt >= iter.currentIndex); val = iter.Next() {
		f(val)
	}
}

func (iter *IndexBasedIterator[T]) StopAt(pos int) {
	iter.stopAt = pos
}

func (iter *IndexBasedIterator[T]) Pos() int {
	return iter.currentIndex
}
