package streams

// Reduce folds all elements of the source into a single value using the provided accumulator function.
// The identity value is used as the initial accumulator value.
//
//	{source}    -  The source to read elements from. Accepts []T, IIterable[T], IStream[T].
//	{identity}  -  The initial value for the accumulator.
//	{f}         -  The accumulator function: func(accumulator, element) accumulator.
func Reduce[T comparable](source any, identity T, f func(T, T) T) T {
	result := identity
	iter := resolveIterator[T](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		result = f(result, val)
	}
	return result
}

// ReduceAny is like Reduce but allows the accumulator to be a different type than the elements.
//
//	{source}    -  The source to read elements from. Accepts []T, IIterable[T], IStream[T].
//	{identity}  -  The initial value for the accumulator.
//	{f}         -  The accumulator function: func(accumulator, element) accumulator.
func ReduceAny[T comparable, R any](source any, identity R, f func(R, T) R) R {
	result := identity
	iter := resolveIterator[T](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		result = f(result, val)
	}
	return result
}

// FlatMap maps each element to a slice and flattens the results into a single IList.
//
//	{source}  -  The source to read elements from. Accepts []From, IIterable[From], IStream[From].
//	{f}       -  A function that maps each element to a slice of To.
func FlatMap[From, To comparable](source any, f func(From) []To) IList[To] {
	var result []To
	iter := resolveIterator[From](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		result = append(result, f(val)...)
	}
	return NewList[To](result)
}

// Concat concatenates multiple sources into a single IStream.
// Each source can be []T, ICollection[T], or IStream[T].
func Concat[T comparable](sources ...any) IStream[T] {
	var result []T
	for _, source := range sources {
		iter := resolveIterator[T](source)
		for val := iter.Current(); iter.HasNext(); val = iter.Next() {
			result = append(result, val)
		}
	}
	return FromArray[T](result)
}

// Zip combines two sources element-wise using the provided function. Stops at the shorter source.
func Zip[A, B, R comparable](sourceA any, sourceB any, f func(A, B) R) IList[R] {
	iterA := resolveIterator[A](sourceA)
	iterB := resolveIterator[B](sourceB)
	var result []R
	for valA, valB := iterA.Current(), iterB.Current(); iterA.HasNext() && iterB.HasNext(); valA, valB = iterA.Next(), iterB.Next() {
		result = append(result, f(valA, valB))
	}
	return NewList[R](result)
}

// GroupBy groups the elements of the source by a key function, returning a map of key to slices.
func GroupBy[T comparable, K comparable](source any, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	iter := resolveIterator[T](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		key := keyFn(val)
		result[key] = append(result[key], val)
	}
	return result
}

// DistinctBy returns a new IList with duplicates removed based on a key function.
func DistinctBy[T comparable, K comparable](source any, keyFn func(T) K) IList[T] {
	seen := make(map[K]struct{})
	var result []T
	iter := resolveIterator[T](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		key := keyFn(val)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, val)
		}
	}
	return NewList[T](result)
}

// Min returns the minimum element from the source using the provided less function.
// Returns the zero value of T and false if the source is empty.
func Min[T comparable](source any, less func(T, T) bool) (T, bool) {
	iter := resolveIterator[T](source)
	if !iter.HasNext() {
		var zero T
		return zero, false
	}
	first := true
	var result T
	iter.ForEachRemaining(func(val T) {
		if first || less(val, result) {
			result = val
			first = false
		}
	})
	return result, true
}

// Max returns the maximum element from the source using the provided less function.
// Returns the zero value of T and false if the source is empty.
func Max[T comparable](source any, less func(T, T) bool) (T, bool) {
	iter := resolveIterator[T](source)
	if !iter.HasNext() {
		var zero T
		return zero, false
	}
	first := true
	var result T
	iter.ForEachRemaining(func(val T) {
		if first || less(result, val) {
			result = val
			first = false
		}
	})
	return result, true
}

// Sum returns the sum of all elements in the source. T must be a numeric type satisfying INumeric.
func Sum[T INumeric](source any) T {
	var result T
	iter := resolveIterator[T](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		result += val
	}
	return result
}

// INumeric comprises the numeric types that support arithmetic operations.
type INumeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Average returns the average of all numeric elements in the source as a float64.
// Returns 0 and false if the source is empty.
func Average[T INumeric](source any) (float64, bool) {
	var sum float64
	count := 0
	iter := resolveIterator[T](source)
	for val := iter.Current(); iter.HasNext(); val = iter.Next() {
		sum += toFloat64(val)
		count++
	}
	if count == 0 {
		return 0, false
	}
	return sum / float64(count), true
}

func toFloat64[T INumeric](val T) float64 {
	switch v := any(val).(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	}
	return 0
}

// Range creates a stream of integers from start (inclusive) to end (exclusive).
func Range(start, end int) IStream[int] {
	if end <= start {
		return FromArray[int](nil)
	}
	arr := make([]int, end-start)
	for i := range arr {
		arr[i] = start + i
	}
	return FromArray[int](arr)
}

// Repeat creates a stream that contains the given value repeated n times.
func Repeat[T comparable](value T, n int) IStream[T] {
	if n <= 0 {
		return FromCollection[T](NewList[T]())
	}
	arr := make([]T, n)
	for i := range arr {
		arr[i] = value
	}
	return FromArray[T](arr)
}

// Union returns a new IList containing the distinct union of elements from both sources.
func Union[T comparable](sourceA, sourceB any) IList[T] {
	set := NewSet[T]()
	var result []T
	addUnique := func(iter IIterator[T]) {
		for val := iter.Current(); iter.HasNext(); val = iter.Next() {
			if !set.Contains(val) {
				set.Add(val)
				result = append(result, val)
			}
		}
	}
	addUnique(resolveIterator[T](sourceA))
	addUnique(resolveIterator[T](sourceB))
	return NewList[T](result)
}

// Intersect returns a new IList containing elements present in both sources.
func Intersect[T comparable](sourceA, sourceB any) IList[T] {
	setB := NewSet[T]()
	iterB := resolveIterator[T](sourceB)
	for val := iterB.Current(); iterB.HasNext(); val = iterB.Next() {
		setB.Add(val)
	}

	seen := NewSet[T]()
	var result []T
	iterA := resolveIterator[T](sourceA)
	for val := iterA.Current(); iterA.HasNext(); val = iterA.Next() {
		if setB.Contains(val) && !seen.Contains(val) {
			seen.Add(val)
			result = append(result, val)
		}
	}
	return NewList[T](result)
}

func resolveIterator[T comparable](source any) IIterator[T] {
	switch src := source.(type) {
	case []T:
		return newArrayIterator[T](src)
	case IStream[T]:
		col := src.ToCollection()
		if col == nil {
			return newArrayIterator[T]()
		}
		return col.Iterator()
	case IIterable[T]:
		return src.Iterator()
	}
	panic("invalid source type")
}
