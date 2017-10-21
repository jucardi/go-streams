package streams

import (
	"math"
	"reflect"
	"runtime"
	"sort"
	"sync"
)

type Stream struct {
	collection ICollection
	filters    []func(interface{}) bool
	exceptions []func(interface{}) bool
	sorts      []sortFunc
	threads    int
}

type sortFunc struct {
	fn   func(interface{}, interface{}) int
	desc bool
}

type sorter struct {
	collection ICollection
	sorts      []sortFunc
}

// From: Creates a Stream from a given collection or ICollection.
//
// - collection: The collection or ICollection to be used to create the stream
// - threads:    If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//               to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//               available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//               best combine it with a `SortBy`. Only needs to be provided once per stream.
func From(collection interface{}, threads ...int) *Stream {
	colReflectType := reflect.TypeOf((*ICollection)(nil)).Elem()

	if reflect.PtrTo(reflect.TypeOf(collection)).Implements(colReflectType) {
		return FromCollection(collection.(ICollection), threads...)
	}

	return FromArray(collection, threads...)
}

// From: Creates a Stream from a given collection.
//
// - collection:   The collection to be used to create the stream
// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//            best combine it with a `SortBy`. Only needs to be provided once per stream.
func FromArray(array interface{}, threads ...int) *Stream {
	return FromCollection(NewCollectionFromArray(array), threads...)
}

// From: Creates a Stream from a given ICollection.
//
// - collection: The ICollection to be used to create the stream
// - threads:    If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//               to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//               available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//               best combine it with a `SortBy`. Only needs to be provided once per stream.
func FromCollection(collection ICollection, threads ...int) *Stream {
	return &Stream{
		collection: collection,
		threads:    getCores(threads...),
	}
}

// SetThreads: Sets the amount of go channels to be used for parallel filtering to a maximum of the available CPUs in the
// host machine. Providing a value <= 0, indicates the maximum amount of available CPUs will be the number that determines
// the amount of go channels to be used. If order matters, best combine it with a `SortBy`. Only needs to be provided once
// per stream.
func (s *Stream) SetThreads(threads int) int {
	return s.updateCores(threads)
}

// Filter: Filters any element that does not meet the condition provided by the function.
//
// - f:       The filtering function to be used.
// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//            best combine it with a `SortBy`. Only needs to be provided once per stream.
func (s *Stream) Filter(f func(interface{}) bool, threads ...int) *Stream {
	s.updateCores(threads...)
	s.filters = append(s.filters, f)
	return s
}

// Except: Filters all elements that meet the condition provided by the function
//
// - f:       The filtering function to be used.
// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//            best combine it with a `SortBy`. Only needs to be provided once per stream.
func (s *Stream) Except(f func(interface{}) bool, threads ...int) *Stream {
	s.updateCores(threads...)
	s.exceptions = append(s.exceptions, f)
	return s
}

// Map: Maps the elements of the collection to a new element, using the mapping function provided
func (s *Stream) Map(f func(interface{}) interface{}) *Stream {
	iterable := s.start()
	var col ICollection

	for old := iterable.Current(); iterable.HasNext(); old = iterable.Next() {
		n := f(old)

		if col == nil {
			col = NewCollection(reflect.TypeOf(n))
		}

		col.Add(n)
	}

	return FromCollection(col)
}

// First: Returns the first element of the resulting stream. Nil if the resulting stream is empty.
func (s *Stream) First(defaultValue ...interface{}) interface{} {
	return s.At(0, defaultValue...)
}

// Last: Returns the last element of the resulting stream. Nil if the resulting stream is empty.
func (s *Stream) Last(defaultValue ...interface{}) interface{} {
	return s.AtReverse(0, defaultValue...)
}

// At: Returns the element at the given index in the resulting stream.
// Returns nil (or default value if provided) if out of bounds.
func (s *Stream) At(index int, defaultValue ...interface{}) interface{} {
	filtered := s.start()
	filtered.Skip(index)

	if val := filtered.Current(); val == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	} else {
		return val
	}
}

// AtReverse: Returns the element at the given position, starting from the last element to the first in the resulting stream.
// Returns Nil if out of bounds.
func (s *Stream) AtReverse(pos int, defaultValue ...interface{}) interface{} {
	filtered := s.start()
	i := filtered.Len() - 1 - pos

	if i >= 0 {
		return filtered.Index(i)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

// Count: Counts the elements of the resulting stream
func (s *Stream) Count() int {
	return s.start().Len()
}

// AnyMatch: Indicates whether any elements of the stream match the given condition function
//
// - f:       The matching function to be used.
func (s *Stream) AnyMatch(f func(interface{}) bool) bool {
	array := s.start()
	return anyMatch(array, 0, array.Len(), []func(interface{}) bool{f}, false)
}

// AllMatch: Indicates whether ALL elements of the stream match the given condition function
//
// - f:       The matching function to be used.
func (s *Stream) AllMatch(f func(interface{}) bool) bool {
	array := s.start()
	return !anyMatch(array, 0, array.Len(), []func(interface{}) bool{f}, true)
}

// NoneMatch:  Indicates whether NONE of elements of the stream match the given condition function.
//
// - f:       The matching function to be used.
func (s *Stream) NoneMatch(f func(interface{}) bool) bool {
	return !s.AnyMatch(f)
}

// Contains: Indicates whether the provided value matches any of the values in the stream
//
// - value:   The value to be found.
func (s *Stream) Contains(value interface{}) bool {
	return s.AnyMatch(func(val interface{}) bool {
		return value == val
	})
}

// ForEach: Iterates over all elements in the stream calling the provided function.
func (s *Stream) ForEach(f func(interface{})) {
	array := s.start()

	for i := 0; i < array.Len(); i++ {
		val := array.Index(i)
		f(val)
	}
}

// ParallelForEach: Iterates over all elements in the stream calling the provided function. Creates multiple go channels to parallelize
// the operation. ParallelForeach does not use any thread values previously provided in any filtering method nor enables parallel filtering
// if any filtering is done prior to the `ParallelForEach` phase. Only use `ParallelForEach` if the order in which the elements are processed
// does not matter, otherwise see `ForEach`.
//
// - threads:   Indicates the amount of go channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates
//              the maximum amount of available CPUs will be the number that determines the amount of go channels to be used.
// - skipWait:  Indicates whether `ParallelForEach` will wait until all channels are done processing.
func (s *Stream) ParallelForEach(f func(interface{}), threads int, skipWait ...bool) {
	var wg sync.WaitGroup
	cores := getCores(threads)
	array := s.start()

	if array.Len() < cores {
		cores = array.Len()
	}

	worker := func(start, end int) {
		defer wg.Done()
		for i := start; i < end && i < array.Len(); i++ {
			f(array.Index(i))
		}
	}

	sliceSize := int(math.Ceil(float64(array.Len()) / float64(cores)))

	wg.Add(cores)

	for i := 0; i < cores; i++ {
		go worker(i*sliceSize, (i+1)*sliceSize)
	}

	if len(skipWait) == 0 || !skipWait[0] {
		wg.Wait()
	}
}

// ToArray: Converts the resulting stream back to an collection
func (s *Stream) ToArray() interface{} {
	return s.start().ToArray()
}

// OrderBy: Sorts the elements ascending in the stream using the provided comparable function.
//
// - desc:  indicates whether the sorting should be done descendant
func (s *Stream) OrderBy(f func(interface{}, interface{}) int, desc ...bool) *Stream {
	s.sorts = nil
	return s.ThenBy(f, desc...)
}

// ThenBy: If two elements are considered equal after previously applying a comparable function,
// attempts to sort ascending the 2 elements with an additional comparable function.
//
// - desc:  indicates whether the sorting should be done descendant
func (s *Stream) ThenBy(f func(interface{}, interface{}) int, desc ...bool) *Stream {
	d := false

	if len(desc) > 0 {
		d = desc[0]
	}

	s.sorts = append(s.sorts, sortFunc{
		fn:   f,
		desc: d,
	})
	return s
}

func (s *Stream) start() ICollection {
	if s.threads != 1 {
		return s.parallelStart(s.threads)
	}

	var array = s.collection
	array = s.filter(array)
	array = s.except(array)
	array = s.sort(array)
	return array
}

func (s *Stream) parallelStart(threads int) ICollection {
	var array = s.collection
	array = s.parallelStartHandler(array, threads)
	array = s.sort(array)
	return array
}

func (s *Stream) filter(array ICollection) ICollection {
	return filterHandler(array, 0, array.Len(), s.filters, false)
}

func (s *Stream) except(array ICollection) ICollection {
	return filterHandler(array, 0, array.Len(), s.exceptions, true)
}

func (s *Stream) parallelStartHandler(array ICollection, threads int) ICollection {
	worker := func(result chan ICollection, start, end int) {
		slice := filterHandler(array, start, end, s.filters, false)
		result <- filterHandler(slice, 0, slice.Len(), s.exceptions, true)
	}

	return parallelHandler(array, threads, worker)
}

func (s *Stream) sort(array ICollection) ICollection {
	if len(s.sorts) == 0 {
		return array
	}

	so := sorter{
		collection: array,
		sorts:      s.sorts,
	}

	sort.Slice(so.collection.ToArray(), so.makeLessFunc())
	return so.collection
}

func (s *Stream) updateCores(threads ...int) int {
	if len(threads) > 0 {
		s.threads = getCores(threads...)
	}
	return s.threads
}

func (s *sorter) makeLessFunc() func(i, j int) bool {
	return func(x, y int) bool {
		val := 0

		for i := 0; val == 0 && i < len(s.sorts); i++ {
			sorter := s.sorts[i]
			val = sorter.fn(s.collection.Index(x), s.collection.Index(y))

			if sorter.desc {
				val = val * -1
			}
		}

		return val < 0
	}
}

func filterHandler(col ICollection, start, end int, filters []func(interface{}) bool, negate bool) ICollection {
	if len(filters) == 0 {
		return col
	}

	ret := NewCollection(col.ElementType())

	for i := start; i < end; i++ {
		x := col.Index(i)
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
			ret.Add(x)
		}
	}

	return ret
}

func anyMatch(col ICollection, start, end int, filters []func(interface{}) bool, negate bool) bool {
	if len(filters) == 0 {
		return false
	}

	for i := start; i < end; i++ {
		x := col.Index(i)
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
			return true
		}
	}

	return false
}

func parallelHandler(col ICollection, threads int, worker func(result chan ICollection, start, end int)) ICollection {
	var ret ICollection = NewCollection(col.ElementType())
	cores := getCores(threads)

	if col.Len() < cores {
		cores = col.Len()
	}

	sliceSize := int(math.Ceil(float64(col.Len()) / float64(cores)))
	c := make(chan ICollection, cores)

	for i := 0; i < cores; i++ {
		go worker(c, i*sliceSize, (i+1)*sliceSize)
	}

	for i := 0; i < cores; i++ {
		ret.AddCollection(<-c)
	}

	return ret
}

func getCores(threads ...int) int {
	if len(threads) == 0 {
		return 1
	}

	maxCores := runtime.NumCPU()

	if maxCores < threads[0] || threads[0] <= 0 {
		return maxCores
	} else {
		return threads[0]
	}
}

// STREAM
//   Distinct
//   Reverse

// OPTIONAL ?? or element
//    Min
//    Max
//    Average
//    FindAny                  For parallel operations. Post MVP

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
//   OrderBy           DONE
//   OrderByDescending DONE
//   ThenBy            DONE
//   ThenByDescending  DONE
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
//    FindFirst, First   DONE
//    ElementAt, At      DONE
//    ElementAtOrDefault DONE
//    AtOrDefault        DONE
//    Last               DONE
//    LastOrDefault      DONE
//
// VOID
//    ForEach           DONE
//
// ARRAY
//    ToArray           DONE
