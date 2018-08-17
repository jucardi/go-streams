package streams

import "reflect"

// IIterator defines the contract to be used to iterate over an set.
//
//    Usage:
//
//        for x := collectionIterator.Current(); collectionIterator.HasNext(); x = collectionIterator.Next() {
//        }
//
type IIterator interface {
	// Current retrieves the current element of the collectionIterator
	Current() interface{}

	// MoveNext moves the pointer of the collectionIterator to the next element of the set. Returns `false` if no more elements are present in the set.
	MoveNext() bool

	// HasNext indicates whether the iterable has a next element without moving the pointer.
	HasNext() bool

	// Moves to the next element of the set and returns its value.
	// Returns `nil` if no more elements are present in the set.
	Next() interface{}

	// Skip skips the following N items
	Skip(n int) IIterator

	// ElementType returns the type of the elements in the iterable
	ElementType() reflect.Type
}
