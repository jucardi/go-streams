package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayIterator_Basic(t *testing.T) {
	iter := newArrayIterator[int]([]int{1, 2, 3})

	assert.Equal(t, 1, iter.Current())
	assert.True(t, iter.HasNext())
	assert.Equal(t, 2, iter.Next())
	assert.Equal(t, 2, iter.Current())
	assert.Equal(t, 3, iter.Next())
	// HasNext still true at last element; the standard pattern consumes it in the post-condition
	assert.True(t, iter.HasNext())
	// After moving past the last element, HasNext is false
	iter.Next()
	assert.False(t, iter.HasNext())
}

func TestArrayIterator_Empty(t *testing.T) {
	iter := newArrayIterator[int]()

	assert.False(t, iter.HasNext())
	var zero int
	assert.Equal(t, zero, iter.Current())
	assert.Equal(t, zero, iter.Next())
}

func TestArrayIterator_Skip(t *testing.T) {
	iter := newArrayIterator[int]([]int{10, 20, 30, 40, 50})
	iter.Skip(2)
	assert.Equal(t, 30, iter.Current())
}

func TestArrayIterator_ForEachRemaining(t *testing.T) {
	iter := newArrayIterator[int]([]int{1, 2, 3, 4, 5})
	var result []int
	iter.ForEachRemaining(func(x int) {
		result = append(result, x)
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestArrayIterator_ForEachRemaining_FromMiddle(t *testing.T) {
	iter := newArrayIterator[int]([]int{1, 2, 3, 4, 5})
	iter.Skip(2)
	var result []int
	iter.ForEachRemaining(func(x int) {
		result = append(result, x)
	})
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestArrayIterator_MoveNext(t *testing.T) {
	iter := newArrayIterator[int]([]int{1, 2})
	assert.True(t, iter.MoveNext())
	assert.Equal(t, 2, iter.Current())
	// MoveNext at last element: HasNext is still true (points at element), so it advances past
	assert.True(t, iter.MoveNext())
	// Now past the end
	assert.False(t, iter.MoveNext())
}

func TestArrayIterator_CurrentOutOfBounds(t *testing.T) {
	iter := newArrayIterator[int]([]int{1})
	iter.Skip(5)
	var zero int
	assert.Equal(t, zero, iter.Current())
}

func TestCollectionIterator_Basic(t *testing.T) {
	col := NewList[string]([]string{"a", "b", "c"})
	iter := newCollectionIterator[string](col)

	assert.Equal(t, "a", iter.Current())
	assert.Equal(t, "b", iter.Next())
	assert.Equal(t, "c", iter.Next())
	assert.True(t, iter.HasNext())
	iter.Next()
	assert.False(t, iter.HasNext())
}

func TestCollectionIterator_Empty(t *testing.T) {
	iter := newCollectionIterator[int]()
	assert.False(t, iter.HasNext())
}

func TestCollectionIterator_CurrentOutOfBounds(t *testing.T) {
	col := NewList[int]([]int{1})
	iter := newCollectionIterator[int](col)
	iter.Skip(5)
	var zero int
	assert.Equal(t, zero, iter.Current())
}

func TestIndexBasedIterator_StopAt(t *testing.T) {
	iter := newArrayIterator[int]([]int{1, 2, 3, 4, 5})
	posIter := iter.(IIteratorWithPos[int])
	posIter.StopAt(2)
	assert.Equal(t, 0, posIter.Pos())

	var result []int
	iter.ForEachRemaining(func(x int) {
		result = append(result, x)
	})
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestNewIterator_FromSlice(t *testing.T) {
	iter := NewIterator[int]([]int{1, 2, 3})
	assert.Equal(t, 1, iter.Current())
	assert.True(t, iter.HasNext())
}

func TestNewIterator_FromList(t *testing.T) {
	col := NewList[int]([]int{10, 20})
	iter := NewIterator[int](col)
	assert.Equal(t, 10, iter.Current())
}

func TestNewIterator_InvalidSource_Panics(t *testing.T) {
	assert.Panics(t, func() {
		NewIterator[int]("invalid")
	})
}
