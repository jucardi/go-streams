## Stream Collections for Go, inspired in Java 8 Streams and .NET Linq

This library provides
- structs to easily represent and manage collections and iterable of elements which size are not necessarily defined.
- Stream structs to support functional-style operations on streams of elements, such as map-reduce transformations on 
collections, filtering, sorting, mapping, foreach parallel operations.

---
***(Nov 2022) Important Update:** This library has been redesigned to support **Golang Generics**, and it is not backwards 
compatible with the previous version. Also requires at least Go 1.18. If you require the older version without generics 
or a version that is compatible with an older version of Go, using Golang Modules you may get the latest stable version
without generics by running the following command:*

```bash
go get github.com/jucardi/go-streams@v1.0.3
```
---

##### Quick Start

To keep up to date with the most recent version:

```bash
go get github.com/jucardi/go-streams
```

Using Golang Modules
```bash
go get github.com/jucardi/go-streams/v2@latest
```

For the latest version without Golang Generics (v1.0.3)
```bash
go get github.com/jucardi/go-streams@v1.0.3
```

##### Quick Overview

Streams facilitate operations on arrays, iterables and collections, such as *filtering*, *sorting*, *mapping*, *foreach*, and parallel operations on the items contained by these arrays, iterables or collections

**A quick example:**

Given the following array

```Go
var fruitArray = []string{"peach", "apple", "pear", "plum", "pineapple", "banana", "kiwi", "orange"}
```

Let's obtain an array of only the elements that start with the letter "p"

```go
fruitsThatStartWithP := streams.

    // Creates a stream from the given array
    From[string](fruitArray).

    // Adds a filter for strings that start with 'p'
    Filter(func(v string) bool {
        return strings.HasPrefix(v, "p")
    }).

    // Sorts alphabetically
    Sort(strings.Compare).

    // Converts back to an array
    ToArray()
```
The resulting array will be `{"peach", "pear", "pineapple", "plum"}`

Here we use an array of string as the source of the stream, perform a filter operation provided by a function that
receives a single element of the collection and determines whether the element should remain in the stream by returning
a boolean.

Now let's do a simple foreach operation

```go
streams.
    From[string](fruitArray).
    Filter(func(v string) bool {
        return strings.HasPrefix(v, "p")
    }).
    ForEach(func(v string) {
        println(v)
    })
```

In this example, once the stream processes the filter, performs a foreach operation with the result. With this operation
we'll obtain the following output in the console

```
peach
pear
plum
pineapple
```

## About the go-streams

The characteristics of a Stream are inspired in the stream features provided by Java 8. The following characteristics 
apply in the go-streams.

- No storage. A stream is not a data structure that stores elements; instead, it conveys elements from a source such as 
a data structure, an array, a generator function, or an I/O channel, through a pipeline of computational operations.
- Functional in nature. An operation on a stream produces a result, but does not modify its source. For example,
filtering a Stream obtained from a collection produces a new Stream without the filtered elements, rather than removing
elements from the source collection.
- Laziness-seeking. Many stream operations, such as filtering and sorting, can be implemented lazily, exposing
opportunities for optimization. For example, "find the first String with three consecutive vowels" need not examine all
the input strings. Stream operations are divided into intermediate (Stream-producing) operations and terminal (value- or
side-effect-producing) operations. Intermediate operations are always lazy.
- Consumable. The elements of a stream are only visited once during the life of a stream. Like an Iterator, a new stream
must be generated to revisit the same elements of the source.

Currently Streams can be obtained in several ways:
- From an array via `streams.FromArray[T comparable](array []T)`
- From an implementation of `streams.ICollection[T comparable]` via `streams.FromCollection[T comparable](col ICollection[T])`
- From a map via `streams.FromMap[K comparable, V any](m map[K]V)`
- Using the generic `streams.From[T comparable](set any)` which accepts arrays or collections

## Stream operations and pipelines

Stream operations are divided into intermediate and terminal operations, and are combined to form stream pipelines. A 
stream pipeline consists of a source (such as an iterable, a collection, an array, a generator function, or an I/O 
channel); followed by zero or more intermediate operations such as `Filter()`, `Except()` or `Sort()` and a terminal
operation such as `ForEach()` or `First()`

Intermediate operations return a stream. They are always lazy; executing an intermediate operation such as `Filter()`
or `Sort()` does not actually perform any action, but instead register the action into the stream that, when traversed,
will execute all filtering and sorting criteria at once. Traversal of the pipeline source does not begin until the
terminal operation of the pipeline is executed, such as `First()`, `Last()`, `ToArray()`, `ForEach()`

Terminal operations, such as `ForEach()`, `First()`, `Last()`, `ParallelForEach()` or `ToArray()`, may traverse the
stream to produce a result or a side-effect. After the terminal operation is performed, the stream pipeline is
considered consumed, and can no longer be used; if you need to traverse the same data source again, you must return to
the data source to get a new stream. In almost all cases, terminal operations are eager, completing their traversal of
the data source and processing of the pipeline before returning

### Intermediate operations

- `Filter(f func(x T) bool)`: Filters any element that does not meet the condition provided by the function.
- `Except(f func(x T) bool)`: Filters all elements that meet the condition provided by the function.
- `Sort(f func(x, y T) int, desc ...bool)`: Sorts the elements in the stream using the provided comparable function.
- `SetThreads(threads int)`: Sets the number of threads to be used when processing the stream.
- `Distinct()`: Ensures that the resulting stream operation only includes unique values.
- `Skip(n int)`: Discards the first `n` elements from the resulting stream and returns a new stream with the remaining
  elements.
- `Limit(n int)`: Returns a new stream containing at most `n` elements from the resulting stream.
- `Reverse()`: Returns a new stream with the elements in reverse order.
- `Peek(f func(x T))`: Performs the provided action on each element without consuming the stream, returning a new stream
  with the same elements. Useful for debugging or logging.
- `TakeWhile(f func(x T) bool)`: Returns a new stream containing elements from the start as long as the condition
  returns true. Stops at the first element that does not match.
- `SkipWhile(f func(x T) bool)`: Skips elements from the start as long as the condition returns true, then returns a
  new stream with all remaining elements.


### Terminal operations

###### Single item returns
- `First(defaultValue ...T)`: Returns the first element of the resulting stream. Returns default T (or defaultValue if
provided) if the resulting stream is empty.
- `Last(defaultValue ...T)`: Returns the last element of the resulting stream. Returns default T (or defaultValue if
provided) if the resulting stream is empty.
- `At(i int, defaultValue ...T)`: Returns the element at the given index in the resulting stream. Returns default T (or
defaultValue if provided) if out of bounds.
- `AtReverse(i int, defaultValue ...T)`: Returns the element at the given position, starting from the last element to
the first in the resulting stream. Returns default T (or defaultValue if provided) if out of bounds.

###### Collection returns
- `ToArray()`: Returns an array of elements from the resulting stream.
- `ToCollection()`: Returns a `ICollection[T]` of elements from the resulting stream.
- `ToIterable()`: Returns a `IIterable[T]` of elements from the resulting stream.
- `ToList()`: Returns a `IList[T]` of elements from the resulting stream.
- `ToDistinct()`: Returns a `ISet[T]` of elements from the resulting stream with only unique values.
- `Chunk(size int)`: Splits the stream into slices of the given size and returns them as a `[][]T`.

###### Boolean returns
- `IsEmpty()`: Indicates whether the resulting stream contains no elements.
- `Contains(value T)`: Indicates whether the provided value matches any of the values in the stream.
- `AnyMatch(f func(x T) bool)`: Indicates whether ANY elements of the stream match the given condition function.
- `AllMatch(f func(x T) bool)`: Indicates whether ALL elements of the stream match the given condition function.
- `NotAllMatch(f func(x T) bool)`: The negation of `AllMatch`. If any of the elements do not match the provided
condition the result will be `true`; `false` otherwise.
- `NoneMatch(f func(x T) bool)`: Indicates whether NONE of elements of the stream match the given condition
function.

###### IThen[T] returns

`IThen[T]` is a handler where functions can be registered and triggered if the stream result meets a certain condition

**E.g:**
```go
streams.
    FromArray[T](array). 
    Filter(filterHandler).
    Filter(anotherFilterHandle).
    Sort(sorterHandler).
    IfEmpty().
    Then(func(resultStream IStream[T]) { 
        // Do something if empty
    }).
    Else(func(resultStream IStream[T]) { 
        // Do something if not empty
    })
```

- `IfEmpty()`: Returns a `IThen[T]` handler if empty.
- `IfAnyMatch(f ConditionalFunc[T])`: Returns a `IThen[T]` handler if any elements match the provided condition.
- `IfAllMatch(f ConditionalFunc[T])`: Returns a `IThen[T]` handler if all elements match the provided condition.
- `IfNoneMatch(f ConditionalFunc[T])`: Returns a `IThen[T]` handler if no elements match the provided condition.

###### Void returns
- `ForEach(f func(x T))`: Iterates over all elements in the stream calling the provided function.
- `ParallelForEach(f func(x T), threads int, skipWait ...bool)`: Iterates over all elements in the stream calling the provided
function. Creates multiple go channels to parallelize the operation. ParallelForeach does not use any thread values 
previously provided in any filtering method nor enables parallel filtering if any filtering is done prior to the 
`ParallelForEach` phase. Only use `ParallelForEach` if the order in which the elements are processed does not matter, 
otherwise see `ForEach`.

###### Int returns
- `Count()`: Counts the elements contained by the resulting stream.

### Standalone functions

These functions operate on sources (arrays, collections, iterators, or streams) and produce results. They are standalone
because Go generics do not allow adding additional type parameters to methods on interfaces/structs.

#### Mapping functions

- `Map[From, To comparable](source any, f ConvertFunc[From, To]) IList[To]`: Maps the elements of the source to new
elements using the mapping function. Outputs a new `IList` with the converted elements. Accepted sources: `[]From`,
`IIterable[From]`, `IIterator[From]`, `IStream[From]`.

- `MapNonComparable[From, To any](source any, f ConvertFunc[From, To]) []To`: Similar to `Map`, but outputs an array
and accepts non-comparable types. Accepted sources: `[]From`, `IIterable[From]`, `IIterator[From]`.

- `MapToPtr[T any](source any) []*T`: Converts a collection of `T` to `[]*T`. Useful for non-comparable structs since
pointers are always comparable and can be used with `IStream`, `ICollection`, etc.

#### Reduce / Aggregate

- `Reduce[T comparable](source any, identity T, f func(T, T) T) T`: Folds all elements into a single value using the
accumulator function, starting from the identity value.

```go
sum := streams.Reduce[int]([]int{1, 2, 3, 4}, 0, func(acc, x int) int { return acc + x })
// sum = 10
```

- `ReduceAny[T comparable, R any](source any, identity R, f func(R, T) R) R`: Like `Reduce` but allows the accumulator
to be a different type than the elements.

```go
csv := streams.ReduceAny[int, string]([]int{1, 2, 3}, "", func(acc string, x int) string {
    if acc != "" { acc += "," }
    return acc + strconv.Itoa(x)
})
// csv = "1,2,3"
```

#### FlatMap

- `FlatMap[From, To comparable](source any, f func(From) []To) IList[To]`: Maps each element to a slice and flattens
the results into a single `IList`.

```go
result := streams.FlatMap[string, string]([]string{"hello", "world"}, func(s string) []string {
    var chars []string
    for _, c := range s { chars = append(chars, string(c)) }
    return chars
})
// result contains ["h", "e", "l", "l", "o", "w", "o", "r", "l", "d"]
```

#### Concat

- `Concat[T comparable](sources ...any) IStream[T]`: Concatenates multiple sources (arrays, collections, or streams)
into a single stream.

```go
result := streams.Concat[int]([]int{1, 2}, []int{3, 4}).ToArray()
// result = [1, 2, 3, 4]
```

#### Zip

- `Zip[A, B, R comparable](sourceA any, sourceB any, f func(A, B) R) IList[R]`: Combines two sources element-wise
using the provided function. Stops at the shorter source.

```go
result := streams.Zip[int, string, string](
    []int{1, 2, 3},
    []string{"a", "b", "c"},
    func(i int, s string) string { return strconv.Itoa(i) + s },
)
// result contains ["1a", "2b", "3c"]
```

#### GroupBy

- `GroupBy[T comparable, K comparable](source any, keyFn func(T) K) map[K][]T`: Groups the elements by a key function,
returning a map of key to slices.

```go
result := streams.GroupBy[string, string](
    []string{"apple", "avocado", "banana", "blueberry"},
    func(s string) string { return string(s[0]) },
)
// result = {"a": ["apple", "avocado"], "b": ["banana", "blueberry"]}
```

#### DistinctBy

- `DistinctBy[T comparable, K comparable](source any, keyFn func(T) K) IList[T]`: Returns a new `IList` with duplicates
removed based on a key function.

```go
result := streams.DistinctBy[string, byte](
    []string{"apple", "avocado", "banana"},
    func(s string) byte { return s[0] },
)
// result contains ["apple", "banana"] (first occurrence per key)
```

#### Aggregation functions

- `Min[T comparable](source any, less func(T, T) bool) (T, bool)`: Returns the minimum element using the provided
comparison function. Returns `(zeroValue, false)` if the source is empty.

- `Max[T comparable](source any, less func(T, T) bool) (T, bool)`: Returns the maximum element using the provided
comparison function. Returns `(zeroValue, false)` if the source is empty.

- `Sum[T INumeric](source any) T`: Returns the sum of all numeric elements. `INumeric` includes `int`, `int8`, `int16`,
`int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`.

- `Average[T INumeric](source any) (float64, bool)`: Returns the average of all numeric elements as `float64`. Returns
`(0, false)` if the source is empty.

```go
min, _ := streams.Min[int]([]int{5, 3, 8, 1}, func(a, b int) bool { return a < b })
// min = 1

max, _ := streams.Max[int]([]int{5, 3, 8, 1}, func(a, b int) bool { return a < b })
// max = 8

sum := streams.Sum[int]([]int{1, 2, 3, 4, 5})
// sum = 15

avg, _ := streams.Average[int]([]int{10, 20, 30})
// avg = 20.0
```

#### Set operations

- `Union[T comparable](sourceA, sourceB any) IList[T]`: Returns a new `IList` containing the distinct union of elements
from both sources, preserving order of first appearance.

- `Intersect[T comparable](sourceA, sourceB any) IList[T]`: Returns a new `IList` containing elements that are present
in both sources.

```go
union := streams.Union[int]([]int{1, 2, 3}, []int{3, 4, 5})
// union contains [1, 2, 3, 4, 5]

intersect := streams.Intersect[int]([]int{1, 2, 3}, []int{2, 3, 4})
// intersect contains [2, 3]
```

#### Factory functions

- `Range(start, end int) IStream[int]`: Creates a stream of integers from `start` (inclusive) to `end` (exclusive).

- `Repeat[T comparable](value T, n int) IStream[T]`: Creates a stream containing the given value repeated `n` times.

```go
rangeStream := streams.Range(0, 5).ToArray()
// [0, 1, 2, 3, 4]

repeatStream := streams.Repeat[string]("hello", 3).ToArray()
// ["hello", "hello", "hello"]
```

#### Sorting utilities

- `Sort[T ISortable](arr []T, desc ...bool)`: Sorts a slice of comparable types in-place. `ISortable` includes all
numeric types and `string`.

- `ComparableFn[T ISortable](desc ...bool) SortFunc[T]`: Creates a `SortFunc[T]` for built-in sortable types that can
be used with `IStream.Sort()`.

```go
arr := []int{5, 3, 1, 4, 2}
streams.Sort(arr) // arr is now [1, 2, 3, 4, 5]

sorted := streams.From[int](arr).Sort(streams.ComparableFn[int]()).ToArray()
```

#### Predefined mappers

Access built-in conversion functions via `Mappers()`:

- `Mappers().IntToString()`: Returns a `ConvertFunc[int, string]`.
- `Mappers().StringToInt(errorHandler ...func(string, error))`: Returns a `ConvertFunc[string, int]`. Optionally accepts
an error handler for invalid conversions.

```go
intArr := streams.Map[string, int](
    []string{"1", "2", "3"},
    streams.Mappers().StringToInt(),
)
```

### Parallelism

Parallel operations are allowed with the streams, once enabled, when the intermediate operations are processed, they
will be done in parallel. Note that parallelism for `ForEach` operations are handled separately, since `ForEach` is a
terminal operation and executing it in parallel cannot guarantee the final order of the intermediate operations.

###### Enabling parallelism for intermediate operations.

There are multiple ways of enabling parallelism in a stream.

- As a variadic argument at the time of creation. (Supported by `From`, `FromArray` and `FromCollection`)

```go
stream := streams.
    FromArray[T](array, 8).      // Creates the stream from the given array, the second (variadic) argument indicates 
                                 // the amount of threads to be used when executing intermediate operations
    Filter(filterHandler).
    Filter(anotherFilterHandle).
    Sort(sorterHandler)
```

- Explicitly to the stream at any point before a terminal operation

```go
stream := streams.
    FromArray[T](array).
    SetThreads(8).           // Sets the amount of threads to be used when executing intermediate operations
    Filter(filterHandler).
    Except(exceptHandler).
    Sort(sorterHandler)
```

When setting the amount of threads, the amount will be capped by the maximum number of CPUs available on the host
machine. Any number equal to or lower than `0` can be provided to use the maximum available threads.

###### Executing a ForEach in parallel

The `ForEach` function by default does not run in parallel, regardless if threads were previously assigned for
intermediate operations. The reason is that intermediate operations produce a result which may have been sorted.
When running a `ForEach` in parallel, the order of execution cannot be guaranteed.

Use `ParallelForEach` only when the order of execution does not matter.

```go
stream := streams.
    FromArray[T](array, 8).
    Filter(filterHandler).
    Except(exceptHandler).
    Sort(sorterHandler).
    ParallelForEach(func(x T) {
        // foreach logic here.
    }, 0)                     // 0 = use maximum available cores
```

### Collections

This library provides several collection types:

- `ICollection[T comparable]`: Base collection interface with `Add`, `Remove`, `Contains`, `Len`, `Clear`, `ToArray`,
  `ForEach`, and iterator support.
- `IList[T comparable]`: Extends `ICollection` with index-based access (`Index`, `RemoveAt`), `Distinct`, and `Stream`.
- `ISet[T comparable]`: A collection that guarantees unique values.
- `IMap[K comparable, V any]`: A map that also implements `IList[*KeyValuePair[K, V]]` with `Get`, `Set`, `ContainsKey`,
  `Keys`, `Delete`, and `ToMap`.

#### Creating collections

```go
list := streams.NewList[int]([]int{1, 2, 3})
set := streams.NewSet[int]()
m := streams.NewMap[string, int](map[string]int{"a": 1, "b": 2})
```
