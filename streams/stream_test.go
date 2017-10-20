package streams

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var testArray = []string{"peach", "apple", "pear", "plum", "pineapple", "banana", "kiwi", "orange"}

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

func TestStream_Map(t *testing.T) {
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
