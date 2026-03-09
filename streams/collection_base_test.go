package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionBase_AddFromIterator(t *testing.T) {
	col := NewList[int]([]int{1, 2})
	iter := newArrayIterator[int]([]int{3, 4, 5})

	result := col.AddFromIterator(iter)
	assert.True(t, result)
	assert.Equal(t, 5, col.Len())
	assert.Equal(t, []int{1, 2, 3, 4, 5}, col.ToArray())
}

func TestCollectionBase_Remove(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3, 4, 5})
	result := col.Remove(2, 4)
	assert.True(t, result)
	arr := col.ToArray()
	assert.NotContains(t, arr, 2)
	assert.NotContains(t, arr, 4)
}

func TestCollectionBase_Remove_NotFound(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	result := col.Remove(99)
	assert.False(t, result)
	assert.Equal(t, 3, col.Len())
}

func TestCollectionBase_RemoveFromIterator(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3, 4, 5})
	iter := newArrayIterator[int]([]int{2, 4})
	result := col.RemoveFromIterator(iter)
	assert.True(t, result)
	arr := col.ToArray()
	assert.NotContains(t, arr, 2)
	assert.NotContains(t, arr, 4)
}

func TestCollectionBase_RemoveIf(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3, 4, 5})
	result := col.RemoveIf(func(x int) bool { return x > 3 })
	assert.True(t, result)
	assert.NotContains(t, col.ToArray(), 4)
	assert.NotContains(t, col.ToArray(), 5)
}

func TestCollectionBase_RemoveIf_NoneMatch(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	result := col.RemoveIf(func(x int) bool { return x > 100 })
	assert.False(t, result)
}

func TestCollectionBase_Contains(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	assert.True(t, col.Contains(2))
	assert.True(t, col.Contains(1, 3))
	assert.False(t, col.Contains(99))
}

func TestCollectionBase_ContainsFromIterator(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3, 4, 5})

	iter1 := newArrayIterator[int]([]int{2, 4})
	assert.True(t, col.ContainsFromIterator(iter1))

	iter2 := newArrayIterator[int]([]int{2, 99})
	assert.False(t, col.ContainsFromIterator(iter2))
}

func TestCollectionBase_IsEmpty(t *testing.T) {
	col := NewList[int]()
	assert.True(t, col.IsEmpty())
	col.Add(1)
	assert.False(t, col.IsEmpty())
}

func TestCollectionBase_Distinct(t *testing.T) {
	col := NewList[int]([]int{1, 2, 2, 3, 3, 3})
	set := col.Distinct()
	assert.Equal(t, 3, set.Len())
	assert.True(t, set.Contains(1, 2, 3))

	set2 := col.Distinct()
	assert.Equal(t, set, set2)
}

func TestCollectionBase_Distinct_CacheClearedOnModify(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	set1 := col.Distinct()
	assert.Equal(t, 3, set1.Len())

	col.Add(4)
	set2 := col.Distinct()
	assert.Equal(t, 4, set2.Len())
}

func TestCollectionBase_Stream(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	stream := col.Stream()
	assert.Equal(t, 3, stream.Count())
}

func TestCollectionBase_ForEach(t *testing.T) {
	col := NewList[int]([]int{10, 20, 30})
	var sum int
	col.ForEach(func(x int) { sum += x })
	assert.Equal(t, 60, sum)
}

func TestCollectionBase_RemoveIf_KeepOrder(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3, 4, 5})
	result := col.RemoveIf(func(x int) bool { return x%2 == 0 }, true)
	assert.True(t, result)
	assert.Equal(t, []int{1, 3, 5}, col.ToArray())
}
