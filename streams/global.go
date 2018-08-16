package streams

import (
	"fmt"
	"reflect"
)

// From Creates a ArrayStream from a given iterable or ICollection.  Panics if the value is not an array, slice, map or IIterable
//
// - set:      The iterable or ICollection to be used to create the stream
// - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//             to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//             available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//             best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func From(set interface{}, threads ...int) IStream {
	colReflectType := reflect.TypeOf((*IIterable)(nil)).Elem()

	if reflect.PtrTo(reflect.TypeOf(set)).Implements(colReflectType) {
		return FromIterable(set.(IIterable), threads...)
	}

	t := reflect.TypeOf(set)

	switch t.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return FromArray(set, threads...)
	case reflect.Map:
		return FromMap(set, threads...)
	default:
		panic("unknown type, streams may only be created from arrays, slices, maps and IIterable implementations")
	}
}

// FromArray Creates a ArrayStream from a given array.  Panics if the value is not an array or slice.
//
// - array:    The array to be used to create the stream
// - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//             to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//             available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//             best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func FromArray(array interface{}, threads ...int) IStream {
	col, err := NewCollectionFromArray(array)
	if err != nil {
		panic(err)
	}
	return FromIterable(col, threads...)
}

// FromArray Creates a ArrayStream from a given array.  Panics if the value is not an array or slice.
//
// - array:    The array to be used to create the stream
// - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//             to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//             available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//             best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func FromMap(m interface{}, threads ...int) IStream {
	col, err := NewKeyValueSetCollection(m)
	if err != nil {
		panic(err)
	}
	return FromIterable(col, threads...)
}

// FromIterable Creates a ArrayStream from a given IIterable.
//
// - iterable: The ICollection to be used to create the stream
// - threads:  If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//             to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//             available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//             best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func FromIterable(iterable IIterable, threads ...int) IStream {
	return &ArrayStream{
		iterable: iterable,
		threads:  getCores(threads...),
	}
}

// NewCollectionFromArray Creates a new ICollection from the given array or slice.
//
// - array:  The array to be used to create the collection
//
func NewCollectionFromArray(array interface{}) (ICollection, error) {
	arrayType := reflect.TypeOf(array)

	if arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.Array {
		return nil, fmt.Errorf("unable to create collection, the input value is not a slice or array, %s", arrayType.Kind().String())
	}

	return &arrayCollection{
		v:           reflect.ValueOf(array),
		elementType: reflect.TypeOf(array).Elem(),
	}, nil
}

// NewCollectionFromArray Creates a new ICollection from the given array or slice.
//
// - m:  The array to be used to create the collection
//
func NewKeyValueSetCollection(m interface{}) (ICollection, error) {
	val := reflect.ValueOf(m)

	if val.Kind() != reflect.Map {
		return nil, fmt.Errorf("unable to create a key value set collection, the input value must be a map, %s", val.Kind().String())
	}

	var array []*KeyValuePair
	for _, key := range val.MapKeys() {
		array = append(array, &KeyValuePair{
			Key:   key.Interface(),
			Value: val.MapIndex(key),
		})
	}

	return NewCollectionFromArray(array)
}

// NewArrayCollection Creates a new empty array collection of the given type
//
// - elementType:  The element type for the items in the collection to be created.
//
func NewArrayCollection(elementType reflect.Type) ICollection {
	return &arrayCollection{
		v:           reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0),
		elementType: elementType,
	}
}
