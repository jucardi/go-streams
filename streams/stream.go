package streams

import (
	"reflect"
)

type Stream struct {
	array       reflect.Value
	elementType reflect.Type
	filters     []func(interface{}) bool
	exceptions  []func(interface{}) bool
}

// Creates a Stream from the given array
func From(array interface{}) *Stream {
	arrayType := reflect.TypeOf(array)

	if arrayType.Kind() != reflect.Slice && arrayType.Kind() != reflect.Array {
		panic("Unable to create Stream from a none Slice or none Array")
	}

	return &Stream{
		array:       reflect.ValueOf(array),
		elementType: reflect.TypeOf(array).Elem(),
	}
}

// Uses the specified function to filter elements within the Stream
func (s *Stream) Filter(f func(interface{}) bool) *Stream {
	s.filters = append(s.filters, f)
	return s
}

func (s *Stream) Except(f func(interface{}) bool) *Stream {
	s.exceptions = append(s.exceptions, f)
	return s
}

func (s *Stream) Map(f func(interface{}) interface{}) *Stream {
	array := s.start()
	var newArr reflect.Value

	for i := 0; i < array.Len(); i++ {
		old := array.Index(i)
		n := reflect.ValueOf(f(old.Interface()))

		if i == 0 {
			newArr = reflect.MakeSlice(reflect.SliceOf(n.Type()), 0, 0)
		}

		newArr = reflect.Append(newArr, n)
	}

	return From(newArr.Interface())
}

func (s *Stream) First() interface{} {
	if filtered := s.start(); filtered.Len() > 0 {
		return filtered.Index(0).Interface()
	} else {
		return nil
	}
}

func (s *Stream) Count() int {
	return s.start().Len()
}

func (s *Stream) AnyMatch(f func(interface{}) bool) bool {
	return s.filterHandler(s.start(), []func(interface{}) bool{f}, false).Len() > 0
}

func (s *Stream) AllMatch(f func(interface{}) bool) bool {
	array := s.start()
	return array.Len() == s.filterHandler(array, []func(interface{}) bool{f}, false).Len()
}

func (s *Stream) NoneMatch(f func(interface{}) bool) bool {
	return !s.AnyMatch(f)
}

func (s *Stream) Contains(value interface{}) bool {
	return s.AnyMatch(func(val interface{}) bool {
		return value == val
	})
}

func (s *Stream) ForEach(f func(interface {})) {
	array := s.start()

	for i := 0; i < array.Len(); i++ {
		val := array.Index(i)
		f(val.Interface())
	}
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

func (s *Stream) filterHandler(array reflect.Value, filters []func(interface{}) bool, negate bool) reflect.Value {
	if len(filters) == 0 {
		return array
	}

	ret := reflect.MakeSlice(array.Type(), 0, 0)

	for i := 0; i < array.Len(); i++ {
		x := array.Index(i)
		var match bool = true

		for _, f := range filters {
			if negate {
				match = match && !f(x.Interface())
			} else {
				match = match && f(x.Interface())
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

// OPTIONAL ?? or element
//    Min
//    Max
//    Average
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
//
//
//
// =========== DONE =============
//
// STREAM
//   Filter, Where     DONE
//   Except            DONE
//   Map               DONE
//
// BOOLEAN
//    AnyMatch, Any     DONE
//    AllMatch, All     DONE
//    NoneMatch         DONE
//    Contains          DONE  -> Like Any, but instead of receiving a func, receives an element to perform an equals operation
//
// INT
//    Count, Size       DONE
//
// ELEMENT
//
//    FindFirst, First         DONE