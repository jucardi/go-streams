package streams

// ICollection represents a collection of elements to be used in a Stream. It can represent a data structure,
// an iterable, a generator function, or an I/O channel, through a pipeline of computational operations.
type ICollection interface {
	IIterable

	// Index returns the value in the position indicated by the index.
	//
	//   - index:  The index of the element to be retrieved.
	//
	Index(index int) interface{}

	// Add appends the element into the iterable. Returns error if the item is not the proper type
	//
	//   - item:  The item to be added to the collection.
	//
	Add(item interface{}) error

	// AddAll appends another iterable into this ICollection instance.
	//
	//   - iterable:  The iterable of elements to be added to the collection.
	//
	AddAll(iterable IIterable) error

	// Remove removes the element at the provided index. Returns the removed item or `nil` if no item was found in that position.
	//
	//   - index:      The index of the element to remove
	//   - keepOrder:  (false by default) Optional flag that indicates if the removal of the element should guarantee the order of the remaining elements. In some cases,
	//                 guaranteeing the order of elements after a removal can me a costly operation since the remaining elements have to be shifted in the collection.
	//
	Remove(index int, keepOrder ...bool) interface{}
}

// IMapCollection represents a collection of `*KeyValuePairs` tied to a `map`
type IMapCollection interface {
	ICollection

	// ToMap returns a map representation of the IMapIterable
	ToMap() interface{}

	// Get returns the value at index `key`, or the value mapped to the key `key` if the collection represents a `map`.
	Get(key interface{}) interface{}

	// Set is mapCollection specific function that allows a value to be added to the map without having to wrap it in a *KeyValuePair
	Set(key, value interface{}) error
}
