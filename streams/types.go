package streams

import "reflect"

type IIterable interface {
	// Retrieves the current element of the iterator
	Current() interface{}
	// Indicates whether the collection has a next
	HasNext() bool
	// Moves to the next position and returns the next element if any
	Next() interface{}
	// Resents the iterator position to the beginning.
	Reset()
	// Skips the following N items
	Skip(n int)
	// ElementType: Returns the type of the elements in the collection
	ElementType() reflect.Type
}

// ICollection represents a collection of elements to be used in a Stream. It can represent a data structure,
// an collection, a generator function, or an I/O channel, through a pipeline of computational operations.
type ICollection interface {
	IIterable
	// Len: Returns the length of the collection
	Len() int
	// Index: Returns v's i'th element.
	Index(index int) interface{}
	// Add: Appends the element into the collection. Returns error if the item is not the proper type
	Add(item interface{}) error
	// AddCollection: Appends another collection into this ICollection instance.
	AddCollection(slice ICollection)
	// Interface: Returns v's current value as an interface{}.
	ToArray() interface{}
}
