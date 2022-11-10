package streams

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

var (
	// To ensure *Stream implements IStream on build
	_ IStream[string] = (*Stream[string])(nil)
)

// Stream is the default stream implementation which allows stream operations on IIterables.
type Stream[T comparable] struct {
	iterable ICollection[T]
	filters  []ConditionalFunc[T]
	sorts    []sortFunc[T]
	distinct bool
	threads  int

	current ICollection[T]
}

type sortFunc[T comparable] struct {
	fn   SortFunc[T]
	desc bool
}

type sorter[T comparable] struct {
	array []T
	sorts []sortFunc[T]
}

func (s *Stream[T]) SetThreads(threads int) IStream[T] {
	s.updateCores(threads)
	return s
}

func (s *Stream[T]) Filter(f ConditionalFunc[T]) IStream[T] {
	s.filters = append(s.filters, f)
	return s
}

func (s *Stream[T]) Except(f ConditionalFunc[T]) IStream[T] {
	s.filters = append(s.filters, func(x T) bool { return !f(x) })
	return s
}

func (s *Stream[T]) Sort(f SortFunc[T], desc ...bool) IStream[T] {
	d := false

	if len(desc) > 0 {
		d = desc[0]
	}

	s.sorts = append(s.sorts, sortFunc[T]{
		fn:   f,
		desc: d,
	})
	return s
}

func (s *Stream[T]) Distinct() IStream[T] {
	s.distinct = true
	return s
}

func (s *Stream[T]) First(defaultValue ...T) T {
	return s.At(0, defaultValue...)
}

func (s *Stream[T]) Last(defaultValue ...T) T {
	return s.AtReverse(0, defaultValue...)
}

func (s *Stream[T]) At(index int, defaultValue ...T) (ret T) {
	iterable := s.process()
	if iterable == nil {
		return
	}
	iterator := iterable.Iterator()
	iterator.Skip(index)

	var defaultV T
	ret = iterator.Current()
	if ret == defaultV && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return
}

func (s *Stream[T]) AtReverse(pos int, defaultValue ...T) (ret T) {
	iterable := s.process()
	iterator := iterable.Iterator()

	i := iterable.Len() - 1 - pos

	if i >= 0 {
		iterator.Skip(i)
		ret = iterator.Current()
	}

	var defaultV T
	if ret == defaultV && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return
}

func (s *Stream[T]) Count() int {
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

func (s *Stream[T]) IsEmpty() bool {
	return s.Count() == 0
}

func (s *Stream[T]) Contains(value T) bool {
	return s.AnyMatch(func(val T) bool {
		return value == val
	})
}

func (s *Stream[T]) AnyMatch(f ConditionalFunc[T]) bool {
	iterable := s.process()
	return anyMatch[T](iterable, 0, iterable.Len(), f, false)
}

func (s *Stream[T]) AllMatch(f ConditionalFunc[T]) bool {
	iterable := s.process()
	return !anyMatch[T](iterable, 0, iterable.Len(), f, true)
}

func (s *Stream[T]) NotAllMatch(f ConditionalFunc[T]) bool {
	return !s.AllMatch(f)
}

func (s *Stream[T]) NoneMatch(f ConditionalFunc[T]) bool {
	return !s.AnyMatch(f)
}

func (s *Stream[T]) IfEmpty() IThen[T] {
	return &thenWrapper[T]{
		conditionMet: s.Count() == 0,
		stream:       FromCollection[T](s.current),
	}
}

func (s *Stream[T]) IfAnyMatch(f ConditionalFunc[T]) IThen[T] {
	return &thenWrapper[T]{
		conditionMet: s.AnyMatch(f),
		stream:       FromCollection[T](s.current),
	}
}

func (s *Stream[T]) IfAllMatch(f ConditionalFunc[T]) IThen[T] {
	return &thenWrapper[T]{
		conditionMet: s.AllMatch(f),
		stream:       FromCollection[T](s.current),
	}
}

func (s *Stream[T]) IfNoneMatch(f ConditionalFunc[T]) IThen[T] {
	return &thenWrapper[T]{
		conditionMet: s.NoneMatch(f),
		stream:       FromCollection[T](s.current),
	}
}

func (s *Stream[T]) ForEach(f IterFunc[T]) {
	iterable := s.process()
	iterator := iterable.Iterator()

	iterator.ForEachRemaining(f)
}

func (s *Stream[T]) ParallelForEach(f IterFunc[T], threads int, skipWait ...bool) {
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

func (s *Stream[T]) ToArray() []T {
	iterable := s.process()
	if iterable == nil {
		return nil
	}
	return iterable.ToArray()
}

func (s *Stream[T]) ToCollection() ICollection[T] {
	return s.process()
}

func (s *Stream[T]) ToIterable() IIterable[T] {
	return s.process()
}

func (s *Stream[T]) ToList() IList[T] {
	col := s.ToCollection()
	switch ret := col.(type) {
	case IList[T]:
		return ret
	}
	return NewList[T](col.ToArray())
}

func (s *Stream[T]) ToDistinct() ISet[T] {
	return s.Distinct().ToCollection().(ISet[T])
}

func (s *Stream[T]) process() ICollection[T] {
	if s.threads != 1 {
		return s.parallelProcess(s.threads)
	}

	iterable := s.iterable
	if iterable == nil {
		return nil
	}
	iterable = s.filter(iterable)
	iterable = s.sort(iterable)
	s.current = iterable
	return iterable
}

func (s *Stream[T]) parallelProcess(threads int) ICollection[T] {
	iterable := s.iterable
	iterable = s.parallelProcessHandler(iterable, threads)
	iterable = s.sort(iterable)
	return iterable
}

func (s *Stream[T]) filter(iterable ICollection[T]) ICollection[T] {
	return s.iterHandler(iterable, 0, iterable.Len())
}

func (s *Stream[T]) iterHandler(iterable ICollection[T], start, end int) ICollection[T] {
	if len(s.filters) == 0 && !s.distinct {
		return s.iterable
	}

	var ret ICollection[T]
	iterator := iterable.Iterator().Skip(start)
	i := start

	if s.distinct {
		ret = NewSet[T]()
	} else {
		ret = NewList[T]()
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
			_ = ret.Add(x)
		}
	}

	return ret
}

func (s *Stream[T]) parallelProcessHandler(iterable ICollection[T], threads int) ICollection[T] {
	worker := func(result chan ICollection[T], start, end int) {
		result <- s.iterHandler(iterable, start, end)
	}

	ret := NewList[T]()
	cores := getCores(threads)

	if iterable.Len() < cores {
		cores = iterable.Len()
	}

	sliceSize := int(math.Ceil(float64(iterable.Len()) / float64(cores)))
	c := make(chan ICollection[T], cores)

	for i := 0; i < cores; i++ {
		go worker(c, i*sliceSize, (i+1)*sliceSize)
	}

	for i := 0; i < cores; i++ {
		func(iter ICollection[T]) {
			iter.ForEach(func(item T) { ret.Add(item) })
		}(<-c)
	}

	return ret
}

func (s *Stream[T]) sort(iterable ICollection[T]) ICollection[T] {
	if len(s.sorts) == 0 {
		return iterable
	}

	so := sorter[T]{
		array: iterable.ToArray(),
		sorts: s.sorts,
	}

	sort.Slice(so.array, so.makeLessFunc())
	v := NewList[T](so.array)
	return v
}

func (s *Stream[T]) updateCores(threads ...int) int {
	if len(threads) > 0 {
		s.threads = getCores(threads...)
	}
	return s.threads
}

func (s *sorter[T]) makeLessFunc() func(int, int) bool {
	return func(x, y int) bool {
		val := 0

		for i := 0; val == 0 && i < len(s.sorts); i++ {
			sorter := s.sorts[i]
			val = sorter.fn(s.array[x], s.array[y])

			if sorter.desc {
				val = val * -1
			}
		}

		return val < 0
	}
}

func anyMatch[T comparable](iterable IIterable[T], start, end int, f ConditionalFunc[T], negate bool) bool {
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

	if threads[0] <= 0 {
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
