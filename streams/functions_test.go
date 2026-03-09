package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce_Sum(t *testing.T) {
	result := Reduce[int]([]int{1, 2, 3, 4}, 0, func(acc, x int) int { return acc + x })
	assert.Equal(t, 10, result)
}

func TestReduce_Concat(t *testing.T) {
	result := Reduce[string]([]string{"a", "b", "c"}, "", func(acc, x string) string { return acc + x })
	assert.Equal(t, "abc", result)
}

func TestReduce_Empty(t *testing.T) {
	result := Reduce[int]([]int{}, 42, func(acc, x int) int { return acc + x })
	assert.Equal(t, 42, result)
}

func TestReduce_FromStream(t *testing.T) {
	stream := FromArray[int]([]int{1, 2, 3}).Filter(func(x int) bool { return x > 1 })
	result := Reduce[int](stream, 0, func(acc, x int) int { return acc + x })
	assert.Equal(t, 5, result)
}

func TestReduce_FromIterable(t *testing.T) {
	col := NewList[int]([]int{10, 20, 30})
	result := Reduce[int](col, 0, func(acc, x int) int { return acc + x })
	assert.Equal(t, 60, result)
}

func TestReduceAny_SumToString(t *testing.T) {
	result := ReduceAny[int, string]([]int{1, 2, 3}, "", func(acc string, x int) string {
		if acc != "" {
			acc += ","
		}
		return acc + Mappers().IntToString()(x)
	})
	assert.Equal(t, "1,2,3", result)
}

func TestFlatMap(t *testing.T) {
	result := FlatMap[string, string]([]string{"hello", "world"}, func(s string) []string {
		var chars []string
		for _, c := range s {
			chars = append(chars, string(c))
		}
		return chars
	})
	assert.Equal(t, 10, result.Len())
	assert.Equal(t, "h", result.ToArray()[0])
}

func TestFlatMap_Empty(t *testing.T) {
	result := FlatMap[int, int]([]int{}, func(x int) []int { return []int{x, x} })
	assert.Equal(t, 0, result.Len())
}

func TestConcat(t *testing.T) {
	result := Concat[int]([]int{1, 2}, []int{3, 4}, []int{5})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result.ToArray())
}

func TestConcat_Empty(t *testing.T) {
	result := Concat[int]()
	assert.Equal(t, 0, result.Count())
}

func TestConcat_WithStreams(t *testing.T) {
	s1 := FromArray[int]([]int{1, 2})
	s2 := FromArray[int]([]int{3, 4})
	result := Concat[int](s1, s2)
	assert.Equal(t, []int{1, 2, 3, 4}, result.ToArray())
}

func TestZip(t *testing.T) {
	result := Zip[int, string, string](
		[]int{1, 2, 3},
		[]string{"a", "b", "c"},
		func(i int, s string) string {
			return Mappers().IntToString()(i) + s
		},
	)
	assert.Equal(t, []string{"1a", "2b", "3c"}, result.ToArray())
}

func TestZip_DifferentLengths(t *testing.T) {
	result := Zip[int, int, int](
		[]int{1, 2, 3, 4, 5},
		[]int{10, 20},
		func(a, b int) int { return a + b },
	)
	assert.Equal(t, []int{11, 22}, result.ToArray())
}

func TestGroupBy(t *testing.T) {
	result := GroupBy[string, string](
		[]string{"apple", "avocado", "banana", "blueberry", "cherry"},
		func(s string) string { return string(s[0]) },
	)
	assert.Len(t, result["a"], 2)
	assert.Len(t, result["b"], 2)
	assert.Len(t, result["c"], 1)
}

func TestGroupBy_Empty(t *testing.T) {
	result := GroupBy[int, int]([]int{}, func(x int) int { return x })
	assert.Empty(t, result)
}

func TestDistinctBy(t *testing.T) {
	result := DistinctBy[string, byte](
		[]string{"apple", "avocado", "banana", "blueberry", "cherry"},
		func(s string) byte { return s[0] },
	)
	assert.Equal(t, 3, result.Len())
	assert.Equal(t, "apple", result.ToArray()[0])
	assert.Equal(t, "banana", result.ToArray()[1])
	assert.Equal(t, "cherry", result.ToArray()[2])
}

func TestMin(t *testing.T) {
	val, ok := Min[int]([]int{5, 3, 8, 1, 4}, func(a, b int) bool { return a < b })
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestMin_Empty(t *testing.T) {
	_, ok := Min[int]([]int{}, func(a, b int) bool { return a < b })
	assert.False(t, ok)
}

func TestMin_SingleElement(t *testing.T) {
	val, ok := Min[int]([]int{42}, func(a, b int) bool { return a < b })
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestMax(t *testing.T) {
	val, ok := Max[int]([]int{5, 3, 8, 1, 4}, func(a, b int) bool { return a < b })
	assert.True(t, ok)
	assert.Equal(t, 8, val)
}

func TestMax_Empty(t *testing.T) {
	_, ok := Max[int]([]int{}, func(a, b int) bool { return a < b })
	assert.False(t, ok)
}

func TestSum_Int(t *testing.T) {
	result := Sum[int]([]int{1, 2, 3, 4, 5})
	assert.Equal(t, 15, result)
}

func TestSum_Float(t *testing.T) {
	result := Sum[float64]([]float64{1.5, 2.5, 3.0})
	assert.Equal(t, 7.0, result)
}

func TestSum_Empty(t *testing.T) {
	result := Sum[int]([]int{})
	assert.Equal(t, 0, result)
}

func TestAverage_Int(t *testing.T) {
	avg, ok := Average[int]([]int{10, 20, 30})
	assert.True(t, ok)
	assert.Equal(t, 20.0, avg)
}

func TestAverage_Float(t *testing.T) {
	avg, ok := Average[float64]([]float64{1.0, 2.0, 3.0})
	assert.True(t, ok)
	assert.Equal(t, 2.0, avg)
}

func TestAverage_Empty(t *testing.T) {
	_, ok := Average[int]([]int{})
	assert.False(t, ok)
}

func TestAverage_Int8(t *testing.T) {
	avg, ok := Average[int8]([]int8{2, 4, 6})
	assert.True(t, ok)
	assert.Equal(t, 4.0, avg)
}

func TestAverage_Int16(t *testing.T) {
	avg, ok := Average[int16]([]int16{10, 20})
	assert.True(t, ok)
	assert.Equal(t, 15.0, avg)
}

func TestAverage_Int32(t *testing.T) {
	avg, ok := Average[int32]([]int32{3, 6, 9})
	assert.True(t, ok)
	assert.Equal(t, 6.0, avg)
}

func TestAverage_Int64(t *testing.T) {
	avg, ok := Average[int64]([]int64{100, 200})
	assert.True(t, ok)
	assert.Equal(t, 150.0, avg)
}

func TestAverage_Uint(t *testing.T) {
	avg, ok := Average[uint]([]uint{5, 10, 15})
	assert.True(t, ok)
	assert.Equal(t, 10.0, avg)
}

func TestAverage_Uint8(t *testing.T) {
	avg, ok := Average[uint8]([]uint8{2, 4})
	assert.True(t, ok)
	assert.Equal(t, 3.0, avg)
}

func TestAverage_Uint16(t *testing.T) {
	avg, ok := Average[uint16]([]uint16{10, 20})
	assert.True(t, ok)
	assert.Equal(t, 15.0, avg)
}

func TestAverage_Uint32(t *testing.T) {
	avg, ok := Average[uint32]([]uint32{6, 12})
	assert.True(t, ok)
	assert.Equal(t, 9.0, avg)
}

func TestAverage_Uint64(t *testing.T) {
	avg, ok := Average[uint64]([]uint64{50, 100})
	assert.True(t, ok)
	assert.Equal(t, 75.0, avg)
}

func TestAverage_Float32(t *testing.T) {
	avg, ok := Average[float32]([]float32{1.0, 3.0})
	assert.True(t, ok)
	assert.InDelta(t, 2.0, avg, 0.001)
}

func TestRange(t *testing.T) {
	result := Range(0, 5).ToArray()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
}

func TestRange_Negative(t *testing.T) {
	result := Range(5, 0)
	assert.Equal(t, 0, result.Count())
}

func TestRange_Same(t *testing.T) {
	result := Range(3, 3)
	assert.Equal(t, 0, result.Count())
}

func TestRepeat(t *testing.T) {
	result := Repeat[string]("hello", 3).ToArray()
	assert.Equal(t, []string{"hello", "hello", "hello"}, result)
}

func TestRepeat_Zero(t *testing.T) {
	result := Repeat[int](42, 0)
	assert.Equal(t, 0, result.Count())
}

func TestRepeat_Negative(t *testing.T) {
	result := Repeat[int](42, -1)
	assert.Equal(t, 0, result.Count())
}

func TestUnion(t *testing.T) {
	result := Union[int]([]int{1, 2, 3}, []int{3, 4, 5})
	assert.Equal(t, 5, result.Len())
	arr := result.ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr)
}

func TestUnion_NoOverlap(t *testing.T) {
	result := Union[int]([]int{1, 2}, []int{3, 4})
	assert.Equal(t, 4, result.Len())
}

func TestIntersect(t *testing.T) {
	result := Intersect[int]([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})
	assert.Equal(t, 2, result.Len())
	assert.Contains(t, result.ToArray(), 3)
	assert.Contains(t, result.ToArray(), 4)
}

func TestIntersect_NoOverlap(t *testing.T) {
	result := Intersect[int]([]int{1, 2}, []int{3, 4})
	assert.Equal(t, 0, result.Len())
}

func TestResolveIterator_Panics(t *testing.T) {
	assert.Panics(t, func() {
		resolveIterator[int]("invalid")
	})
}

func TestResolveIterator_NilStream(t *testing.T) {
	stream := FromCollection[int](NewList[int]())
	iter := resolveIterator[int](stream)
	assert.False(t, iter.HasNext())
}
