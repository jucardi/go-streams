package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapCollection_Set_And_Get(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)

	val, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = m.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	_, ok = m.Get("c")
	assert.False(t, ok)
}

func TestMapCollection_Set_OverwriteExisting(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("a", 99)

	val, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 99, val)
	assert.Equal(t, 1, m.Len())
}

func TestMapCollection_ContainsKey(t *testing.T) {
	m := NewMap[string, int](map[string]int{"x": 10, "y": 20})
	assert.True(t, m.ContainsKey("x"))
	assert.True(t, m.ContainsKey("y"))
	assert.False(t, m.ContainsKey("z"))
}

func TestMapCollection_Delete(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1, "b": 2, "c": 3})
	assert.True(t, m.Delete("b"))
	assert.Equal(t, 2, m.Len())
	assert.False(t, m.ContainsKey("b"))

	assert.False(t, m.Delete("nonexistent"))
}

func TestMapCollection_Keys(t *testing.T) {
	m := NewMap[string, int](map[string]int{"x": 1, "y": 2})
	keys := m.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "x")
	assert.Contains(t, keys, "y")
}

func TestMapCollection_ToMap(t *testing.T) {
	original := map[string]int{"a": 1, "b": 2}
	m := NewMap[string, int](original)
	result := m.ToMap()
	assert.Equal(t, original, result)
}

func TestMapCollection_Add(t *testing.T) {
	m := NewMap[string, int]()
	result := m.Add(&KeyValuePair[string, int]{Key: "k1", Value: 10})
	assert.True(t, result)
	val, ok := m.Get("k1")
	assert.True(t, ok)
	assert.Equal(t, 10, val)

	assert.False(t, m.Add())
}

func TestMapCollection_RemoveAt(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1, "b": 2, "c": 3})

	assert.False(t, m.RemoveAt(-1))
	assert.False(t, m.RemoveAt(100))

	assert.True(t, m.RemoveAt(0))
	assert.Equal(t, 2, m.Len())
}

func TestMapCollection_Len(t *testing.T) {
	m := NewMap[string, int]()
	assert.Equal(t, 0, m.Len())
	m.Set("a", 1)
	assert.Equal(t, 1, m.Len())
}

func TestMapCollection_Clear(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1, "b": 2})
	m.Clear()
	assert.Equal(t, 0, m.Len())
}

func TestMapCollection_ToArray(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1})
	arr := m.ToArray()
	assert.Len(t, arr, 1)
	assert.Equal(t, "a", arr[0].Key)
	assert.Equal(t, 1, arr[0].Value)
}

func TestMapCollection_Index(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1})

	val, exists := m.Index(0)
	assert.True(t, exists)
	assert.NotNil(t, val)

	_, exists = m.Index(-1)
	assert.False(t, exists)

	_, exists = m.Index(10)
	assert.False(t, exists)
}

func TestMapCollection_NewMapEmpty(t *testing.T) {
	m := NewMap[string, string]()
	assert.Equal(t, 0, m.Len())
	m.Set("hello", "world")
	assert.Equal(t, 1, m.Len())
}

func TestMapCollection_Delete_NonExistent(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1})
	assert.False(t, m.Delete("nonexistent"))
	assert.Equal(t, 1, m.Len())
}

func TestMapCollection_Index_KeyValueMatch(t *testing.T) {
	m := NewMap[string, int](map[string]int{"a": 1, "b": 2})
	// Delete a key so the keys slice and map go out of sync, then Index
	m.Delete("a")
	val, exists := m.Index(0)
	assert.True(t, exists)
	assert.Equal(t, "b", val.Key)
	assert.Equal(t, 2, val.Value)
}
