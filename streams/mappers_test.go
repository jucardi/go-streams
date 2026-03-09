package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappers_IntToString(t *testing.T) {
	fn := Mappers().IntToString()
	assert.Equal(t, "42", fn(42))
	assert.Equal(t, "0", fn(0))
	assert.Equal(t, "-1", fn(-1))
}

func TestMappers_StringToInt(t *testing.T) {
	fn := Mappers().StringToInt()
	assert.Equal(t, 42, fn("42"))
	assert.Equal(t, 0, fn("0"))
	assert.Equal(t, 0, fn("not_a_number"))
}

func TestMappers_StringToInt_WithErrorHandler(t *testing.T) {
	var capturedErr error
	var capturedStr string

	fn := Mappers().StringToInt(func(s string, err error) {
		capturedStr = s
		capturedErr = err
	})

	result := fn("invalid")
	assert.Equal(t, 0, result)
	assert.Equal(t, "invalid", capturedStr)
	assert.NotNil(t, capturedErr)
}

func TestMappers_StringToInt_WithNilErrorHandler(t *testing.T) {
	fn := Mappers().StringToInt(nil)
	result := fn("invalid")
	assert.Equal(t, 0, result)
}

func TestMap_FromArray(t *testing.T) {
	result := Map[int, string]([]int{1, 2, 3}, Mappers().IntToString())
	assert.Equal(t, []string{"1", "2", "3"}, result.ToArray())
}

func TestMap_FromIterable(t *testing.T) {
	col := NewList[int]([]int{10, 20})
	result := Map[int, string](col, Mappers().IntToString())
	assert.Equal(t, []string{"10", "20"}, result.ToArray())
}

func TestMap_FromIterator(t *testing.T) {
	iter := newArrayIterator[int]([]int{5, 6})
	result := Map[int, string](iter, Mappers().IntToString())
	assert.Equal(t, []string{"5", "6"}, result.ToArray())
}

func TestMap_FromStream(t *testing.T) {
	stream := FromArray[int]([]int{1, 2, 3, 4}).
		Filter(func(x int) bool { return x > 2 })
	result := Map[int, string](stream, Mappers().IntToString())
	assert.Equal(t, []string{"3", "4"}, result.ToArray())
}

func TestMap_InvalidSource_Panics(t *testing.T) {
	assert.Panics(t, func() {
		Map[int, string]("invalid", Mappers().IntToString())
	})
}

func TestMapNonComparable_FromArray(t *testing.T) {
	result := MapNonComparable[int, string]([]int{1, 2}, Mappers().IntToString())
	assert.Equal(t, []string{"1", "2"}, result)
}

func TestMapNonComparable_FromIterable(t *testing.T) {
	col := NewList[int]([]int{3, 4})
	result := MapNonComparable[int, string](col, Mappers().IntToString())
	assert.Equal(t, []string{"3", "4"}, result)
}

func TestMapNonComparable_FromIterator(t *testing.T) {
	iter := newArrayIterator[int]([]int{7, 8})
	result := MapNonComparable[int, string](iter, Mappers().IntToString())
	assert.Equal(t, []string{"7", "8"}, result)
}

func TestMapNonComparable_InvalidSource_Panics(t *testing.T) {
	assert.Panics(t, func() {
		MapNonComparable[int, string]("invalid", Mappers().IntToString())
	})
}
