package streams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testArray = []string{"peach", "apple", "pear", "plum", "pineapple", "banana", "kiwi", "orange"}

func TestSample(t *testing.T) {
	From[string](testArray).
		Filter(func(v string) bool {
			return strings.HasPrefix(v, "p")
		}).
		Sort(strings.Compare).
		ForEach(func(x string) {
			println(x)
		})
}

func TestStream_Contains(t *testing.T) {
	contains := From[string](testArray).Contains("apple")
	assert.True(t, contains)
}

func TestStream_AllMatch(t *testing.T) {
	var trueFunc = func(x string) bool {
		return true
	}

	var appleFunc = func(x string) bool {
		return "apple" == x
	}

	allMatch := From[string](testArray).AllMatch(trueFunc)
	notAllMatch := From[string](testArray).AllMatch(appleFunc)
	assert.True(t, allMatch)
	assert.False(t, notAllMatch)
}

func TestStream_AnyMatch(t *testing.T) {
	var appleFunc = func(x string) bool {
		return "apple" == x
	}

	match := From[string](testArray).AnyMatch(appleFunc)
	assert.True(t, match)
}

func TestStream_NoneMatch(t *testing.T) {
	var falseFunc = func(x string) bool {
		return false
	}

	noneMatch := From[string](testArray).NoneMatch(falseFunc)
	assert.True(t, noneMatch)
}

func TestStream_Filter(t *testing.T) {
	var appleFunc = func(x string) bool {
		return "apple" == x
	}

	stream := From[string](testArray).
		Filter(appleFunc)

	assert.Equal(t, 1, stream.Count())
	assert.Equal(t, "apple", stream.First())
}

func TestStream_Except(t *testing.T) {
	var appleFunc = func(x string) bool {
		return "apple" != x
	}

	stream := From[string](testArray).
		Except(appleFunc)

	assert.Equal(t, 1, stream.Count())
	assert.Equal(t, "apple", stream.First())
}

func TestStream_FirstAndLast(t *testing.T) {
	stream := From[string](testArray)
	emptyStream := From[string]([]string{})

	assert.Equal(t, 8, stream.Count())
	assert.Equal(t, "peach", stream.First())
	assert.Equal(t, "peach", stream.First("some-value"))
	assert.Equal(t, "orange", stream.Last())
	assert.Equal(t, "orange", stream.Last("some-value"))

	assert.Equal(t, 0, emptyStream.Count())
	assert.Empty(t, emptyStream.First())
	assert.Equal(t, "some-value", emptyStream.First("some-value"))
	assert.Empty(t, emptyStream.Last())
	assert.Equal(t, "some-value", emptyStream.Last("some-value"))
}

func TestStream_MapFn(t *testing.T) {
	mapFunc := func(i string) int {
		return 5
	}

	stream := Map[string, int](
		From[string](testArray), mapFunc).Stream()

	assert.Equal(t, len(testArray), stream.Count())
	assert.Equal(t, 5, stream.First())
}

func TestStream_ForEach(t *testing.T) {
	buffer1 := new(bytes.Buffer)

	for _, v := range testArray {
		buffer1.WriteString(v)
	}

	buffer2 := new(bytes.Buffer)

	From[string](testArray).ForEach(func(v string) {
		buffer2.WriteString(v)
	})

	assert.Equal(t, buffer1.String(), buffer2.String())
}

func TestStream_ParallelForEach(t *testing.T) {
	sampleSize := 10000

	type testObj struct {
		Index     int
		Processed bool
	}

	bigArray := make([]*testObj, sampleSize)

	for i := 0; i < sampleSize; i++ {
		bigArray[i] = &testObj{
			Index:     i,
			Processed: false,
		}
	}

	stream := From[*testObj](bigArray)
	stream.ParallelForEach(func(v *testObj) {
		v.Processed = true
	}, 0)

	for _, v := range bigArray {
		assert.True(t, v.Processed)
	}
}

// This test may fail when running with coverage with IntelliJ due to the coverage capture that may affect
// the performance of go channels. Running normally on a 2 CPU host, demonstrates an efficiency of around 200% vs non-parallel.
func TestStream_ParallelFiltering(t *testing.T) {
	testParallelFiltering(t, 100, true)
	testParallelFiltering(t, 50000, false)
}

func testParallelFiltering(t *testing.T, sampleSize int, sort bool) {
	cores := getCores(-1)

	// Skip this test if the machine only has one available CPU.
	if cores == 1 {
		return
	}

	bigArray := make([]int, sampleSize)

	for i := 0; i < sampleSize; i++ {
		bigArray[i] = rand.Intn(100)
	}

	filter1 := func(v int) bool {
		return v < 50
	}

	filter2 := func(v int) bool {
		return v > 10
	}

	const timesToTry = 5

	var avElapsed1, avElapsed2 time.Duration
	successTimes := 0

	for i := 0; i < timesToTry; i++ {
		start := time.Now()
		var result1, result2 []int
		if sort {
			result1 = From[int](bigArray).Filter(filter1).Filter(filter2).Sort(ComparableFn[int]()).ToArray()
		} else {
			result1 = From[int](bigArray).Filter(filter1).Filter(filter2).ToArray()
		}

		elapsed1 := time.Since(start)

		start = time.Now()
		if sort {
			result2 = From[int](bigArray, -1).Filter(filter1).Filter(filter2).Sort(ComparableFn[int]()).ToArray()
		} else {
			result2 = From[int](bigArray, -1).Filter(filter1).Filter(filter2).ToArray()
		}

		elapsed2 := time.Since(start)

		if elapsed1 > elapsed2 {
			successTimes++
		}

		avElapsed1 += elapsed1
		avElapsed2 += elapsed2

		// Validates than both results were the same
		if sort {
			assert.EqualValues(t, result1, result2)
		} else {
			assert.Equal(t, len(result1), len(result2))
		}
	}

	fmt.Printf("\n\nWith sample of %d and sorting:%v\n", sampleSize, sort)
	fmt.Println("  Non-Parallel Filtering average time: ", avElapsed1/timesToTry)
	fmt.Println("  Parallel Filtering average time:  ", avElapsed2/timesToTry)
	fmt.Println("  Parallel took", int64(100*avElapsed2/avElapsed1), "% of Non-Parallel time")
	println()

	// Validates parallel filtering was faster than non-parallel
	// assert.True(t, avElapsed1/timesToTry > avElapsed2/timesToTry)
}

func TestStream_OrderBy(t *testing.T) {
	expected := []string{"apple", "banana", "kiwi", "orange", "peach", "pear", "pineapple", "plum"}
	sorted := From[string](testArray).Sort(strings.Compare).ToArray()

	assert.Equal(t, expected, sorted)
}

func TestStream_OrderByDesc(t *testing.T) {
	expected := []string{"plum", "pineapple", "pear", "peach", "orange", "kiwi", "banana", "apple"}
	sorted := From[string](testArray).Sort(strings.Compare, true).ToArray()

	assert.Equal(t, expected, sorted)
}

func TestStream_MapCollection(t *testing.T) {
	m := map[string]string{
		"a": "123",
		"b": "456",
		"c": "789",
		"d": "000",
		"e": "111",
	}

	result := FromMap[string, string](m).
		Filter(func(pair *KeyValuePair[string, string]) bool {
			return pair.Key != "d"
		}).
		Sort(func(a, b *KeyValuePair[string, string]) int {
			x, _ := strconv.Atoi(a.Value)
			y, _ := strconv.Atoi(b.Value)
			return x - y
		}).
		ToArray()

	jsonResult, _ := json.MarshalIndent(result, "", " ")
	assert.JSONEq(t, `[{ "Key": "e", "Value": "111" }, { "Key": "a", "Value": "123" }, { "Key": "b", "Value": "456" }, { "Key": "c", "Value": "789" }]`, string(jsonResult))
}

func TestStream_Distinct(t *testing.T) {
	arr := append(testArray, "apple", "banana", "kiwi", "apple", "banana", "apple")
	result := From[string](arr).Distinct().ToArray()

	assert.Equal(t, len(result), len(testArray))
}

func TestStream_DistinctWithSort(t *testing.T) {
	arr := append(testArray, "apple", "banana", "kiwi", "apple", "banana", "apple")
	expected := []string{"apple", "banana", "kiwi", "orange", "peach", "pear", "pineapple", "plum"}
	sorted := From[string](arr).Sort(strings.Compare).Distinct().ToArray()
	assert.Equal(t, expected, sorted)
}

type testStruct struct {
	A string
	B int
	C string
}

func TestMapToPtr(t *testing.T) {
	arr := []testStruct{
		{
			A: "abcde",
			B: 1234,
			C: "zxcv",
		},
		{
			A: "something",
			B: 54321,
			C: "something_else",
		},
	}
	ret := MapToPtr[testStruct](arr)

	sourceData, _ := json.Marshal(arr)
	targetData, _ := json.Marshal(ret)
	assert.JSONEq(t, string(sourceData), string(targetData))
	assert.Equal(t, "[]streams.testStruct", reflect.TypeOf(arr).String())
	assert.Equal(t, "[]*streams.testStruct", reflect.TypeOf(ret).String())
}

func TestStream_SetThreads(t *testing.T) {
	stream := FromArray[int]([]int{1, 2, 3, 4, 5}).SetThreads(2)
	result := stream.Filter(func(x int) bool { return x > 2 }).ToArray()
	assert.Len(t, result, 3)
}

func TestStream_NotAllMatch(t *testing.T) {
	result := From[int]([]int{1, 2, 3}).NotAllMatch(func(x int) bool { return x > 2 })
	assert.True(t, result)

	result = From[int]([]int{3, 4, 5}).NotAllMatch(func(x int) bool { return x > 2 })
	assert.False(t, result)
}

func TestStream_IsEmpty(t *testing.T) {
	assert.True(t, FromArray[int]([]int{}).IsEmpty())
	assert.False(t, FromArray[int]([]int{1}).IsEmpty())
}

func TestStream_At(t *testing.T) {
	stream := FromArray[int]([]int{10, 20, 30, 40, 50})
	assert.Equal(t, 10, stream.At(0))
	assert.Equal(t, 30, stream.At(2))
	assert.Equal(t, 50, stream.At(4))
}

func TestStream_At_OutOfBounds(t *testing.T) {
	stream := FromArray[int]([]int{10, 20})
	assert.Equal(t, 0, stream.At(-1))
	assert.Equal(t, 0, stream.At(5))
	assert.Equal(t, 99, stream.At(5, 99))
}

func TestStream_AtReverse(t *testing.T) {
	stream := FromArray[int]([]int{10, 20, 30, 40, 50})
	assert.Equal(t, 50, stream.AtReverse(0))
	assert.Equal(t, 40, stream.AtReverse(1))
	assert.Equal(t, 10, stream.AtReverse(4))
}

func TestStream_AtReverse_OutOfBounds(t *testing.T) {
	stream := FromArray[int]([]int{10, 20})
	assert.Equal(t, 0, stream.AtReverse(5))
	assert.Equal(t, 99, stream.AtReverse(5, 99))
}

func TestStream_ToIterable(t *testing.T) {
	stream := FromArray[int]([]int{1, 2, 3})
	iterable := stream.ToIterable()
	assert.NotNil(t, iterable)
}

func TestStream_ToList(t *testing.T) {
	stream := FromArray[int]([]int{1, 2, 3})
	list := stream.ToList()
	assert.Equal(t, 3, list.Len())
	assert.Equal(t, []int{1, 2, 3}, list.ToArray())
}

func TestStream_ToList_NilIterable(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	list := s.ToList()
	assert.Equal(t, 0, list.Len())
}

func TestStream_ToDistinct(t *testing.T) {
	arr := []int{1, 2, 2, 3, 3, 3}
	set := FromArray[int](arr).ToDistinct()
	assert.Equal(t, 3, set.Len())
	assert.True(t, set.Contains(1, 2, 3))
}

func TestStream_ToDistinct_Empty(t *testing.T) {
	set := FromArray[int]([]int{}).ToDistinct()
	assert.Equal(t, 0, set.Len())
}

func TestStream_Skip(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5}).Skip(2).ToArray()
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestStream_Skip_Zero(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Skip(0).ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStream_Skip_BeyondLength(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Skip(10).ToArray()
	assert.Empty(t, result)
}

func TestStream_Limit(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5}).Limit(3).ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStream_Limit_Zero(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Limit(0).ToArray()
	assert.Empty(t, result)
}

func TestStream_Limit_BeyondLength(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Limit(10).ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStream_Reverse(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5}).Reverse().ToArray()
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}

func TestStream_Reverse_Empty(t *testing.T) {
	result := FromArray[int]([]int{}).Reverse().ToArray()
	assert.Empty(t, result)
}

func TestStream_Reverse_Single(t *testing.T) {
	result := FromArray[int]([]int{42}).Reverse().ToArray()
	assert.Equal(t, []int{42}, result)
}

func TestStream_Peek(t *testing.T) {
	var peeked []int
	result := FromArray[int]([]int{1, 2, 3}).
		Peek(func(x int) { peeked = append(peeked, x) }).
		ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, []int{1, 2, 3}, peeked)
}

func TestStream_Peek_Empty(t *testing.T) {
	var peeked []int
	result := FromArray[int]([]int{}).
		Peek(func(x int) { peeked = append(peeked, x) }).
		ToArray()
	assert.Empty(t, result)
	assert.Empty(t, peeked)
}

func TestStream_TakeWhile(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5}).
		TakeWhile(func(x int) bool { return x < 4 }).
		ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStream_TakeWhile_AllMatch(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).
		TakeWhile(func(x int) bool { return x < 10 }).
		ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStream_TakeWhile_NoneMatch(t *testing.T) {
	result := FromArray[int]([]int{5, 6, 7}).
		TakeWhile(func(x int) bool { return x < 5 }).
		ToArray()
	assert.Empty(t, result)
}

func TestStream_SkipWhile(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5}).
		SkipWhile(func(x int) bool { return x < 3 }).
		ToArray()
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestStream_SkipWhile_AllSkipped(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).
		SkipWhile(func(x int) bool { return x < 10 }).
		ToArray()
	assert.Empty(t, result)
}

func TestStream_SkipWhile_NoneSkipped(t *testing.T) {
	result := FromArray[int]([]int{5, 6, 7}).
		SkipWhile(func(x int) bool { return x < 5 }).
		ToArray()
	assert.Equal(t, []int{5, 6, 7}, result)
}

func TestStream_Chunk(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5}).Chunk(2)
	assert.Len(t, result, 3)
	assert.Equal(t, []int{1, 2}, result[0])
	assert.Equal(t, []int{3, 4}, result[1])
	assert.Equal(t, []int{5}, result[2])
}

func TestStream_Chunk_ExactDivision(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4}).Chunk(2)
	assert.Len(t, result, 2)
	assert.Equal(t, []int{1, 2}, result[0])
	assert.Equal(t, []int{3, 4}, result[1])
}

func TestStream_Chunk_SizeLargerThanStream(t *testing.T) {
	result := FromArray[int]([]int{1, 2}).Chunk(10)
	assert.Len(t, result, 1)
	assert.Equal(t, []int{1, 2}, result[0])
}

func TestStream_Chunk_Empty(t *testing.T) {
	result := FromArray[int]([]int{}).Chunk(3)
	assert.Nil(t, result)
}

func TestStream_Chunk_ZeroSize(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Chunk(0)
	assert.Nil(t, result)
}

func TestStream_NilIterable_Count(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	assert.Equal(t, 0, s.Count())
}

func TestStream_NilIterable_AnyMatch(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	assert.False(t, s.AnyMatch(func(x int) bool { return true }))
}

func TestStream_NilIterable_AllMatch(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	assert.True(t, s.AllMatch(func(x int) bool { return true }))
}

func TestStream_NilIterable_ForEach(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	s.ForEach(func(x int) {
		t.Fatal("should not be called")
	})
}

func TestStream_NilIterable_ToArray(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	assert.Nil(t, s.ToArray())
}

func TestStream_NilIterable_ParallelForEach(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	s.ParallelForEach(func(x int) {
		t.Fatal("should not be called")
	}, 2)
}

func TestStream_Chained_Pipeline(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Filter(func(x int) bool { return x%2 == 0 }).
		Sort(ComparableFn[int](), true).
		Skip(1).
		Limit(3).
		ToArray()
	assert.Equal(t, []int{8, 6, 4}, result)
}

func TestStream_ProcessCaching(t *testing.T) {
	callCount := 0
	stream := FromArray[int]([]int{1, 2, 3}).
		Filter(func(x int) bool {
			callCount++
			return true
		})

	_ = stream.Count()
	count1 := callCount
	_ = stream.ToArray()
	assert.Equal(t, count1, callCount, "process() should cache results")
}

func TestStream_ParallelForEach_EmptyCollection(t *testing.T) {
	called := false
	FromArray[int]([]int{}).ParallelForEach(func(x int) {
		called = true
	}, 2)
	assert.False(t, called)
}

func TestStream_NilIterable_Skip(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.Skip(1).ToArray()
	assert.Empty(t, result)
}

func TestStream_NilIterable_Limit(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.Limit(1).ToArray()
	assert.Empty(t, result)
}

func TestStream_NilIterable_Reverse(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.Reverse().ToArray()
	assert.Empty(t, result)
}

func TestStream_NilIterable_Peek(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.Peek(func(x int) {}).ToArray()
	assert.Empty(t, result)
}

func TestStream_NilIterable_TakeWhile(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.TakeWhile(func(x int) bool { return true }).ToArray()
	assert.Empty(t, result)
}

func TestStream_NilIterable_SkipWhile(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.SkipWhile(func(x int) bool { return true }).ToArray()
	assert.Empty(t, result)
}

func TestStream_NilIterable_Chunk(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.Chunk(2)
	assert.Nil(t, result)
}

func TestStream_NilIterable_IfEmpty(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	thenCalled := false
	s.IfEmpty().Then(func(st IStream[int]) {
		thenCalled = true
	})
	assert.True(t, thenCalled)
}

func TestStream_NilIterable_ToDistinct(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	result := s.ToDistinct()
	assert.Equal(t, 0, result.Len())
}

func TestStream_AtReverse_NilIterable(t *testing.T) {
	s := &Stream[int]{iterable: nil}
	assert.Equal(t, 0, s.AtReverse(0))
	assert.Equal(t, 99, s.AtReverse(0, 99))
}

func TestStream_ParallelProcess_WithFilters(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, -1).
		Filter(func(x int) bool { return x%2 == 0 }).
		Sort(ComparableFn[int]()).
		ToArray()
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestStream_ParallelProcess_NoFilters(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}, -1).
		Sort(ComparableFn[int](), true).
		ToArray()
	assert.Equal(t, []int{3, 2, 1}, result)
}

func TestStream_ToDistinct_WithSort(t *testing.T) {
	result := FromArray[int]([]int{3, 1, 2, 1, 3}).Sort(ComparableFn[int]()).ToDistinct()
	assert.Equal(t, 3, result.Len())
}

func TestStream_Count_NegativeLen(t *testing.T) {
	s := FromArray[int]([]int{1, 2, 3})
	assert.Equal(t, 3, s.Count())
}

func TestStream_Skip_Negative(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Skip(-1).ToArray()
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStream_Limit_Negative(t *testing.T) {
	result := FromArray[int]([]int{1, 2, 3}).Limit(-1).ToArray()
	assert.Empty(t, result)
}

func TestStream_ParallelForEach_SkipWait(t *testing.T) {
	FromArray[int]([]int{1, 2, 3, 4, 5}).ParallelForEach(func(x int) {}, 2, true)
}
