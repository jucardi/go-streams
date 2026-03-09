package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort_IntAscending(t *testing.T) {
	arr := []int{5, 3, 1, 4, 2}
	Sort(arr)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr)
}

func TestSort_IntDescending(t *testing.T) {
	arr := []int{5, 3, 1, 4, 2}
	Sort(arr, true)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, arr)
}

func TestSort_StringAscending(t *testing.T) {
	arr := []string{"banana", "apple", "cherry"}
	Sort(arr)
	assert.Equal(t, []string{"apple", "banana", "cherry"}, arr)
}

func TestSort_StringDescending(t *testing.T) {
	arr := []string{"banana", "apple", "cherry"}
	Sort(arr, true)
	assert.Equal(t, []string{"cherry", "banana", "apple"}, arr)
}

func TestSort_Float(t *testing.T) {
	arr := []float64{3.14, 1.0, 2.71}
	Sort(arr)
	assert.Equal(t, []float64{1.0, 2.71, 3.14}, arr)
}

func TestSort_Empty(t *testing.T) {
	arr := []int{}
	Sort(arr)
	assert.Empty(t, arr)
}

func TestComparableFn_Ascending(t *testing.T) {
	fn := ComparableFn[int]()
	assert.Equal(t, -1, fn(1, 2))
	assert.Equal(t, 0, fn(2, 2))
	assert.Equal(t, 1, fn(3, 2))
}

func TestComparableFn_Descending(t *testing.T) {
	fn := ComparableFn[int](true)
	assert.Equal(t, 1, fn(1, 2))
	assert.Equal(t, 0, fn(2, 2))
	assert.Equal(t, -1, fn(3, 2))
}
