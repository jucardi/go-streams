package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrom_Array(t *testing.T) {
	stream := From[int]([]int{1, 2, 3})
	assert.Equal(t, 3, stream.Count())
}

func TestFrom_Collection(t *testing.T) {
	col := NewList[int]([]int{1, 2, 3})
	stream := From[int](col)
	assert.Equal(t, 3, stream.Count())
}

func TestFrom_InvalidSource_Panics(t *testing.T) {
	assert.Panics(t, func() {
		From[int]("invalid")
	})
}

func TestFromMap_Map(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	stream := FromMap[string, int](m)
	assert.Equal(t, 2, stream.Count())
}

func TestFromMap_KeyValuePairArray(t *testing.T) {
	pairs := []*KeyValuePair[string, int]{
		{Key: "x", Value: 10},
		{Key: "y", Value: 20},
	}
	stream := FromMap[string, int](pairs)
	assert.Equal(t, 2, stream.Count())
}

func TestFromMap_Collection(t *testing.T) {
	col := NewMap[string, int](map[string]int{"a": 1})
	stream := FromMap[string, int](col)
	assert.Equal(t, 1, stream.Count())
}

func TestFromMap_InvalidSource_Panics(t *testing.T) {
	assert.Panics(t, func() {
		FromMap[string, int]("invalid")
	})
}

func TestFromArray(t *testing.T) {
	stream := FromArray[string]([]string{"hello", "world"})
	assert.Equal(t, 2, stream.Count())
	assert.Equal(t, "hello", stream.First())
}

func TestFromCollection(t *testing.T) {
	col := NewList[int]([]int{5, 10, 15})
	stream := FromCollection[int](col)
	assert.Equal(t, 3, stream.Count())
}

func TestNewList_Empty(t *testing.T) {
	list := NewList[int]()
	assert.Equal(t, 0, list.Len())
	assert.True(t, list.IsEmpty())
}

func TestNewList_WithArray(t *testing.T) {
	list := NewList[int]([]int{1, 2, 3})
	assert.Equal(t, 3, list.Len())
}

func TestNewMap_WithValues(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1, "b": 2})
	assert.Equal(t, 2, m.Len())
}

func TestMapToPtr_ConvertsCorrectly(t *testing.T) {
	type item struct {
		Name string
	}
	arr := []item{{Name: "a"}, {Name: "b"}}
	result := MapToPtr[item](arr)
	assert.Len(t, result, 2)
	assert.Equal(t, "a", result[0].Name)
	assert.Equal(t, "b", result[1].Name)
}
