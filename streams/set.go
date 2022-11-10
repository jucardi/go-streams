package streams

import (
	"sync"
)

var (
	// To ensure *set implements ISet on build
	_ ISet[string]        = (*set[string])(nil)
	_ ICollection[string] = (*set[string])(nil)
)

func NewSet[T comparable]() ISet[T] {
	return &set[T]{m: map[T]struct{}{}}
}

type set[T comparable] struct {
	m  map[T]struct{}
	mx sync.RWMutex
}

func (c *set[T]) Iterator() IIterator[T] {
	return newArrayIterator[T](c.ToArray())
}

func (c *set[T]) ForEach(f IterFunc[T]) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	for item := range c.m {
		f(item)
	}
}

func (c *set[T]) Add(items ...T) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	l := len(c.m)

	for _, item := range items {
		c.m[item] = struct{}{}
	}

	return len(c.m) > l
}

func (c *set[T]) AddFromIterator(iterator IIterator[T]) bool {
	ret := false
	iterator.ForEachRemaining(func(item T) {
		ret = ret || c.Add(item)
	})
	return ret
}

func (c *set[T]) Remove(items ...T) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	l := len(c.m)
	for _, item := range items {
		delete(c.m, item)
	}

	return len(c.m) < l
}

func (c *set[T]) RemoveFromIterator(iterator IIterator[T]) bool {
	ret := false
	iterator.ForEachRemaining(func(item T) {
		ret = ret || c.Remove(item)
	})
	return ret
}

func (c *set[T]) RemoveIf(condition ConditionalFunc[T], _ ...bool) bool {
	var toRemove []T
	c.ForEach(func(item T) {
		if condition(item) {
			toRemove = append(toRemove, item)
		}
	})
	if len(toRemove) == 0 {
		return false
	}

	return c.Remove(toRemove...)
}

func (c *set[T]) Contains(item ...T) bool {
	c.mx.RLock()
	defer c.mx.RUnlock()

	for _, x := range item {
		if _, ok := c.m[x]; !ok {
			return false
		}
	}
	return true
}

func (c *set[T]) ContainsFromIterator(iterator IIterator[T]) bool {
	c.mx.RLock()
	defer c.mx.RUnlock()

	for val := iterator.Current(); iterator.HasNext(); val = iterator.Next() {
		if _, ok := c.m[val]; !ok {
			return false
		}
	}
	return true
}

func (c *set[T]) Len() int {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return len(c.m)
}

func (c *set[T]) Clear() {
	c.m = map[T]struct{}{}
}

func (c *set[T]) ToArray() (ret []T) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	for item := range c.m {
		ret = append(ret, item)
	}
	return
}

func (c *set[T]) IsEmpty() bool {
	return len(c.m) == 0
}
