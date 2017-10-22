package streams

import "reflect"

type IIterator interface {
	// Retrieves the current element of the iterator
	Current() interface{}
	// Indicates whether the iterable has a next
	HasNext() bool
	// Moves to the next position and returns the next element if any
	Next() interface{}
	// Resents the iterator position to the beginning.
	Reset() IIterator
	// Skips the following N items
	Skip(n int) IIterator
	// ElementType: Returns the type of the elements in the iterable
	ElementType() reflect.Type
}

// ICollection represents a iterable of elements to be used in a Stream. It can represent a data structure,
// an iterable, a generator function, or an I/O channel, through a pipeline of computational operations.
type ICollection interface {
	IIterable
	// Index: Returns v's i'th element.
	Index(index int) interface{}
	// Add: Appends the element into the iterable. Returns error if the item is not the proper type
	Add(item interface{}) error
	// AddAll: Appends another iterable into this ICollection instance.
	AddAll(slice IIterable)
}

type IIterable interface {
	// Returns a new iterator for the iterable
	Iterator() IIterator
	// Len: Returns the size of the iterable if the size is finite and known, otherwise returns -1.
	Len() int
	// ElementType: Returns the type of the elements in the iterable
	ElementType() reflect.Type
	// Interface: Returns v's current value as an interface{}.
	ToArray() interface{}
}
