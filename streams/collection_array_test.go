package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayCollection_Add(t *testing.T) {
	col := NewList[int]()
	assert.True(t, col.Add(1, 2, 3))
	assert.Equal(t, 3, col.Len())
	assert.Equal(t, []int{1, 2, 3}, col.ToArray())
}

func TestArrayCollection_Index(t *testing.T) {
	col := NewList[int]([]int{10, 20, 30})

	val, exists := col.(*arrayCollection[int]).Index(0)
	assert.True(t, exists)
	assert.Equal(t, 10, val)

	val, exists = col.(*arrayCollection[int]).Index(2)
	assert.True(t, exists)
	assert.Equal(t, 30, val)

	_, exists = col.(*arrayCollection[int]).Index(-1)
	assert.False(t, exists)

	_, exists = col.(*arrayCollection[int]).Index(3)
	assert.False(t, exists)
}

func TestArrayCollection_RemoveAt_KeepOrder(t *testing.T) {
	col := NewList[int]([]int{10, 20, 30, 40, 50})

	result := col.RemoveAt(2, true)
	assert.True(t, result)
	assert.Equal(t, 4, col.Len())
	assert.Equal(t, []int{10, 20, 40, 50}, col.ToArray())
}

func TestArrayCollection_RemoveAt_Fast(t *testing.T) {
	col := NewList[int]([]int{10, 20, 30, 40, 50})

	result := col.RemoveAt(1)
	assert.True(t, result)
	assert.Equal(t, 4, col.Len())
	assert.Contains(t, col.ToArray(), 10)
	assert.Contains(t, col.ToArray(), 30)
	assert.Contains(t, col.ToArray(), 40)
	assert.Contains(t, col.ToArray(), 50)
	assert.NotContains(t, col.ToArray(), 20)
}

func TestArrayCollection_RemoveAt_OutOfBounds(t *testing.T) {
	col := NewList[int]([]int{10, 20, 30})

	assert.False(t, col.RemoveAt(-1))
	assert.False(t, col.RemoveAt(3))
	assert.False(t, col.RemoveAt(100))
	assert.Equal(t, 3, col.Len())
}

func TestArrayCollection_Clear(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	assert.Equal(t, 3, col.Len())
	col.Clear()
	assert.Equal(t, 0, col.Len())
	assert.True(t, col.IsEmpty())
}

func TestArrayCollection_IsEmpty(t *testing.T) {
	col := NewList[int]()
	assert.True(t, col.IsEmpty())

	col.Add(1)
	assert.False(t, col.IsEmpty())
}

func TestArrayCollection_Stream(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	stream := col.Stream()
	assert.Equal(t, 3, stream.Count())
}

func TestArrayCollection_ToArray(t *testing.T) {
	col := NewList[string]([]string{"a", "b"})
	arr := col.ToArray()
	assert.Equal(t, []string{"a", "b"}, arr)
}

func TestArrayCollection_RemoveFast_SingleElement(t *testing.T) {
	col := NewList[int]([]int{42})
	assert.True(t, col.RemoveAt(0))
	assert.Equal(t, 0, col.Len())
}

func TestArrayCollection_RemoveKeepOrder_First(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	assert.True(t, col.RemoveAt(0, true))
	assert.Equal(t, []int{2, 3}, col.ToArray())
}

func TestArrayCollection_RemoveKeepOrder_Last(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	assert.True(t, col.RemoveAt(2, true))
	assert.Equal(t, []int{1, 2}, col.ToArray())
}
