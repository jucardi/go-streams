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
