package streams

// IIterator defines the contract to be used to iterate over a set.
//
//	Usage:
//
//	    for x := collectionIterator.Current(); collectionIterator.HasNext(); x = collectionIterator.Next() {
//	    }
type IIterator[T any] interface {
	// Current retrieves the current element of the collectionIterator
	Current() T

	// MoveNext moves the pointer of the collectionIterator to the next element of the set. Returns `false` if no more elements are present in the set.
	MoveNext() bool

	// HasNext indicates whether the iterable has a next element without moving the pointer.
	HasNext() bool

	// Next moves to the next element of the set and returns its value.
	// Returns `nil` if no more elements are present in the set.
	Next() T

	// Skip skips the following N items
	Skip(n int) IIterator[T]

	// ForEachRemaining iterates over the elements in the iterator starting from its current position.
	ForEachRemaining(f IterFunc[T])
}

// IIteratorWithPos defines an iterator that keeps track of the position in the list while iterating. Supports stopping at a certain index when iterating
type IIteratorWithPos[T any] interface {
	IIterator[T]
	// StopAt If stopAt is provided and Pos() os implemented by the
	// iterator, this function will stop when the Pos() == stopAt[0], otherwise it will iterate until there are no more elements in the iterator
	StopAt(pos int)
	// Pos indicates the current position in the iterator if the iterator implements it, otherwise returns -1
	Pos() int
}

// IIterable represent an iterable of elements in a set. By default Collections are considered iterables.
// Iterables do not require to have a defined size. They can represent a collection, a generator function, or an I/O channel.
type IIterable[T any] interface {
	// Iterator returns a new collection iterator for the iterable
	Iterator() IIterator[T]

	// ForEach Iterates over the elements in the collection
	ForEach(f IterFunc[T])
}

// ICollection represents a collection of elements of type T. It can represent a data structure, an iterable, a generator function, or an I/O channel, through
// a pipeline of computational operations.
type ICollection[T comparable] interface {
	IIterable[T]

	// Add appends the element into the iterable. Returns error if the item is not the proper type
	Add(item ...T) bool

	// AddFromIterator appends elements from an iterator into this IList instance.
	AddFromIterator(iterator IIterator[T]) bool

	// Remove removes the provided element(s) from the collection. Returns true if there was any match that was removed
	Remove(item ...T) bool

	// RemoveFromIterator removes elements in the collection that match any elements in the provided iterator
	RemoveFromIterator(iterator IIterator[T]) bool

	// RemoveIf removes all the elements of this collection that satisfy the given condition. Returns true if an item was removed, otherwise false
	//
	//   {condition}  -  The condition function that will determine whether to remove items
	//   {keepOrder}  -  (false by default) Optional flag that indicates if the removal of the element should guarantee the order of the remaining elements. In
	//                   some cases, guaranteeing the order of elements after a removal can be a costly operation since the remaining elements have to be shifted
	//                   in the collection.
	RemoveIf(condition ConditionalFunc[T], keepOrder ...bool) bool

	// Contains indicates if this collection contains the provided item(s).
	Contains(item ...T) bool

	// ContainsFromIterator indicates if this collection contains the items contained in the provided iterator.
	ContainsFromIterator(iterator IIterator[T]) bool

	// Len returns the size of the iterable if the size is finite and known, otherwise returns -1.
	Len() int

	// Clear removes all elements from this collection
	Clear()

	// ToArray returns an array containing all the elements in this collection.
	ToArray() []T

	// IsEmpty indicates whether this collection has any elements
	IsEmpty() bool
}

// IList represents a collection of elements of type T, similar to ICollection[T], but it includes operations that require an index
type IList[T comparable] interface {
	ICollection[T]

	// Index returns the value in the position indicated by the index.
	Index(index int) (val T, exists bool)

	// RemoveAt removes the element at the provided index. Returns true if an item was removed, otherwise false
	//
	//   {index}      -  The index of the element to remove
	//   {keepOrder}  -  (false by default) Optional flag that indicates if the removal of the element should guarantee the order of the remaining elements. In
	//                   some cases, guaranteeing the order of elements after a removal can be a costly operation since the remaining elements have to be shifted
	//                   in the list.
	RemoveAt(index int, keepOrder ...bool) bool

	// Distinct returns a set of all unique values in this collection
	Distinct() ISet[T]

	// Stream returns a sequential Stream with this collection as its source.
	Stream() IStream[T]
}

// ISet represents a collection of T with only unique values
type ISet[T comparable] interface {
	ICollection[T]
}

// IStream defines the functions of a stream implementation
type IStream[T comparable] interface {
	// SetThreads Sets the amount of go channels to be used for parallel filtering. Providing a value <= 0, indicates the
	// maximum amount of available CPUs will be the number that determines
	// the amount of go channels to be used. If order matters, best combine it with a `SortBy`. Only needs to be provided once
	// per stream.
	// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
	//            to be used. <= 0 indicates the maximum amount of available CPUs will be the number that determines the
	//            amount of go channels to be used. If order matters,
	//            best combine it with a `SortBy`. Only needs to be provided once per stream.
	SetThreads(threads int) IStream[T]

	// Filter appends a filtering function to the stream, where any element that does not meet the condition provided by
	// the function (return false) will be filtered when processing the stream
	Filter(f ConditionalFunc[T]) IStream[T]

	// Except has the opposite effect than 'Filter'. Appends a filtering function to the stream, where any element that
	// does not meet the condition provided by the function (return true) will be filtered when processing the stream.
	Except(f ConditionalFunc[T]) IStream[T]

	// Sort sorts the elements in the stream using the provided comparable function.
	//
	// - desc:  indicates whether the sorting should be done descendant
	Sort(f SortFunc[T], desc ...bool) IStream[T]

	// Distinct ensures that the finalizing operation of the stream includes only unique elements
	Distinct() IStream[T]

	// First Returns the first element of the resulting stream.
	// Returns default T if the resulting stream is empty (or defaultValue if provided)
	First(defaultValue ...T) T

	// Last Returns the last element of the resulting stream.
	// Returns default T if the resulting stream is empty (or defaultValue if provided)
	Last(defaultValue ...T) T

	// At Returns the element at the given index in the resulting stream.
	// Returns default T if out of bounds (or defaultValue if provided)
	At(index int, defaultValue ...T) T

	// AtReverse Returns the element at the given position, starting from the last element to the first in the resulting stream.
	// Returns default T if out of bounds (or defaultValue if provided)
	AtReverse(pos int, defaultValue ...T) T

	// Count Counts the elements of the resulting stream
	Count() int

	// IsEmpty indicates whether the result of the stream produced no elements
	IsEmpty() bool

	// Contains indicates whether the provided value matches any of the values in the stream
	//
	// - value:   The value to be found.
	Contains(value T) bool

	// AnyMatch Indicates whether any elements of the stream match the given condition function.
	//
	// - f:       The matching function to be used.
	AnyMatch(f ConditionalFunc[T]) bool

	// AllMatch Indicates whether ALL elements of the stream match the given condition function
	//
	// - f:       The matching function to be used.
	AllMatch(f ConditionalFunc[T]) bool

	// NotAllMatch is the negation of `AllMatch`. If any of the elements do not match the provided condition the result
	// will be `true`; `false` otherwise.
	//
	// - f:       The matching function to be used.
	NotAllMatch(f ConditionalFunc[T]) bool

	// NoneMatch indicates whether NONE of elements of the stream match the given condition function.
	//
	// - f:       The matching function to be used.
	NoneMatch(f ConditionalFunc[T]) bool

	// IfEmpty returns a `Then` handler where actions like `Then` or `Else` can be triggered if the stream empty
	IfEmpty() IThen[T]

	// IfAnyMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if any element match the
	// provided condition based on what the result of `AnyMatch` would return
	//
	// - f:       The matching function to be used.
	IfAnyMatch(f ConditionalFunc[T]) IThen[T]

	// IfAllMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if all elements match the
	// provided condition based on what the result of `AllMatch` would return
	//
	// - f:       The matching function to be used.
	IfAllMatch(f ConditionalFunc[T]) IThen[T]

	// IfNoneMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if no elements match the
	// provided condition based on what the result of `NoneMatch` would return
	//
	// - f:       The matching function to be used.
	IfNoneMatch(f ConditionalFunc[T]) IThen[T]

	// ForEach iterates over all elements in the stream calling the provided function.
	ForEach(f IterFunc[T])

	// ParallelForEach Iterates over all elements in the stream calling the provided function. Creates multiple go channels to parallelize
	// the operation. ParallelForeach does not use any thread values previously provided in any filtering method nor enables parallel filtering
	// if any filtering is done prior to the `ParallelForEach` phase. Only use `ParallelForEach` if the order in which the elements are processed
	// does not matter, otherwise see `ForEach`.
	//
	// - threads:   Indicates the amount of go channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates
	//              the maximum amount of available CPUs will be the number that determines the amount of go channels to be used.
	// - skipWait:  Indicates whether `ParallelForEach` will wait until all channels are done processing.
	ParallelForEach(f IterFunc[T], threads int, skipWait ...bool)

	// ToArray Returns an array of elements from the resulting stream
	ToArray() []T

	// ToCollection returns a `ICollection` of elements from the resulting stream
	ToCollection() ICollection[T]

	// ToIterable returns a `IIterable` of elements from the resulting stream
	ToIterable() IIterable[T]

	// ToList returns a `IList` of elements from the resulting stream
	ToList() IList[T]

	// ToDistinct processes the stream and outputs a set of unique values
	ToDistinct() ISet[T]
}

// KeyValuePair is a structure which contains a pair of key-values from a map
type KeyValuePair[K, V comparable] struct {
	Key   K
	Value V
}

// IMap defines the contract for a generic map which also represents a collection of `*KeyValuePairs`
type IMap[K, V comparable] interface {
	IList[*KeyValuePair[K, V]]

	// ToMap returns a map containing the elements contained by this instance
	ToMap() map[K]V

	// Get returns the value at index `key`, or the value mapped to the key `key` if the collection represents a `map`.
	Get(k K) (val V, exists bool)

	// Set is mapCollection specific function that allows a value to be added to the map without having to wrap it in a *KeyValuePair
	Set(key K, value V) bool

	// ContainsKey indicates whether the map contains the specified key
	ContainsKey(k K) bool

	// Keys returns the list of keys contained by the map
	Keys() []K

	// Delete removes the item matching the specified key from the map. Returns false iff the key is not contained by the map
	Delete(k K) bool

	// Clear removes all elements from the map
	Clear()
}

// ConditionalFunc is an alias to `func(interface{}) bool` which serves to define if a condition is met for an element in the collection.
type ConditionalFunc[T comparable] func(T) bool

// ConvertFunc is an alias to `func(interface{}) interface{}` which serves to define a type conversion for the `Map` function.
type ConvertFunc[From, To any] func(From) To

// IterFunc is an alias to `func(T)` which serves to define an iteration over a collection.
type IterFunc[T any] func(T)

// SortFunc is an alias to `func(interface{}, interface{}) int` which serves to define a comparison between two elements in the collection. Used for sorting purposes.
type SortFunc[T comparable] func(T, T) int
