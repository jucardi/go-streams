package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThen_ConditionMet(t *testing.T) {
	thenCalled := false
	elseCalled := false

	tw := &thenWrapper[int]{
		conditionMet: true,
		stream:       FromArray[int]([]int{1, 2, 3}),
	}

	tw.Then(func(s IStream[int]) {
		thenCalled = true
		assert.Equal(t, 3, s.Count())
	}).Else(func(s IStream[int]) {
		elseCalled = true
	})

	assert.True(t, thenCalled)
	assert.False(t, elseCalled)
}

func TestThen_ConditionNotMet(t *testing.T) {
	thenCalled := false
	elseCalled := false

	tw := &thenWrapper[int]{
		conditionMet: false,
		stream:       FromArray[int]([]int{1, 2, 3}),
	}

	tw.Then(func(s IStream[int]) {
		thenCalled = true
	}).Else(func(s IStream[int]) {
		elseCalled = true
		assert.Equal(t, 3, s.Count())
	})

	assert.False(t, thenCalled)
	assert.True(t, elseCalled)
}

func TestStream_IfEmpty_Empty(t *testing.T) {
	thenCalled := false
	FromArray[int]([]int{}).
		IfEmpty().
		Then(func(s IStream[int]) {
			thenCalled = true
		})
	assert.True(t, thenCalled)
}

func TestStream_IfEmpty_NotEmpty(t *testing.T) {
	elseCalled := false
	FromArray[int]([]int{1, 2}).
		IfEmpty().
		Then(func(s IStream[int]) {
			t.Fatal("should not be called")
		}).
		Else(func(s IStream[int]) {
			elseCalled = true
		})
	assert.True(t, elseCalled)
}

func TestStream_IfAnyMatch(t *testing.T) {
	thenCalled := false
	FromArray[int]([]int{1, 2, 3}).
		IfAnyMatch(func(x int) bool { return x == 2 }).
		Then(func(s IStream[int]) {
			thenCalled = true
		})
	assert.True(t, thenCalled)
}

func TestStream_IfAnyMatch_NoMatch(t *testing.T) {
	elseCalled := false
	FromArray[int]([]int{1, 2, 3}).
		IfAnyMatch(func(x int) bool { return x == 99 }).
		Then(func(s IStream[int]) {
			t.Fatal("should not be called")
		}).
		Else(func(s IStream[int]) {
			elseCalled = true
		})
	assert.True(t, elseCalled)
}

func TestStream_IfAllMatch(t *testing.T) {
	thenCalled := false
	FromArray[int]([]int{2, 4, 6}).
		IfAllMatch(func(x int) bool { return x%2 == 0 }).
		Then(func(s IStream[int]) {
			thenCalled = true
		})
	assert.True(t, thenCalled)
}

func TestStream_IfAllMatch_NotAll(t *testing.T) {
	elseCalled := false
	FromArray[int]([]int{2, 3, 6}).
		IfAllMatch(func(x int) bool { return x%2 == 0 }).
		Else(func(s IStream[int]) {
			elseCalled = true
		})
	assert.True(t, elseCalled)
}

func TestStream_IfNoneMatch(t *testing.T) {
	thenCalled := false
	FromArray[int]([]int{1, 3, 5}).
		IfNoneMatch(func(x int) bool { return x%2 == 0 }).
		Then(func(s IStream[int]) {
			thenCalled = true
		})
	assert.True(t, thenCalled)
}

func TestStream_IfNoneMatch_SomeMatch(t *testing.T) {
	elseCalled := false
	FromArray[int]([]int{1, 2, 5}).
		IfNoneMatch(func(x int) bool { return x%2 == 0 }).
		Else(func(s IStream[int]) {
			elseCalled = true
		})
	assert.True(t, elseCalled)
}
