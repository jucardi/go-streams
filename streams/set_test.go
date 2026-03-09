package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet_Add(t *testing.T) {
	s := NewSet[int]()
	assert.True(t, s.Add(1, 2, 3))
	assert.Equal(t, 3, s.Len())

	assert.False(t, s.Add(1))
	assert.Equal(t, 3, s.Len())
}

func TestSet_AddFromIterator(t *testing.T) {
	s := NewSet[int]()
	iter := newArrayIterator[int]([]int{1, 2, 3})
	assert.True(t, s.AddFromIterator(iter))
	assert.Equal(t, 3, s.Len())
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)

	assert.True(t, s.Remove(2))
	assert.Equal(t, 2, s.Len())
	assert.False(t, s.Contains(2))

	assert.False(t, s.Remove(99))
}

func TestSet_RemoveFromIterator(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3, 4, 5)

	iter := newArrayIterator[int]([]int{2, 4})
	assert.True(t, s.RemoveFromIterator(iter))
	assert.Equal(t, 3, s.Len())
	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))
}

func TestSet_RemoveIf(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3, 4, 5)

	removed := s.RemoveIf(func(x int) bool { return x > 3 })
	assert.True(t, removed)
	assert.Equal(t, 3, s.Len())

	notRemoved := s.RemoveIf(func(x int) bool { return x > 100 })
	assert.False(t, notRemoved)
}

func TestSet_Contains(t *testing.T) {
	s := NewSet[string]()
	s.Add("a", "b", "c")

	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("a", "b"))
	assert.False(t, s.Contains("d"))
	assert.False(t, s.Contains("a", "d"))
}

func TestSet_ContainsFromIterator(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)

	iter1 := newArrayIterator[int]([]int{1, 2})
	assert.True(t, s.ContainsFromIterator(iter1))

	iter2 := newArrayIterator[int]([]int{1, 99})
	assert.False(t, s.ContainsFromIterator(iter2))
}

func TestSet_Len(t *testing.T) {
	s := NewSet[int]()
	assert.Equal(t, 0, s.Len())

	s.Add(1, 2)
	assert.Equal(t, 2, s.Len())
}

func TestSet_Clear(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	s.Clear()
	assert.Equal(t, 0, s.Len())
	assert.True(t, s.IsEmpty())
}

func TestSet_ToArray(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	arr := s.ToArray()
	assert.Len(t, arr, 3)
	assert.Contains(t, arr, 1)
	assert.Contains(t, arr, 2)
	assert.Contains(t, arr, 3)
}

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[int]()
	assert.True(t, s.IsEmpty())
	s.Add(1)
	assert.False(t, s.IsEmpty())
}

func TestSet_Iterator(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)

	iter := s.Iterator()
	count := 0
	for iter.HasNext() {
		iter.Next()
		count++
	}
	assert.Equal(t, 3, count)
}

func TestSet_ForEach(t *testing.T) {
	s := NewSet[int]()
	s.Add(10, 20, 30)

	var sum int
	s.ForEach(func(x int) { sum += x })
	assert.Equal(t, 60, sum)
}
