package streams

import "reflect"

// IIterable represent an iterable of elements in a set. By default Collections are considered iterables.
// Iterables do not require to have a defined size. They can represent a collection, a generator function, or an I/O channel.
type IIterable interface {
	// Returns a new collectionIterator for the iterable
	Iterator() IIterator

	// Len returns the size of the iterable if the size is finite and known, otherwise returns -1.
	Len() int

	// ElementType returns the type of the elements in the iterable
	ElementType() reflect.Type

	// ToArray returns an array representation of the iterable
	ToArray(defaultArray ...interface{}) interface{}
}

type IMapIterable interface {
	IIterable

	// ToMap returns a map representation of the IMapIterable
	ToMap() interface{}
}
