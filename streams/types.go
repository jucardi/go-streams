package streams

import "reflect"

// IIterator defines the contract to be used to iterate over an iterable item. They indicate which one is the current item,
// and whether it has a next item or not.
//
// The following example indicates one of the proper way to iterate over an iterator:
//
//     iterator := iterable.Iterator()
//
//     for x := iterator.Current(); iterator.HasNext() && i < end; x = iterator.Next() {
//
type IIterator interface {
	// Retrieves the current element of the iterator
	Current() interface{}
	// Indicates whether the iterable has a next
	HasNext() bool
	// Moves to the next position and returns the next element, if any
	Next() interface{}
	// Resents the iterator position to the beginning.
	Reset() IIterator
	// Skips the following N items
	Skip(n int) IIterator
	// ElementType returns the type of the elements in the iterable
	ElementType() reflect.Type
}

// ICollection represents a collection of elements to be used in a Stream. It can represent a data structure,
// an iterable, a generator function, or an I/O channel, through a pipeline of computational operations.
type ICollection interface {
	IIterable
	// Index returns v's i'th element.
	Index(index int) interface{}
	// Add appends the element into the iterable. Returns error if the item is not the proper type
	Add(item interface{}) error
	// AddAll appends another iterable into this ICollection instance.
	AddAll(slice IIterable)
}

// IIterable represent an iterable of elements to be used in a Stream. By default Collections are considered iterables.
// Iterables do not require to have a defined size. They can represent a collection, a generator function, or an I/O channel.
type IIterable interface {
	// Returns a new iterator for the iterable
	Iterator() IIterator
	// Len returns the size of the iterable if the size is finite and known, otherwise returns -1.
	Len() int
	// ElementType returns the type of the elements in the iterable
	ElementType() reflect.Type
	// Interface returns v's current value as an interface{}.
	ToArray() interface{}
}
