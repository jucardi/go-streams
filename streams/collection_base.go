package streams

/*
This file contains base implementations of collections, so only members in the abstract collection interfaces would need
to be implemented.

Base implementations provided:
  - CollectionBase:           Has a default implementation of all IList functions except for Index, Add, RemoveAt,
                              Len, Clear and ToArray.

  - CollectionBaseNoIterator: Similar to CollectionBase but also requires Iterator() to be implemented.

Usage:


    Step 1:  Simply add either struct as the base of the new struct implementing IList.

             Eg:   type customCollection[T comparable] struct {
                       *CollectionBase[T]
                       . . .
                   }

    Step 2:  Implement a constructor that properly sets the references between the base and the
             implementation using SetAbstract

             func NewCustomCollection[T comparable]() IList[T] {
                 // create a new instance of the collection implementation
                 ret := &customCollection[T]{}

                 // create an instance of the base to be used
                 base := &CollectionBase[T]{}

                 // set the custom collection implementation instance as the abstract members of the base.
                 base.SetAbstract(ret)

                 // set the base instance as the base for the custom implementation
                 ret.CollectionBase = base

                 // done
                 return ret
             }
*/

// IAbstractCollection defines the members that need to be defined in a collection implementation if the implementation
// is using the type CollectionBase as the base struct. This base uses the default implementation of a collection
// iterator and does not require Iterator() to be implemented
type IAbstractCollection[T comparable] interface {
	Index(index int) (val T, exists bool)
	Add(item ...T) bool
	RemoveAt(index int, keepOrder ...bool) bool
	Len() int
	Clear()
	ToArray() []T
}

// IAbstractCollectionWithIterator defines the members that need to be defined in a collection implementation if the
// implementation is using the type CollectionBaseNoIterator as the base struct. This base assumes the collection that
// is being implemented includes the Iterator() function
type IAbstractCollectionWithIterator[T comparable] interface {
	IAbstractCollection[T]
	Iterator() IIterator[T]
}

type iAbstractCollectionIterator[T comparable] interface {
	Iterator() IIterator[T]
}

type CollectionBaseNoIterator[T comparable] struct {
	IAbstractCollection[T]
	iAbstractCollectionIterator[T]
	set ISet[T]
}

func (c *CollectionBaseNoIterator[T]) modified() {
	c.set = nil
}

func (c *CollectionBaseNoIterator[T]) ForEach(f IterFunc[T]) {
	c.Iterator().ForEachRemaining(f)
}

func (c *CollectionBaseNoIterator[T]) AddFromIterator(iterator IIterator[T]) (ret bool) {
	iterator.ForEachRemaining(func(item T) {
		ret = ret || c.Add(item)
	})
	return
}

func (c *CollectionBaseNoIterator[T]) Remove(items ...T) bool {
	return c.RemoveIf(func(x T) bool {
		for _, y := range items {
			if x == y {
				return true
			}
		}
		return false
	})
}

func (c *CollectionBaseNoIterator[T]) RemoveFromIterator(iterator IIterator[T]) (ret bool) {
	iterator.ForEachRemaining(func(item T) {
		ret = ret || c.Remove(item)
	})
	return ret
}

func (c *CollectionBaseNoIterator[T]) RemoveIf(condition ConditionalFunc[T], keepOrder ...bool) bool {
	removed := 0
	count := c.Len()
	for i := 0; i < count-removed; i++ {
		val, _ := c.Index(i)
		if condition(val) && c.RemoveAt(i+removed, keepOrder...) {
			removed++
		}
	}
	return removed > 0
}

func (c *CollectionBaseNoIterator[T]) Contains(item ...T) bool {
	return c.Distinct().Contains(item...)
}

func (c *CollectionBaseNoIterator[T]) ContainsFromIterator(iterator IIterator[T]) bool {
	ret := true
	iterator.ForEachRemaining(func(item T) {
		ret = ret && c.Contains(item)
	})
	return ret
}

func (c *CollectionBaseNoIterator[T]) IsEmpty() bool {
	return c.Len() == 0
}

func (c *CollectionBaseNoIterator[T]) Distinct() ISet[T] {
	if c.set != nil {
		return c.set
	}
	set := NewSet[T]()
	set.Add(c.ToArray()...)
	c.set = set
	return set
}

func (c *CollectionBaseNoIterator[T]) Stream() IStream[T] {
	return FromCollection[T](c)
}

func (c *CollectionBaseNoIterator[T]) SetAbstract(col IAbstractCollectionWithIterator[T]) {
	c.IAbstractCollection = col
	c.iAbstractCollectionIterator = col
}

type CollectionBase[T comparable] struct {
	CollectionBaseNoIterator[T]
}

func (c *CollectionBase[T]) Iterator() IIterator[T] {
	return newCollectionIterator[T](c)
}

func (c *CollectionBase[T]) SetAbstract(col IAbstractCollection[T]) {
	c.CollectionBaseNoIterator.IAbstractCollection = col
	c.CollectionBaseNoIterator.iAbstractCollectionIterator = c
}
