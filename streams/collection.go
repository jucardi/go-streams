package streams

import (
	"errors"
	"reflect"
)

var (
	// ErrorWrongType is returned when the wrong type of object is passed
	ErrorWrongType = errors.New("wrong type")
)

type genericArrayCollection struct {
	reflect.Value
	elementType reflect.Type
}

type genericArrayIterator struct {
	col          *genericArrayCollection
	elementType  reflect.Type
	currentIndex int
}

// NewCollectionFromArray Creates a new ICollection from the given iterable. Panics if the provided interface is not an iterable or slice.
func NewCollectionFromArray(array interface{}) ICollection {
	arrayType := reflect.TypeOf(array)

	if arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.Array {
		panic("Unable to create iterable from a none Slice or none Array")
	}

	return &genericArrayCollection{
		Value:       reflect.ValueOf(array),
		elementType: reflect.TypeOf(array).Elem(),
	}
}

// NewCollection Creates a new empty iterable of the given type
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

// AddAll: Appends another iterable into this ICollection instance.
func (g *genericArrayCollection) AddAll(slice IIterable) {
	g.Value = reflect.AppendSlice(g.Value, reflect.ValueOf(slice.ToArray()))
}

// ElementType: Returns the type of the elements in the iterable
func (g *genericArrayCollection) ElementType() reflect.Type {
	return g.elementType
}

func (g *genericArrayCollection) Iterator() IIterator {
	return &genericArrayIterator{
		col:         g,
		elementType: g.elementType,
	}
}

// ElementType: Returns the type of the elements in the iterable
func (g *genericArrayIterator) ElementType() reflect.Type {
	return g.elementType
}

func (g *genericArrayCollection) ToArray() interface{} {
	return g.Value.Interface()
}

// Retrieves the current element of the iterator
func (g *genericArrayIterator) Current() interface{} {
	if g.currentIndex >= g.col.Len() {
		return nil
	}

	return g.col.Index(g.currentIndex)
}

// Indicates whether the iterable has a next
func (g *genericArrayIterator) HasNext() bool {
	return g.col.Len()-1 >= g.currentIndex
}

// Moves to the next position and indicates whether the iterable has a next element and returns the next element if so
func (g *genericArrayIterator) Next() interface{} {
	g.currentIndex++

	if g.currentIndex < 0 || g.currentIndex >= g.col.Len() {
		return nil
	}

	return g.col.Index(g.currentIndex)
}

// Resents the iterator position to the beginning.
func (g *genericArrayIterator) Reset() IIterator {
	g.currentIndex = 0
	return g
}

// Skips the following N items
func (g *genericArrayIterator) Skip(n int) IIterator {
	g.currentIndex += n
	return g
}
