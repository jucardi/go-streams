package streams

import (
	"errors"
	"reflect"
)

type arrayCollection struct {
	v           reflect.Value
	elementType reflect.Type
}

func (g *arrayCollection) Len() int {
	if g.v.IsValid() {
		return g.v.Len()
	}
	return 0
}

func (g *arrayCollection) Index(index int) interface{} {
	if index < 0 || index >= g.Len() {
		return nil
	}

	return g.v.Index(index).Interface()
}

func (g *arrayCollection) Remove(index int, keepOrder ...bool) interface{} {
	if len(keepOrder) > 0 && keepOrder[0] {
		return g.removeKeepOrder(index)
	}
	return g.removeFast(index)
}

func (g *arrayCollection) Add(item interface{}) error {
	if item == nil {
		return errors.New("unable to add nil value")
	}
	if g.elementType == nil {
		g.elementType = reflect.TypeOf(item)
		g.v = reflect.MakeSlice(reflect.SliceOf(g.elementType), 0, 0)
	}
	if reflect.PtrTo(reflect.TypeOf(item)).AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	g.v = reflect.Append(g.v, reflect.ValueOf(item))
	return nil
}

func (g *arrayCollection) AddAll(slice IIterable) error {
	if !slice.ElementType().AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	g.v = reflect.AppendSlice(g.v, reflect.ValueOf(slice.ToArray()))
	return nil
}

func (g *arrayCollection) ElementType() reflect.Type {
	return g.elementType
}

func (g *arrayCollection) Iterator() IIterator {
	return newCollectionIterator(g)
}

func (g *arrayCollection) ToArray(defaultArray ...interface{}) interface{} {
	if (!g.v.IsValid() || (g.v.IsValid() && g.v.IsNil())) && len(defaultArray) > 0 {
		return defaultArray[0]
	}
	return g.v.Interface()
}

// removeFast swaps the element to remove with the last element, then shrinks the array size by one. The order of the elements is not ensured with this method
func (g *arrayCollection) removeFast(index int) interface{} {
	if index < 0 || index >= g.Len() {
		return nil
	}

	last := g.v.Index(g.Len() - 1)
	toRemove := g.v.Index(index)
	ret := toRemove.Interface()
	toRemove.Set(last)
	g.v = g.v.Slice(0, g.Len()-1)
	return ret
}

// removeKeepOrder creates a slice from the beginning of the slice up to the element before the provided index, then it creates another slice from the index+1 element to the end.
// This function guarantees the original order of the elements but it can be a costly operation since the elements in the original slice need to be shifted one position below.
func (g *arrayCollection) removeKeepOrder(index int) interface{} {
	if index < 0 || index >= g.Len() {
		return nil
	}

	ret := g.Index(index)
	firstHalf := g.v.Slice(0, index)
	secondHalf := g.v.Slice(index, g.Len())
	g.v = reflect.Append(firstHalf, secondHalf)
	return ret
}
