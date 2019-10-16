package streams

import (
	"math"
	"reflect"
	"runtime"
	"sort"
	"sync"
)

// Stream is the default stream implementation which allows stream operations on IIterables.
type Stream struct {
	iterable IIterable
	filters  []ConditionalFunc
	sorts    []sortFunc
	distinct bool
	threads  int
}

type sortFunc struct {
	fn   SortFunc
	desc bool
}

type sorter struct {
	array interface{}
	sorts []sortFunc
}

type iAdd interface {
	Add(item interface{}) error
}

type mapAdd struct {
	m reflect.Value
}

func (m *mapAdd) Add(item interface{}) error {
	m.m.SetMapIndex(reflect.ValueOf(item), reflect.ValueOf(true))
	return nil
}

// SetThreads Sets the amount of go channels to be used for parallel filtering to a maximum of the available CPUs in the
// host machine. Providing a value <= 0, indicates the maximum amount of available CPUs will be the number that determines
// the amount of go channels to be used. If order matters, best combine it with a `SortBy`. Only needs to be provided once
// per stream.
//
// - threads: The amount of threads to use
//
func (s *Stream) SetThreads(threads int) int {
	return s.updateCores(threads)
}

// Filter Filters any element that does not meet the condition provided by the function.
//
// - f:       The filtering function to be used.
// - threads: (Optional) If provided, enables parallel filtering for all filter operations. Indicates the amount of go
//            channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum
//            amount of available CPUs will be the number that determines the amount of go channels to be used. If order
//            matters, best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func (s *Stream) Filter(f ConditionalFunc, threads ...int) IStream {
	s.updateCores(threads...)
	s.filters = append(s.filters, f)
	return s
}

// Except Filters all elements that meet the condition provided by the function.
//
// - f:       The filtering function to be used.
// - threads: (Optional) If provided, enables parallel filtering for all filter operations. Indicates the amount of go
//            channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum
//            amount of available CPUs will be the number that determines the amount of go channels to be used. If order
//            matters, best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func (s *Stream) Except(f ConditionalFunc, threads ...int) IStream {
	s.updateCores(threads...)
	s.filters = append(s.filters, func(x interface{}) bool { return !f(x) })
	return s
}

// Map Maps the elements of the iterable to a new element, using the mapping function provided
//
// - f:       The conversion function to use.
// - threads: (Optional) If provided, enables parallel filtering for all filter operations. Indicates the amount of go
//            channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum
//            amount of available CPUs will be the number that determines the amount of go channels to be used. If order
//            matters, best combine it with a `SortBy`. Only needs to be provided once per stream.
//
func (s *Stream) Map(f ConvertFunc, threads ...int) IStream {
	iterable := s.process()
	var col ICollection

	iterator := iterable.Iterator()
	for old := iterator.Current(); iterator.HasNext(); old = iterator.Next() {
		n := f(old)

		if col == nil {
			col = NewArrayCollection(reflect.TypeOf(n))
		}

		_ = col.Add(n)
	}

	return FromIterable(col)
}

// Distinct Returns a stream consisting of the distinct elements
func (s *Stream) Distinct() IStream {
	s.distinct = true
	return s
}

// First Returns the first element of the resulting stream.
// Returns nil (or default value if provided) if the resulting stream is empty.
//
// - defaultValue:  (Optional) The default value to return if empty.
//
func (s *Stream) First(defaultValue ...interface{}) interface{} {
	return s.At(0, defaultValue...)
}

// Last Returns the last element of the resulting stream.
// Returns nil (or default value if provided) if the resulting stream is empty.
//
// - defaultValue:  (Optional) The default value to return if empty.
//
func (s *Stream) Last(defaultValue ...interface{}) interface{} {
	return s.AtReverse(0, defaultValue...)
}

// At Returns the element at the given index in the resulting stream.
// Returns nil (or default value if provided) if out of bounds.
//
// - index:         The index of the element to return
// - defaultValue:  (Optional) The default value to return if out of bounds.
//
func (s *Stream) At(index int, defaultValue ...interface{}) interface{} {
	iterable := s.process()
	if iterable == nil {
		return nil
	}
	iterator := iterable.Iterator()
	iterator.Skip(index)

	val := iterator.Current()
	if val == nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

// AtReverse Returns the element at the given position, starting from the last element to the first in the resulting stream.
// Returns nil (or default value if provided) if out of bounds.
//
// - post:          The position of the element to return from the last element.
// - defaultValue:  (Optional) The default value to return if out of bounds.
//
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
//
func (s *Stream) AnyMatch(f ConditionalFunc) bool {
	iterable := s.process()
	return anyMatch(iterable, 0, iterable.Len(), f, false)
}

// AllMatch Indicates whether ALL elements of the stream match the given condition function
//
// - f:       The matching function to be used.
//
func (s *Stream) AllMatch(f ConditionalFunc) bool {
	iterable := s.process()
	return !anyMatch(iterable, 0, iterable.Len(), f, true)
}

// IfAllMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if `AllMatch`
// based on what the result of `AllMatch` would be with the provided conditional function
//
// - f:       The matching function to be used.
func (s *Stream) IfAllMatch(f ConditionalFunc) IThen {
	return &thenWrapper{
		conditionMet: s.AllMatch(f),
	}
}

// NotAllMatch is the negation of `AllMatch`. If any of the elements don not match the provided condition
// the result will be `true`; `false` otherwise.
//
// - f:       The matching function to be used.
func (s *Stream) NotAllMatch(f ConditionalFunc) bool {
	return !s.AllMatch(f)
}

// IfNotAllMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if `AllMatch`
// based on what the result of `AllMatch` would be with the provided conditional function
//
// - f:       The matching function to be used.
func (s *Stream) IfNotAllMatch(f ConditionalFunc) IThen {
	return &thenWrapper{
		conditionMet: s.NotAllMatch(f),
	}
}

// NoneMatch Indicates whether NONE of elements of the stream match the given condition function.
//
// - f:       The matching function to be used.
//
func (s *Stream) NoneMatch(f ConditionalFunc) bool {
	return !s.AnyMatch(f)
}

// Contains Indicates whether the provided value matches any of the values in the stream
//
// - value:   The value to be found.
//
func (s *Stream) Contains(value interface{}) bool {
	return s.AnyMatch(func(val interface{}) bool {
		return value == val
	})
}

// ForEach Iterates over all elements in the stream calling the provided function.
//
// - f:       The iterator function to be used.
//
func (s *Stream) ForEach(f IterFunc) {
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
// - f:         The iterator function to be used.
// - threads:   Indicates the amount of go channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates
//              the maximum amount of available CPUs will be the number that determines the amount of go channels to be used.
// - skipWait:  Indicates whether `ParallelForEach` will wait until all channels are done processing.
//
func (s *Stream) ParallelForEach(f IterFunc, threads int, skipWait ...bool) {
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
func (s *Stream) ToArray(defaultArray ...interface{}) interface{} {
	iterable := s.process()
	if iterable == nil {
		return nil
	}
	return iterable.ToArray(defaultArray...)
}

// ToCollection Returns a `ICollection` of elements from the resulting stream
func (s *Stream) ToCollection() ICollection {
	iterable := s.process()
	colReflectType := reflect.TypeOf((*ICollection)(nil)).Elem()

	if reflect.PtrTo(reflect.TypeOf(iterable)).Implements(colReflectType) {
		return iterable.(ICollection)
	}

	ret := NewArrayCollection(iterable.ElementType())
	_ = ret.AddAll(iterable)

	return ret
}

// ToIterable Returns a `IIterable` of elements from the resulting stream
func (s *Stream) ToIterable() IIterable {
	return s.process()
}

// OrderBy Sorts the elements in the stream using the provided comparable function.
//
// - f:     The sorting function to be used
// - desc:  Indicates whether the sorting should be done descendant
//
func (s *Stream) OrderBy(f SortFunc, desc ...bool) IStream {
	s.sorts = nil
	return s.ThenBy(f, desc...)
}

// ThenBy If two elements are considered equal after previously applying a comparable function,
// attempts to sort ascending the 2 elements with an additional comparable function.
//
// - f:     The sorting function to be used
// - desc:  Indicates whether the sorting should be done descendant
//
func (s *Stream) ThenBy(f SortFunc, desc ...bool) IStream {
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

	iterable := s.iterable
	if iterable == nil {
		return nil
	}
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
	return s.iterHandler(iterable, 0, iterable.Len())
}

func (s *Stream) iterHandler(iterable IIterable, start, end int) IIterable {
	if len(s.filters) == 0 && !s.distinct {
		return iterable
	}

	var adder iAdd

	ret := NewArrayCollection(iterable.ElementType())
	iterator := iterable.Iterator().Skip(start)
	i := start
	adder = ret

	if s.distinct {
		mType := reflect.MapOf(iterable.ElementType(), reflect.TypeOf(true))
		adder = &mapAdd{m: reflect.MakeMap(mType)}
	}

	for x := iterator.Current(); iterator.HasNext() && i < end; x = iterator.Next() {
		i++
		match := true

		for _, f := range s.filters {
			match = match && f(x)

			if !match {
				break
			}
		}

		if match {
			_ = adder.Add(x)
		}
	}

	if s.distinct {
		mAdder := adder.(*mapAdd)
		for _, v := range mAdder.m.MapKeys() {
			_ = ret.Add(v.Interface())
		}
	}

	return ret
}

func (s *Stream) parallelProcessHandler(iterable IIterable, threads int) IIterable {
	worker := func(result chan IIterable, start, end int) {
		result <- s.iterHandler(iterable, start, end)
	}

	ret := NewArrayCollection(iterable.ElementType())
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
	v, _ := NewCollectionFromArray(so.array)
	return v
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

func anyMatch(iterable IIterable, start, end int, f ConditionalFunc, negate bool) bool {
	iterator := iterable.Iterator().Skip(start)
	i := start

	for x := iterator.Current(); iterator.HasNext() && i < end; x = iterator.Next() {
		match := true

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
//
// STREAM
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
