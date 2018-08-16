package streams

import "reflect"

type arrayCollection struct {
	v           reflect.Value
	elementType reflect.Type
}

func (g *arrayCollection) Len() int {
	return g.v.Len()
}

func (g *arrayCollection) Get(index interface{}) interface{} {
	i, ok := index.(int)
	if !ok {
		panic("key must be `int` for array collections")
	}
	if i >= g.Len() {
		return nil
	}

	return g.v.Index(i).Interface()
}

func (g *arrayCollection) Add(item interface{}) error {
	if reflect.PtrTo(reflect.TypeOf(item)).AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	g.v = reflect.Append(g.v, reflect.ValueOf(item))
	return nil
}

// AddAll appends another iterable into this ICollection instance.
func (g *arrayCollection) AddAll(slice IIterable) error {
	if !slice.ElementType().AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	g.v = reflect.AppendSlice(g.v, reflect.ValueOf(slice.ToArray()))
	return nil
}

// ElementType returns the type of the elements in the iterable
func (g *arrayCollection) ElementType() reflect.Type {
	return g.elementType
}

func (g *arrayCollection) Iterator() IIterator {
	return newArrayIterator(g, g.elementType)
}

func (g *arrayCollection) ToArray() interface{} {
	return g.v.Interface()
}
