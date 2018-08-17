package streams

import "reflect"

// collectionIterator is the default implementation for an iterator which helps to iterate over an ICollection implementation.
type collectionIterator struct {
	col          ICollection
	currentIndex int
}

// Current retrieves the current element of the collectionIterator.
func (g *collectionIterator) Current() interface{} {
	if g.currentIndex >= g.col.Len() {
		return nil
	}

	return g.col.Index(g.currentIndex)
}

// HasNext indicates whether the iterable has a next element without moving the pointer.
func (g *collectionIterator) HasNext() bool {
	return g.col.Len()-1 >= g.currentIndex
}

// MoveNext moves the pointer of the collectionIterator to the next element of the set. Returns `false` if no more elements are present in the set.
func (g *collectionIterator) MoveNext() bool {
	if !g.HasNext() {
		return false
	}

	g.currentIndex++
	return true
}

// Moves to the next element of the set and returns its value. Returns `nil` if no more elements are present in the set.
func (g *collectionIterator) Next() interface{} {
	if !g.MoveNext() {
		return nil
	}

	return g.Current()
}

// Skips the following N items
func (g *collectionIterator) Skip(n int) IIterator {
	g.currentIndex += n
	return g
}

// ElementType returns the type of the elements in the iterable
func (g *collectionIterator) ElementType() reflect.Type {
	return g.col.ElementType()
}

// Resets the iterator position to the beginning. Not available in the IIterator interface since not all iterators support resetting to the beginning.
func (g *collectionIterator) Reset() IIterator {
	g.currentIndex = 0
	return g
}

func newCollectionIterator(col ICollection) IIterator {
	return &collectionIterator{
		col:          col,
		currentIndex: 0,
	}
}
