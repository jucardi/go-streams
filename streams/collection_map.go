package streams

import "reflect"

type mapCollection struct {
	keySet ICollection
	v      reflect.Value
}

var keyValuePairType = reflect.TypeOf((*KeyValuePair)(nil))

func (g *mapCollection) init() *mapCollection {
	g.updateKeys()
	return g
}

func (g *mapCollection) updateKeys() {
	// Ideally, a collection implementation from a Map would have been defined that knows how to iterate over a K,V set to avoid a full map iteration. However, there is not way to
	// iterate over a K,V set of a map through reflection, instead the only thing available is the function `MapKeys`.
	//
	// Please note that given the nature of `MapKeys`, a full iteration over the map will always happen when creating a new map collection. Also a full iteration will happen if
	// a value was added to the original map instead of using the `Add` and `AddAll` functions provided in this collection.
	//

	if g.keySet != nil && g.keySet.Len() == g.v.Len() {
		// Validates if the keySet is out of sync from the source map. Should never happen after creating the `mapCollection` for the first time if adding or removing items are done
		// through the `mapCollection` instance.
		return
	}

	keySet, _ := NewCollectionFromArray(g.v.MapKeys())
	g.keySet = keySet
}

func (g *mapCollection) Len() int {
	return g.v.Len()
}

func (g *mapCollection) Index(index int) interface{} {
	g.updateKeys()
	if index >= g.keySet.Len() {
		return nil
	}

	key := g.keySet.Index(index).(reflect.Value)
	return &KeyValuePair{
		Key:   key.Interface(),
		Value: g.v.MapIndex(key).Interface(),
	}
}

func (g *mapCollection) Remove(index int, keepOrder ...bool) interface{} {
	// For a HashMap, keepOrder has no effect since the map balances itself after an item is removed.
	g.updateKeys()
	key := g.keySet.Remove(index)
	ret := g.v.MapIndex(reflect.ValueOf(key)).Interface()
	g.v.SetMapIndex(key.(reflect.Value), reflect.Value{})
	return ret
}

func (g *mapCollection) Get(key interface{}) interface{} {
	g.updateKeys()
	return g.v.MapIndex(reflect.ValueOf(key)).Interface()
}

func (g *mapCollection) Add(item interface{}) error {
	if reflect.PtrTo(reflect.TypeOf(item)).AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	pair := item.(*KeyValuePair)
	keyVal := reflect.ValueOf(pair.Key)
	valueVal := reflect.ValueOf(pair.Value)
	_ = g.keySet.Add(keyVal)
	g.v.SetMapIndex(keyVal, valueVal)
	return nil
}

// AddAll appends another iterable into this ICollection instance.
func (g *mapCollection) AddAll(slice IIterable) error {
	if !slice.ElementType().AssignableTo(g.ElementType()) {
		return ErrorWrongType
	}

	iter := slice.Iterator()
	for x := iter.Current(); iter.HasNext(); x = iter.Next() {
		if err := g.Add(x); err != nil {
			return err
		}
	}
	return nil
}

// ElementType returns the type of the elements in the iterable
func (g *mapCollection) ElementType() reflect.Type {
	return keyValuePairType
}

func (g *mapCollection) Iterator() IIterator {
	return newCollectionIterator(g)
}

func (g *mapCollection) ToArray(_ ...interface{}) interface{} {
	var array []*KeyValuePair
	for _, key := range g.v.MapKeys() {
		array = append(array, &KeyValuePair{
			Key:   key.Interface(),
			Value: g.v.MapIndex(key).Interface(),
		})
	}
	return array
}
