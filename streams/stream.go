package streams

import (
	"math"
	"reflect"
	"runtime"
	"sort"
	"sync"
)

type Stream struct {
	iterable IIterable
	filters  []func(interface{}) bool
	sorts    []sortFunc
	threads  int
}

type sortFunc struct {
	fn   func(interface{}, interface{}) int
	desc bool
}

type sorter struct {
	array interface{}
	sorts []sortFunc
}

// From Creates a Stream from a given iterable or ICollection.
//
// - iterable: The iterable or ICollection to be used to create the stream
// - threads:    If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//               to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//               available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//               best combine it with a `SortBy`. Only needs to be provided once per stream.
func From(collection interface{}, threads ...int) *Stream {
	colReflectType := reflect.TypeOf((*IIterable)(nil)).Elem()

	if reflect.PtrTo(reflect.TypeOf(collection)).Implements(colReflectType) {
		return FromIterable(collection.(IIterable), threads...)
	}

	return FromArray(collection, threads...)
}

// FromArray Creates a Stream from a given array.
//
// - iterable:   The iterable to be used to create the stream
// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//            best combine it with a `SortBy`. Only needs to be provided once per stream.
func FromArray(array interface{}, threads ...int) *Stream {
	return FromIterable(NewCollectionFromArray(array), threads...)
}

// FromIterable Creates a Stream from a given IIterable.
//
// - iterable: The ICollection to be used to create the stream
// - threads:    If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//               to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//               available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//               best combine it with a `SortBy`. Only needs to be provided once per stream.
func FromIterable(iterable IIterable, threads ...int) *Stream {
	return &Stream{
		iterable: iterable,
		threads:  getCores(threads...),
	}
}

// SetThreads Sets the amount of go channels to be used for parallel filtering to a maximum of the available CPUs in the
// host machine. Providing a value <= 0, indicates the maximum amount of available CPUs will be the number that determines
// the amount of go channels to be used. If order matters, best combine it with a `SortBy`. Only needs to be provided once
// per stream.
func (s *Stream) SetThreads(threads int) int {
	return s.updateCores(threads)
}

// Filter Filters any element that does not meet the condition provided by the function.
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

// Except Filters all elements that meet the condition provided by the function.
//
// - f:       The filtering function to be used.
// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//            best combine it with a `SortBy`. Only needs to be provided once per stream.
func (s *Stream) Except(f func(interface{}) bool, threads ...int) *Stream {
	s.updateCores(threads...)
	s.filters = append(s.filters, func(x interface{}) bool { return !f(x) })
	return s
}

// Map Maps the elements of the iterable to a new element, using the mapping function provided
//
// - f:       The filtering function to be used.
// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
//            best combine it with a `SortBy`. Only needs to be provided once per stream.
func (s *Stream) Map(f func(interface{}) interface{}, threads ...int) *Stream {
	iterable := s.process()
	var col ICollection

	iterator := iterable.Iterator()
	for old := iterator.Current(); iterator.HasNext(); old = iterator.Next() {
		n := f(old)

		if col == nil {
			col = NewCollection(reflect.TypeOf(n))
		}

		col.Add(n)
	}

	return FromIterable(col)
}

// First Returns the first element of the resulting stream.
// Returns nil (or default value if provided) if the resulting stream is empty.
func (s *Stream) First(defaultValue ...interface{}) interface{} {
	return s.At(0, defaultValue...)
}

// Last Returns the last element of the resulting stream.
// Returns nil (or default value if provided) if the resulting stream is empty.
func (s *Stream) Last(defaultValue ...interface{}) interface{} {
	return s.AtReverse(0, defaultValue...)
}

// At Returns the element at the given index in the resulting stream.
// Returns nil (or default value if provided) if out of bounds.
func (s *Stream) At(index int, defaultValue ...interface{}) interface{} {
	iterator := s.process().Iterator()
	iterator.Skip(index)

	val := iterator.Current()
	if val == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

// AtReverse Returns the element at the given position, starting from the last element to the first in the resulting stream.
// Returns nil (or default value if provided) if out of bounds.
func (s *Stream) AtReverse(pos int, defaultValue ...interface{}) interface{} {
	// TODO: Return error if Len is unavailable
	iterable := s.process()
	iterator := iterable.Iterator()

	i := iterable.Len() - 1 - pos

	if i >= 0 {
		iterator.Skip(i)
		return iterator.Current()
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

// Count Counts the elements of the resulting stream
func (s *Stream) Count() int {
	iterable := s.process()

	if iterable.Len() >= 0 {
		return iterable.Len()
	}

	iterator := iterable.Iterator()
	size := 0

	for ; iterator.HasNext(); iterator.Next() {
		size++
	}

	return size
}

// AnyMatch Indicates whether any elements of the stream match the given condition function.
//
// - f:       The matching function to be used.
func (s *Stream) AnyMatch(f func(interface{}) bool) bool {
	iterable := s.process()
	return anyMatch(iterable, 0, iterable.Len(), f, false)
}

// AllMatch Indicates whether ALL elements of the stream match the given condition function
//
// - f:       The matching function to be used.
func (s *Stream) AllMatch(f func(interface{}) bool) bool {
	iterable := s.process()
	return !anyMatch(iterable, 0, iterable.Len(), f, true)
}

// NoneMatch Indicates whether NONE of elements of the stream match the given condition function.
//
// - f:       The matching function to be used.
func (s *Stream) NoneMatch(f func(interface{}) bool) bool {
	return !s.AnyMatch(f)
}

// Contains Indicates whether the provided value matches any of the values in the stream
//
// - value:   The value to be found.
func (s *Stream) Contains(value interface{}) bool {
	return s.AnyMatch(func(val interface{}) bool {
		return value == val
	})
}

// ForEach Iterates over all elements in the stream calling the provided function.
func (s *Stream) ForEach(f func(interface{})) {
	iterable := s.process()
	iterator := iterable.Iterator()

	for val := iterator.Current(); iterator.HasNext(); val = iterator.Next() {
		f(val)
	}
}

// ParallelForEach Iterates over all elements in the stream calling the provided function. Creates multiple go channels to parallelize
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
	iterable := s.process()

	if iterable.Len() < cores {
		cores = iterable.Len()
	}

	worker := func(start, end int) {
		defer wg.Done()
		iterator := iterable.Iterator()
		iterator.Skip(start)
		i := start

		for val := iterator.Current(); iterator.HasNext() && i < end; val = iterator.Next() {
			i++
			f(val)
		}
	}

	sliceSize := int(math.Ceil(float64(iterable.Len()) / float64(cores)))

	wg.Add(cores)

	for i := 0; i < cores; i++ {
		go worker(i*sliceSize, (i+1)*sliceSize)
	}

	if len(skipWait) == 0 || !skipWait[0] {
		wg.Wait()
	}
}

// ToArray Returns an array of elements from the resulting stream
func (s *Stream) ToArray() interface{} {
	return s.process().ToArray()
}

// ToCollection Returns a `ICollection` of elements from the resulting stream
func (s *Stream) ToCollection() ICollection {
	iterable := s.process()
	colReflectType := reflect.TypeOf((*ICollection)(nil)).Elem()

	if reflect.PtrTo(reflect.TypeOf(iterable)).Implements(colReflectType) {
		return iterable.(ICollection)
	}

	ret := NewCollection(iterable.ElementType())
	ret.AddAll(iterable)

	return ret
}

// ToIterable Returns a `IIterable` of elements from the resulting stream
func (s *Stream) ToIterable() IIterable {
	return s.process()
}

// OrderBy Sorts the elements in the stream using the provided comparable function.
//
// - desc:  indicates whether the sorting should be done descendant
func (s *Stream) OrderBy(f func(interface{}, interface{}) int, desc ...bool) *Stream {
	s.sorts = nil
	return s.ThenBy(f, desc...)
}

// ThenBy If two elements are considered equal after previously applying a comparable function,
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

func (s *Stream) process() IIterable {
	if s.threads != 1 {
		return s.parallelProcess(s.threads)
	}

	var iterable = s.iterable
	iterable = s.filter(iterable)
	iterable = s.sort(iterable)
	return iterable
}

func (s *Stream) parallelProcess(threads int) IIterable {
	var iterable = s.iterable
	iterable = s.parallelProcessHandler(iterable, threads)
	iterable = s.sort(iterable)
	return iterable
}

func (s *Stream) filter(iterable IIterable) IIterable {
	return s.filterHandler(iterable, 0, iterable.Len())
}

func (s *Stream) filterHandler(iterable IIterable, start, end int) IIterable {
	if len(s.filters) == 0 {
		return iterable
	}

	ret := NewCollection(iterable.ElementType())
	iterator := iterable.Iterator().Skip(start)
	i := start

	for x := iterator.Current(); iterator.HasNext() && i < end; x = iterator.Next() {
		i++
		var match bool = true

		for _, f := range s.filters {
			match = match && f(x)

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

func (s *Stream) parallelProcessHandler(iterable IIterable, threads int) IIterable {
	worker := func(result chan IIterable, start, end int) {
		result <- s.filterHandler(iterable, start, end)
	}

	ret := NewCollection(iterable.ElementType())
	cores := getCores(threads)

	if iterable.Len() < cores {
		cores = iterable.Len()
	}

	sliceSize := int(math.Ceil(float64(iterable.Len()) / float64(cores)))
	c := make(chan IIterable, cores)

	for i := 0; i < cores; i++ {
		go worker(c, i*sliceSize, (i+1)*sliceSize)
	}

	for i := 0; i < cores; i++ {
		ret.AddAll(<-c)
	}

	return ret
}

func (s *Stream) sort(iterable IIterable) IIterable {
	if len(s.sorts) == 0 {
		return iterable
	}

	so := sorter{
		array: iterable.ToArray(),
		sorts: s.sorts,
	}

	sort.Slice(so.array, so.makeLessFunc())
	return NewCollectionFromArray(so.array)
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
			arr := reflect.ValueOf(s.array)
			val = sorter.fn(arr.Index(x).Interface(), arr.Index(y).Interface())

			if sorter.desc {
				val = val * -1
			}
		}

		return val < 0
	}
}

func anyMatch(iterable IIterable, start, end int, f func(interface{}) bool, negate bool) bool {
	iterator := iterable.Iterator().Skip(start)
	i := start

	for x := iterator.Current(); iterator.HasNext() && i < end; x = iterator.Next() {
		var match bool = true

		if negate {
			match = match && !f(x)
		} else {
			match = match && f(x)
		}

		if match {
			return true
		}
	}

	return false
}

func getCores(threads ...int) int {
	if len(threads) == 0 {
		return 1
	}

	maxCores := runtime.NumCPU()

	if maxCores < threads[0] || threads[0] <= 0 {
		return maxCores
	}
	return threads[0]
}

// TODO:
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
