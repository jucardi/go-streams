package streams

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"reflect"
)

var testArray = []string{"peach", "apple", "pear", "plum"}

func TestStream_Contains(t *testing.T) {
	contains := From(testArray).Contains("apple")
	assert.True(t, contains)
}

func TestStream_AllMatch(t *testing.T) {
	var trueFunc = func(x reflect.Value) bool {
		return true
	}

	var appleFunc = func(x reflect.Value) bool {
		return "apple" == x.String()
	}

	allMatch := From(testArray).AllMatch(trueFunc)
	notAllMatch := From(testArray).AllMatch(appleFunc)
	assert.True(t, allMatch)
	assert.False(t, notAllMatch)
}

func TestStream_AnyMatch(t *testing.T) {
	var appleFunc = func(x reflect.Value) bool {
		return "apple" == x.String()
	}

	match := From(testArray).AnyMatch(appleFunc)
	assert.True(t, match)
}

func TestStream_NoneMatch(t *testing.T) {
	var falseFunc = func(x reflect.Value) bool {
		return false
	}

	noneMatch := From(testArray).NoneMatch(falseFunc)
	assert.True(t, noneMatch)
}

func TestStream_Filter(t *testing.T) {
	var appleFunc = func(x reflect.Value) bool {
		return "apple" == x.String()
	}

	stream := From(testArray).
		Filter(appleFunc)

	assert.Equal(t, 1, stream.Count())
	assert.Equal(t, "apple", stream.First().String())
}

func TestStream_Except(t *testing.T) {
	var appleFunc = func(x reflect.Value) bool {
		return "apple" != x.String()
	}

	stream := From(testArray).
		Except(appleFunc)

	assert.Equal(t, 1, stream.Count())
	assert.Equal(t, "apple", stream.First().String())
}
