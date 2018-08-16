package streams

// ICollection represents a collection of elements to be used in a ArrayStream. It can represent a data structure,
// an iterable, a generator function, or an I/O channel, through a pipeline of computational operations.
type ICollection interface {
	IIterable
	// Get returns the value at index `key`, or the value mapped to the key `key` if the collection represents a `map`.
	Get(key interface{}) interface{}
	// Add appends the element into the iterable. Returns error if the item is not the proper type
	Add(item interface{}) error
	// AddAll appends another iterable into this ICollection instance.
	AddAll(slice IIterable) error
}

// IMapCollection represents a collection of `*KeyValuePairs` tied to a `map`
type IMapCollection interface {
	ICollection

	// ToMap returns a map representation of the IMapIterable
	ToMap() interface{}

	// Set is mapCollection specific function that allows a value to be added to the map without having to wrap it in a *KeyValuePair
	Set(key, value interface{}) error
}
