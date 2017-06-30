package streams

import (
	"reflect"
)

type Stream struct {
	arrayType   reflect.Type
	array       reflect.Value
	elementType reflect.Type
	filters     []func(reflect.Value) bool
	exceptions  []func(reflect.Value) bool
}

func From(array interface{}) *Stream {
	arrayType := reflect.TypeOf(array)

	if arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.Array {
		panic("Unable to create Stream from a none Slice or none Array")
	}

	return &Stream{
		array:       reflect.ValueOf(array),
		arrayType:   arrayType,
		elementType: reflect.TypeOf(array).Elem(),
	}
}

func (s *Stream) Filter(f func(reflect.Value) bool) *Stream {
	s.filters = append(s.filters, f)
	return s
}

func (s *Stream) Except(f func(reflect.Value) bool) *Stream {
	s.exceptions = append(s.exceptions, f)
	return s
}

func (s *Stream) First() reflect.Value {
	if filtered := s.start(); filtered.Len() > 0 {
		return filtered.Index(0)
	} else {
		return reflect.ValueOf(nil)
	}
}

func (s *Stream) Count() int {
	return s.start().Len()
}

func (s *Stream) AnyMatch(f func(reflect.Value) bool) bool {
	return s.filterHandler(s.start(), []func(reflect.Value) bool{f}, false).Len() > 0
}

func (s *Stream) AllMatch(f func(reflect.Value) bool) bool {
	array := s.start()
	return array.Len() == s.filterHandler(array, []func(reflect.Value) bool{f}, false).Len()
}

func (s *Stream) NoneMatch(f func(reflect.Value) bool) bool {
	return !s.AnyMatch(f)
}

func (s *Stream) Contains(value interface{}) bool {
	return s.AnyMatch(func(val reflect.Value) bool {
		return value == val.Interface()
	})
}

// region Private functions

func (s *Stream) start() reflect.Value {
	var array = s.array
	array = s.filter(array)
	array = s.except(array)
	return array
}

func (s *Stream) filter(array reflect.Value) reflect.Value {
	return s.filterHandler(array, s.filters, false)
}

func (s *Stream) except(array reflect.Value) reflect.Value {
	return s.filterHandler(array, s.exceptions, true)
}

func (s *Stream) filterHandler(array reflect.Value, filters []func(reflect.Value) bool, negate bool) reflect.Value {
	if len(filters) == 0 {
		return array
	}

	ret := reflect.MakeSlice(array.Type(), 0, 0)

	for i := 0; i < array.Len(); i++ {
		x := array.Index(i)
		var match bool = true

		for _, f := range filters {
			if negate {
				match = match && !f(x)
			} else {
				match = match && f(x)
			}

			if !match {
				break
			}
		}

		if match {
			ret = reflect.Append(ret, x)
		}
	}

	return ret
}

// endregion

// STREAM
//   Filter, Where     DONE
//   Except            DONE
//   Map
//   Distinct
//   Sorted, OrderBy    (orders by key)
//   Sorted, OrderBy    (comparator)
//   OrderByDescending
//   ThenBy
//   ThenByDescending
//   Reverse

// VOID
//    ForEach

// ARRAY
//    ToArray

// INT
//    Count, Size

// BOOLEAN
//    AnyMatch, Any     DONE
//    AllMatch, All     DONE
//    NoneMatch         DONE
//    Contains          DONE  -> Like Any, but instead of receiving a func, receives an element to perform an equals operation

// OPTIONAL ?? or element
//    Min
//    Max
//    Average
//    FindFirst, First         DONE
//    FindAny                  For parallel operations. Post MVP
//    ElementAt, At
//    ElementAtOrDefault, AtOrDefault
//    Last   (no args, returns last, func args, iterates from back to forth and returns the first match)
//    LastOrDefault

// Concat --> Concatenates two sequences
// Reduce, Aggregate       --->   Sum, min, max, average, string concatenation, with and without seed value
// Skip(long n) -> skips the first N elements.
// Peek -> iterates and does something returning back the stream. Mainly for debugging
// Limit -> limits the size of the stream.

// GROUP OPERATIONS
//    GroupBy
//    GroupJoin
//    Intersect    (default equals or with comparer function)
//    Union
