package streams

import (
	"errors"
)

var (
	// ErrorWrongType is returned when the wrong type of object is passed
	ErrorWrongType = errors.New("wrong type of element")
)

type KeyValuePair struct {
	Key   interface{}
	Value interface{}
}

// ConditionalFunc is an alias to `func(interface{}) bool` which serves to define if a condition is met for an element in the collection.
type ConditionalFunc func(interface{}) bool

// ConvertFunc is an alias to `func(interface{}) interface{}` which serves to define a type conversion for the `Map` function.
type ConvertFunc func(interface{}) interface{}

// IterFunc is an alias to `func(interface{})` which serves to define an iteration over a collection.
type IterFunc func(interface{})

// SortFunc is an alias to `func(interface{}, interface{}) int` which serves to define a comparison between two elements in the collection. Used for sorting purposes.
type SortFunc func(interface{}, interface{}) int
