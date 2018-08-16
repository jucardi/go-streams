package streams

import "reflect"

type arrayIterator struct {
	col          *arrayCollection
	elementType  reflect.Type
	currentIndex int
}

// Current retrieves the current element of the iterator.
func (g *arrayIterator) Current() interface{} {
	if g.currentIndex >= g.col.Len() {
		return nil
	}

	return g.col.Get(g.currentIndex)
}

// HasNext indicates whether the iterable has a next element without moving the pointer.
func (g *arrayIterator) HasNext() bool {
	return g.col.Len()-1 >= g.currentIndex
}

// MoveNext moves the pointer of the iterator to the next element of the set. Returns `false` if no more elements are present in the set.
func (g *arrayIterator) MoveNext() bool {
	if !g.HasNext() {
		return false
	}

	g.currentIndex++
	return true
}

// Moves to the next element of the set and returns its value. Returns `nil` if no more elements are present in the set.
func (g *arrayIterator) Next() interface{} {
	if !g.MoveNext() {
		return nil
	}

	return g.Current()
}

// Skips the following N items
func (g *arrayIterator) Skip(n int) IIterator {
	g.currentIndex += n
	return g
}

// ElementType returns the type of the elements in the iterable
func (g *arrayIterator) ElementType() reflect.Type {
	return g.elementType
}

// Resents the iterator position to the beginning.
func (g *arrayIterator) Reset() IIterator {
	g.currentIndex = 0
	return g
}

func newArrayIterator(col *arrayCollection, t reflect.Type) IIterator {
	return &arrayIterator{
		col:          col,
		elementType:  t,
		currentIndex: 0,
	}
}
