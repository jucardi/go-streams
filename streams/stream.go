package streams

// IStream defines the functions of a stream implementation
type IStream interface {
	// SetThreads Sets the amount of go channels to be used for parallel filtering to a maximum of the available CPUs in the
	// host machine. Providing a value <= 0, indicates the maximum amount of available CPUs will be the number that determines
	// the amount of go channels to be used. If order matters, best combine it with a `SortBy`. Only needs to be provided once
	// per stream.
	SetThreads(threads int) int

	// Filter Filters any element that does not meet the condition provided by the function.
	//
	// - f:       The filtering function to be used.
	// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
	//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
	//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
	//            best combine it with a `SortBy`. Only needs to be provided once per stream.
	Filter(f ConditionalFunc, threads ...int) IStream

	// Except Filters all elements that meet the condition provided by the function.
	//
	// - f:       The filtering function to be used.
	// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
	//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
	//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
	//            best combine it with a `SortBy`. Only needs to be provided once per stream.
	Except(f ConditionalFunc, threads ...int) IStream

	// Map Maps the elements of the iterable to a new element, using the mapping function provided
	//
	// - f:       The filtering function to be used.
	// - threads: If provided, enables parallel filtering for all filter operations. Indicates the amount of go channels
	//            to be used to a maximum of the available CPUs in the host machine. <= 0 indicates the maximum amount of
	//            available CPUs will be the number that determines the amount of go channels to be used. If order matters,
	//            best combine it with a `SortBy`. Only needs to be provided once per stream.
	Map(f ConvertFunc, threads ...int) IStream

	// Distinct Returns a stream consisting of the distinct elements
	Distinct() IStream

	// First Returns the first element of the resulting stream.
	// Returns nil (or default value if provided) if the resulting stream is empty.
	First(defaultValue ...interface{}) interface{}

	// Last Returns the last element of the resulting stream.
	// Returns nil (or default value if provided) if the resulting stream is empty.
	Last(defaultValue ...interface{}) interface{}

	// At Returns the element at the given index in the resulting stream.
	// Returns nil (or default value if provided) if out of bounds.
	At(index int, defaultValue ...interface{}) interface{}

	// AtReverse Returns the element at the given position, starting from the last element to the first in the resulting stream.
	// Returns nil (or default value if provided) if out of bounds.
	AtReverse(pos int, defaultValue ...interface{}) interface{}

	// Count Counts the elements of the resulting stream
	Count() int

	// AnyMatch Indicates whether any elements of the stream match the given condition function.
	//
	// - f:       The matching function to be used.
	AnyMatch(f ConditionalFunc) bool

	// AllMatch Indicates whether ALL elements of the stream match the given condition function
	//
	// - f:       The matching function to be used.
	AllMatch(f ConditionalFunc) bool

	// IfAllMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if `AllMatch`
	// based on what the result of `AllMatch` would be with the provided conditional function
	//
	// - f:       The matching function to be used.
	IfAllMatch(f ConditionalFunc) IThen

	// NotAllMatch is the negation of `AllMatch`. If any of the elements don not match the provided condition
	// the result will be `true`; `false` otherwise.
	//
	// - f:       The matching function to be used.
	NotAllMatch(f ConditionalFunc) bool

	// IfNotAllMatch returns a `Then` handler where actions like `Then` or `Else` can be triggered if `AllMatch`
	// based on what the result of `AllMatch` would be with the provided conditional function
	//
	// - f:       The matching function to be used.
	IfNotAllMatch(f ConditionalFunc) IThen

	// NoneMatch Indicates whether NONE of elements of the stream match the given condition function.
	//
	// - f:       The matching function to be used.
	NoneMatch(f ConditionalFunc) bool

	// Contains Indicates whether the provided value matches any of the values in the stream
	//
	// - value:   The value to be found.
	Contains(value interface{}) bool

	// ForEach Iterates over all elements in the stream calling the provided function.
	ForEach(f IterFunc)

	// ParallelForEach Iterates over all elements in the stream calling the provided function. Creates multiple go channels to parallelize
	// the operation. ParallelForeach does not use any thread values previously provided in any filtering method nor enables parallel filtering
	// if any filtering is done prior to the `ParallelForEach` phase. Only use `ParallelForEach` if the order in which the elements are processed
	// does not matter, otherwise see `ForEach`.
	//
	// - threads:   Indicates the amount of go channels to be used to a maximum of the available CPUs in the host machine. <= 0 indicates
	//              the maximum amount of available CPUs will be the number that determines the amount of go channels to be used.
	// - skipWait:  Indicates whether `ParallelForEach` will wait until all channels are done processing.
	ParallelForEach(f IterFunc, threads int, skipWait ...bool)

	// ToArray Returns an array of elements from the resulting stream
	//
	// - defaultArray:  (optional) an array instance to return in case that after a stream operation
	//                  would result in an empty array.
	ToArray(defaultArray ...interface{}) interface{}

	// ToCollection Returns a `ICollection` of elements from the resulting stream
	ToCollection() ICollection

	// ToIterable Returns a `IIterable` of elements from the resulting stream
	ToIterable() IIterable

	// OrderBy Sorts the elements in the stream using the provided comparable function.
	//
	// - desc:  indicates whether the sorting should be done descendant
	OrderBy(f SortFunc, desc ...bool) IStream

	// ThenBy If two elements are considered equal after previously applying a comparable function,
	// attempts to sort ascending the 2 elements with an additional comparable function.
	//
	// - desc:  indicates whether the sorting should be done descendant
	ThenBy(f SortFunc, desc ...bool) IStream
}
