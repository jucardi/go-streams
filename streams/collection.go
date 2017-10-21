package streams

import (
	"errors"
	"reflect"
)

var (
	ErrorWrongType = errors.New("wrong type")
)

type genericArrayCollection struct {
	reflect.Value
	elementType  reflect.Type
	currentIndex int
}

// Creates a new ICollection from the given collection. Panics if the provided interface is not an collection or slice.
func NewCollectionFromArray(array interface{}) ICollection {
	arrayType := reflect.TypeOf(array)

	if arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.Array {
		panic("Unable to create collection from a none Slice or none Array")
	}

	return &genericArrayCollection{
		Value:       reflect.ValueOf(array),
		elementType: reflect.TypeOf(array).Elem(),
	}
}

// Creates a new empty collection of the given type
func NewCollection(elementType reflect.Type) ICollection {
	return &genericArrayCollection{
		Value:       reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0),
		elementType: elementType,
	}
}

func (g *genericArrayCollection) Index(i int) interface{} {
	if i >= g.Len() {
		return nil
	}
	return g.Value.Index(i).Interface()
}

func (g *genericArrayCollection) Add(item interface{}) error {
	if reflect.PtrTo(reflect.TypeOf(item)).AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	g.Value = reflect.Append(g.Value, reflect.ValueOf(item))
	return nil
}

// AddCollection: Appends another collection into this ICollection instance.
func (g *genericArrayCollection) AddCollection(slice ICollection) {
	g.Value = reflect.AppendSlice(g.Value, reflect.ValueOf(slice.ToArray()))
}

// ElementType: Returns the type of the elements in the collection
func (g *genericArrayCollection) ElementType() reflect.Type {
	return g.elementType
}

// Retrieves the current element of the iterator
func (g *genericArrayCollection) Current() interface{} {
	if g.currentIndex >= g.Len() {
		return nil
	}

	return g.Index(g.currentIndex)
}

// Indicates whether the collection has a next
func (g *genericArrayCollection) HasNext() bool {
	return g.Len()-1 >= g.currentIndex
}

// Moves to the next position and indicates whether the collection has a next element and returns the next element if so
func (g *genericArrayCollection) Next() interface{} {
	g.currentIndex++

	if g.currentIndex < 0 || g.currentIndex >= g.Len() {
		return nil
	}

	return g.Index(g.currentIndex)
}

// Resents the iterator position to the beginning.
func (g *genericArrayCollection) Reset() {
	g.currentIndex = 0
}

// Skips the following N items
func (g *genericArrayCollection) Skip(n int) {
	g.currentIndex += n
}

func (g *genericArrayCollection) ToArray() interface{} {
	return g.Value.Interface()
}
