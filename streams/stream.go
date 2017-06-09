package streams

type Stream struct {
	array      []interface{}
	filters    []func(interface{}) bool
	exceptions []func(interface{}) bool
}


func From(array []interface{}) *Stream {
	return &Stream{array: array}
}

func (s *Stream) Filter(f func(interface{}) bool) *Stream {
	s.filters = append(s.filters, f)
	return s
}

func (s *Stream) Except(f func(interface{}) bool) *Stream {
	s.filters = append(s.exceptions, f)
	return s
}

// region Private functions

func (s *Stream) start() []interface{} {
	var array = s.array
	array = s.filter(array)
	array = s.except(array)

	return array
}

func (s *Stream) filter(array []interface{}) []interface{} {
	return genericFilter(array, s.filters, false)
}

func (s *Stream) except(array []interface{}) []interface{} {
	return genericFilter(array, s.exceptions, true)
}

func genericFilter(array []interface{}, filters []func(interface{}) bool, negate bool) []interface{} {
	if len(filters) == 0 {
		return array
	}

	var ret []interface{}

	for _, x := range array {
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
			ret = append(ret, x)
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
//    AnyMatch, Any
//    AllMatch, All
//    NoneMatch
//    Contains  -> Like Any, but instead of receiving a func, receives an element to perform an equals operation

// OPTIONAL ?? or element
//    Min
//    Max
//    Average
//    FindFirst, First
//    FindAny
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
