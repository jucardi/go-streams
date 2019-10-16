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
	From(testArray).
		Filter(func(v interface{}) bool {
			return strings.HasPrefix(v.(string), "p")
		}).
		OrderBy(func(a interface{}, b interface{}) int {
			return strings.Compare(a.(string), b.(string))
		}).
		ForEach(func(x interface{}) {
			println(x.(string))
		})
}

func TestStream_Contains(t *testing.T) {
	contains := From(testArray).Contains("apple")
	assert.True(t, contains)
}

func TestStream_AllMatch(t *testing.T) {
	var trueFunc = func(x interface{}) bool {
		return true
	}

	var appleFunc = func(x interface{}) bool {
		return "apple" == x
	}

	allMatch := From(testArray).AllMatch(trueFunc)
	notAllMatch := From(testArray).AllMatch(appleFunc)
	assert.True(t, allMatch)
	assert.False(t, notAllMatch)
}

func TestStream_AnyMatch(t *testing.T) {
	var appleFunc = func(x interface{}) bool {
		return "apple" == x
	}

	match := From(testArray).AnyMatch(appleFunc)
	assert.True(t, match)
}

func TestStream_NoneMatch(t *testing.T) {
	var falseFunc = func(x interface{}) bool {
		return false
	}

	noneMatch := From(testArray).NoneMatch(falseFunc)
	assert.True(t, noneMatch)
}

func TestStream_Filter(t *testing.T) {
	var appleFunc = func(x interface{}) bool {
		return "apple" == x
	}

	stream := From(testArray).
		Filter(appleFunc)

	assert.Equal(t, 1, stream.Count())
	assert.Equal(t, "apple", stream.First())
}

func TestStream_Except(t *testing.T) {
	var appleFunc = func(x interface{}) bool {
		return "apple" != x
	}

	stream := From(testArray).
		Except(appleFunc)

	assert.Equal(t, 1, stream.Count())
	assert.Equal(t, "apple", stream.First())
}

func TestStream_FirstAndLast(t *testing.T) {
	stream := From(testArray)
	emptyStream := From([]string{})

	assert.Equal(t, 8, stream.Count())
	assert.Equal(t, "peach", stream.First())
	assert.Equal(t, "peach", stream.First("some-value"))
	assert.Equal(t, "orange", stream.Last())
	assert.Equal(t, "orange", stream.Last("some-value"))

	assert.Equal(t, 0, emptyStream.Count())
	assert.Nil(t, emptyStream.First())
	assert.Equal(t, "some-value", emptyStream.First("some-value"))
	assert.Nil(t, emptyStream.Last())
	assert.Equal(t, "some-value", emptyStream.Last("some-value"))
}

func TestStream_MapFn(t *testing.T) {
	mapFunc := func(i interface{}) interface{} {
		return 5
	}

	stream := From(testArray).Map(mapFunc)

	assert.Equal(t, len(testArray), stream.Count())
	assert.Equal(t, 5, stream.First())
}

func TestStream_ForEach(t *testing.T) {
	buffer1 := new(bytes.Buffer)

	for _, v := range testArray {
		buffer1.WriteString(v)
	}

	buffer2 := new(bytes.Buffer)

	From(testArray).ForEach(func(v interface{}) {
		buffer2.WriteString(v.(string))
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

	stream := From(bigArray)
	stream.ParallelForEach(func(v interface{}) {
		v.(*testObj).Processed = true
	}, 0)

	for _, v := range bigArray {
		assert.True(t, v.Processed)
	}
}

// This test may fail when running with coverage with IntelliJ due to the coverage capture that may affect
// the performance of go channels. Running normally on a 2 CPU host, demonstrates an efficiency of around 200% vs non-parallel.
func TestStream_ParallelFiltering(t *testing.T) {
	cores := getCores(-1)

	// Skip this test if the machine only has one available CPU.
	if cores == 1 {
		return
	}

	sampleSize := 5000000
	bigArray := make([]int, sampleSize)

	for i := 0; i < sampleSize; i++ {
		bigArray[i] = rand.Intn(100)
	}

	filter1 := func(v interface{}) bool {
		return v.(int) < 50
	}

	filter2 := func(v interface{}) bool {
		return v.(int) < 10
	}

	const timesToTry = 5

	var avElapsed1, avElapsed2 time.Duration
	successTimes := 0

	for i := 0; i < timesToTry; i++ {
		start := time.Now()
		result1 := From(bigArray).Filter(filter1).Filter(filter2).ToArray().([]int)
		elapsed1 := time.Since(start)

		start = time.Now()
		result2 := From(bigArray, -1).Filter(filter1).Filter(filter2).ToArray().([]int)
		elapsed2 := time.Since(start)

		if elapsed1 > elapsed2 {
			successTimes++
		}

		avElapsed1 += elapsed1
		avElapsed2 += elapsed2

		// Validates than both results were the same
		assert.Equal(t, len(result1), len(result2))
	}

	fmt.Println("Non-Parallel Filtering average time: ", avElapsed1/timesToTry)
	fmt.Println("Parallel Filtering average time:  ", avElapsed2/timesToTry)
	fmt.Println("Parallel took", int64(100*avElapsed2/avElapsed1), "% of Non-Parallel time")

	// Validates parallel filtering was faster than non-parallel
	assert.True(t, avElapsed1/timesToTry > avElapsed2/timesToTry)
}

func TestStream_OrderBy(t *testing.T) {
	sortFn := func(a interface{}, b interface{}) int {
		return strings.Compare(a.(string), b.(string))
	}

	expected := []string{"apple", "banana", "kiwi", "orange", "peach", "pear", "pineapple", "plum"}
	sorted := From(testArray).OrderBy(sortFn).ToArray().([]string)

	assert.Equal(t, expected, sorted)
}

func TestStream_OrderByDesc(t *testing.T) {
	sortFn := func(a interface{}, b interface{}) int {
		return strings.Compare(a.(string), b.(string))
	}

	expected := []string{"plum", "pineapple", "pear", "peach", "orange", "kiwi", "banana", "apple"}
	sorted := From(testArray).OrderBy(sortFn, true).ToArray().([]string)

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

	result := From(m).
		Filter(func(i interface{}) bool {
			// assert the elements in the filter are *KeyValuePair
			assert.Equal(t, keyValuePairType, reflect.TypeOf(i))

			pair := i.(*KeyValuePair)
			return pair.Key.(string) != "d"
		}).
		OrderBy(func(a interface{}, b interface{}) int {
			// assert the elements in the filter are *KeyValuePair
			assert.Equal(t, keyValuePairType, reflect.TypeOf(a))
			assert.Equal(t, keyValuePairType, reflect.TypeOf(b))
			aPair, bPair := a.(*KeyValuePair), b.(*KeyValuePair)

			x, _ := strconv.Atoi(aPair.Value.(string))
			y, _ := strconv.Atoi(bPair.Value.(string))
			return x - y
		}).
		ToArray()

	jsonResult, _ := json.MarshalIndent(result, "", " ")
	assert.JSONEq(t, `[{ "Key": "e", "Value": "111" }, { "Key": "a", "Value": "123" }, { "Key": "b", "Value": "456" }, { "Key": "c", "Value": "789" }]`, string(jsonResult))
}

func TestStream_Distinct(t *testing.T) {
	arr := append(testArray, "apple", "banana", "kiwi", "apple", "banana", "apple")
	result := From(arr).Distinct().ToArray().([]string)

	assert.Equal(t, len(result), len(testArray))
}

func TestStream_DistinctWithSort(t *testing.T) {
	arr := append(testArray, "apple", "banana", "kiwi", "apple", "banana", "apple")
	sortFn := func(a interface{}, b interface{}) int {
		return strings.Compare(a.(string), b.(string))
	}

	expected := []string{"apple", "banana", "kiwi", "orange", "peach", "pear", "pineapple", "plum"}
	sorted := From(arr).OrderBy(sortFn).Distinct().ToArray().([]string)

	assert.Equal(t, expected, sorted)
}
