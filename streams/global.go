package streams

// From Creates a Stream from a given iterable or IList.  Panics if the value is not an array, slice, map or IIterable
//
//   - set:      The iterable or IList to be used to create the stream
//   - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//     to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//     available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//     best combine it with a `SortBy`. Only needs to be provided once per stream.
func From[T comparable](set any, threads ...int) (ret IStream[T]) {
	switch val := set.(type) {
	case []T:
		return FromArray(val, threads...)
	case ICollection[T]:
		return FromCollection(val, threads...)
	}
	panic("invalid source to create a stream")
}

func FromMap[K, V comparable](set any, threads ...int) (ret IStream[*KeyValuePair[K, V]]) {
	switch val := set.(type) {
	case []*KeyValuePair[K, V]:
		return FromArray(val, threads...)
	case ICollection[*KeyValuePair[K, V]]:
		return FromCollection(val, threads...)
	case map[K]V:
		// default:
		col := NewMap[K, V](val)
		return FromCollection[*KeyValuePair[K, V]](col)
	}
	panic("invalid source to create a stream")
}

// FromArray Creates a Stream from a given array.  Panics if the input is not an array or slice.
//
//   - array:    The array to be used to create the stream
//   - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//     to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//     available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//     best combine it with a `SortBy`. Only needs to be provided once per stream.
func FromArray[T comparable](array []T, threads ...int) IStream[T] {
	col := NewList[T](array)

	return FromCollection[T](col, threads...)
}

// FromCollection Creates a Stream from a given IIterable.
//
//   - iterable: The IList to be used to create the stream
//   - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//     to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//     available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//     best combine it with a `SortBy`. Only needs to be provided once per stream.
func FromCollection[T comparable](iterable ICollection[T], threads ...int) IStream[T] {
	return &Stream[T]{
		iterable: iterable,
		threads:  getCores(threads...),
	}
}

// NewList Creates a new empty array collection of the given type
func NewList[T comparable](arr ...[]T) IList[T] {
	ret := &arrayCollection[T]{}
	base := &CollectionBase[T]{}
	base.SetAbstract(ret)
	ret.CollectionBase = base
	if len(arr) > 0 {
		ret.arr = arr[0]
	}
	return ret
}

// NewMap creates a new map collection of the given type
func NewMap[K, V comparable](m ...map[K]V) IMap[K, V] {
	ret := &mapCollection[K, V]{}
	base := &CollectionBase[*KeyValuePair[K, V]]{}
	base.SetAbstract(ret)
	ret.CollectionBase = base
	if len(m) > 0 {
		ret.m = m[0]
	} else {
		ret.m = map[K]V{}
	}
	_ = ret.Keys()
	return ret
}

// NewIterator creates an iterator of T using the provided source.
//
//	{source}  -  The source to read elements from. This function accepts the following sources where T is comparable:
//	                - []T
//	                - IIterable[T]
//	                - ICollection[T]
//	                - IList[T]
//	                - IIterator[T]
//
// panics for any other source type
func NewIterator[T comparable](source any) (ret IIterator[T]) {
	switch src := source.(type) {
	case []T:
		return newArrayIterator[T](src)
	case IList[T]:
		return newCollectionIterator[T](src)
	}
	panic("invalid source to create a stream")
}

// Map maps the elements of the source to a new element, using the mapping function provided. Outputs a new collection
// with collection the new elements.
//
//	{source}  -  The source to read elements from. This function accepts the following sources where From and To are
//	             comparable:
//	                - []From
//	                - IIterable[From]
//	                - ICollection[From]
//	                - IList[From]
//	                - IIterator[From]
//	                - IStream
//
// panics for any other source type
//
// NOTE: On previous versions of this library before generics, this function used to be a part of IStream. Unfortunately
// golang generics do not allow passing generics to functions that are part of a structure or interface, so this function
// had to be moved out of IStream so the target type of the mapping function could be passed as a generic argument.
func Map[From, To comparable](source any, f ConvertFunc[From, To]) IList[To] {
	switch src := source.(type) {
	case []From:
		return NewList[To](mapIterable[From, To](NewList[From](src), f))
	case IIterable[From]:
		return NewList[To](mapIterable[From, To](src, f))
	case IIterator[From]:
		return NewList[To](mapIterator[From, To](src, f))
	case IStream[From]:
		return mapStream[From, To](src, f)
	}
	panic("invalid mapping source")
}

// MapNonComparable is similar to Map, maps the elements of the source to a new element, using the mapping function
// provided. Outputs an array with collection the new elements instead of a collection and the source accepts
// non-comparable types.
//
//	{source}  -  The source to read elements from. This function accepts the following sources where From and To
//	             accept any type including non-comparable:
//	                - []From
//	                - IIterable[From]
//	                - ICollection[From]
//	                - IList[From]
//	                - IIterator[From]
//
// panics for any other source type
func MapNonComparable[From, To any](source any, f ConvertFunc[From, To]) []To {
	switch src := source.(type) {
	case []From:
		return mapIterator[From, To](newArrayIterator[From](src), f)
	case IIterable[From]:
		return mapIterable[From, To](src, f)
	case IIterator[From]:
	}
	panic("invalid mapping source")
}

// MapToPtr converts a collection of T to a collection of *T. Many structs are not `comparable` which makes T unsupported
// by IStream, ICollection, IList and ISet. Pointers however are considered comparable, so this function outputs an array
// of *T which can be used in the types mentioned.
//
//	{source}  -  The source to read elements from. This function accepts the following sources where From and To
//	             accept any type including non-comparable:
//	                - []From
//	                - IIterable[From]
//	                - ICollection[From]
//	                - IList[From]
//	                - IIterator[From]
//
// panics for any other source type
func MapToPtr[T any](source any) []*T {
	return MapNonComparable[T, *T](source,
		func(obj T) *T {
			ret := &obj
			return ret
		})
}

// Mappers returns a handler which providers predefined ConvertFunc mappers for different value types that can be used
// mapping functions such as `Map[F, T]` and `MapNonComparable[F, T]`
func Mappers() IMappers {
	return defaultMappers
}
