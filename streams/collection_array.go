package streams

var (
	// To ensure *arrayCollection implements IList on build
	_ IList[string] = (*arrayCollection[string])(nil)
)

type arrayCollection[T comparable] struct {
	*CollectionBase[T]
	arr []T
}

func (c *arrayCollection[T]) Index(index int) (ret T, exists bool) {
	if index < 0 || index >= c.Len() {
		return
	}

	return c.arr[index], true
}

func (c *arrayCollection[T]) Add(item ...T) bool {
	c.arr = append(c.arr, item...)
	c.modified()
	return true
}

func (c *arrayCollection[T]) RemoveAt(index int, keepOrder ...bool) bool {
	if index < 0 || index >= c.Len() {
		return false
	}

	if len(keepOrder) > 0 && keepOrder[0] {
		c.removeKeepOrder(index)
	} else {
		c.removeFast(index)
	}
	c.modified()
	return true
}

func (c *arrayCollection[T]) Len() int {
	return len(c.arr)
}

func (c *arrayCollection[T]) Clear() {
	c.arr = nil
}

func (c *arrayCollection[T]) ToArray() []T {
	return c.arr
}

func (c *arrayCollection[T]) IsEmpty() bool {
	return c.Len() == 0
}

func (c *arrayCollection[T]) Stream() IStream[T] {
	return FromCollection[T](c)
}

// removeFast swaps the element to remove with the last element, then shrinks the array size by one. The order of the elements is not ensured with this method
func (c *arrayCollection[T]) removeFast(index int) (ret T) {
	c.arr = append(c.arr[0:index], c.arr[index:]...)
	last := c.arr[c.Len()-1]
	ret = c.arr[index]
	c.arr[index] = last
	c.arr = c.arr[:len(c.arr)-1]
	return
}

// removeKeepOrder creates a slice from the beginning of the slice up to the element before the provided index, then it creates another slice from the index+1 element to the end.
// This function guarantees the original order of the elements but it can be a costly operation since the elements in the original slice need to be shifted one position below.
func (c *arrayCollection[T]) removeKeepOrder(index int) (ret T) {
	ret = c.arr[index]
	c.arr = append(c.arr[0:index], c.arr[index:]...)
	return
}
