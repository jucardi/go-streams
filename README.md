## Stream Collections for Go, inspired in Java 8 Streams and .NET Linq

This library provides
- structs to easily represent and manage collections and iterable of elements which size are not necessarily defined.
- Stream structs to support functional-style operations on streams of elements, such as map-reduce transformations on collections, filtering, sorting, mapping, foreach parallel operations.

##### Quick Start

To keep up to date with the most recent version:

```bash
go get github.com/jucardi/go-streams
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
	FromArray(fruitArray).

	// Adds a filter for strings that start with 'p'
	Filter(func(v interface{}) bool {
		return strings.HasPrefix(v.(string), "p")
	}).

	// Orders alphabetically
	OrderBy(func(a interface{}, b interface{}) int {
		return strings.Compare(a.(string), b.(string))
	}).

	// Converts back to an array
	ToArray().([]string)
```
The resulting array will be `{"peach", "pear", "pineapple", "plum"}`

Here we use an array of string as the source of the stream, perform a filter operation provided by a function that receives
a single element of the collection and determines whether the element should remain in the stream by returning a boolean.

Now let's do a simple forach operation

```go
streams.
	FromArray(fruitArray).
	Filter(func(v interface{}) bool {
		return strings.HasPrefix(v.(string), "p")
	}).
	ForEach(func(v interface) {
		println(v)
	})
```

In this example, once the stream processes the filter, performs a foreach operation with the result. With this operation we'll obtain the following output in the console

```
peach
pear
plum
pineaple
```

## About the go-streams

The characteristics of a Stream are inspired in the stream features provided by Java 8. The following characteristics apply
in the go-streams.

- No storage. A stream is not a data structure that stores elements; instead, it conveys elements from a source such as a data structure, an array, a generator function, or an I/O channel, through a pipeline of computational operations.
- Functional in nature. An operation on a stream produces a result, but does not modify its source. For example, filtering a Stream obtained from a collection produces a new Stream without the filtered elements, rather than removing elements from the source collection.
- Laziness-seeking. Many stream operations, such as filtering and sorting, can be implemented lazily, exposing opportunities for optimization. For example, "find the first String with three consecutive vowels" need not examine all the input strings. Stream operations are divided into intermediate (Stream-producing) operations and terminal (value- or side-effect-producing) operations. Intermediate operations are always lazy.
- Consumable. The elements of a stream are only visited once during the life of a stream. Like an Iterator, a new stream must be generated to revisit the same elements of the source.

Currently Streams can be obtain in two ways:
- From an array via `streams.FromArray(array interface{})`, where the `array` parameter can be any array or slice implementation.
- From an implementation of `streams.IIterable` via `streams.FromIterable(iterable IIterable)`

**Streams specifically for basic types, such as `StringStream` and `IntStream` will be coming soon**

## Stream operations and pipelines

Stream operations are divided into intermediate and terminal operations, and are combined to form stream pipelines. A stream pipeline consists of a source (such as an iterable, a collection, an array, a generator function, or an I/O channel); followed by zero or more intermediate operations such as `Filter()`, `Exclude()` or `Sort()` and a terminal operation such as `ForEach()` or `First()`

Intermediate operations return a stream. They are always lazy; executing an intermediate operation such as `Filter()` or `OrderBy()` does not actually perform any action, but instead register the action into the stream that, when traversed, will execute all filtering and sorting criteria at once. Traversal of the pipeline source does not begin until the terminal operation of the pipeline is executed, such as `First()`, `Last()`, `ToArray()`, `ForEach()`

Terminal operations, such as `ForEach()`, `First()`, `Last()`, `ParallelForEach()` or `ToArray()`, may traverse the stream to produce a result or a side-effect. After the terminal operation is performed, the stream pipeline is considered consumed, and can no longer be used; if you need to traverse the same data source again, you must return to the data source to get a new stream. In almost all cases, terminal operations are eager, completing their traversal of the data source and processing of the pipeline before returning

### Intermediate operations

- `Filter(function (x) bool)`: Filters any element that does not meet the condition provided by the function.
- `Except(function (x) bool)`: Filters all elements that meet the condition provided by the function.
- `OrderBy(function (x, y) int, desc ...bool)`: Sorts the elements ascending in the stream using the provided comparable function.
- `ThenBy(function (x, y) int, desc ...bool)`: If two elements are considered equal after previously applying a comparable function, attempts to sort ascending the 2 elements with an additional comparable function.

### Terminal operations

###### Single item returns
- `First(defaultValue ...interface{})`: Returns the first element of the resulting stream. Returns nil (or default value if provided) if the resulting stream is empty.
- `Last(defaultValue ...interface{})`: Returns the last element of the resulting stream. Returns nil (or default value if provided) if the resulting stream is empty.
- `At(i int, defaultValue ...interface{})`: Returns the element at the given index in the resulting stream. Returns nil (or default value if provided) if out of bounds.
- `AtReverse(i int, defaultValue ...interface{})`: Returns the element at the given position, starting from the last element to the first in the resulting stream. Returns nil (or default value if provided) if out of bounds.

###### Collection returns
- `ToArray()`: Returns an array of elements from the resulting stream
- `ToCollection()`: Returns a `ICollection` of elements from the resulting stream
- `ToIterable()`: Returns a `IIterable` of elements from the resulting stream
- `Map(function (x) interface{})`: Maps the elements of the resulting stream using the provided mapping function. Returns a new stream so new intermediate and terminal operations can be ran.

**Note:** *For now, `Map` works as a terminal operation, which will trigger filtering and sorting so each one of the elements in the resulting array can be mapped and a new Stream for the new value type can be returned. Post MVP, Map will be refactored so it works as an intermediate operation and gets triggered along filtering and sorting when another terminal operation is invoked.*

###### Boolean returns
- `AnyMatch(f func(interface{}) bool)`: Indicates whether ANY elements of the stream match the given condition function.
- `AllMatch(f func(interface{}) bool)`: Indicates whether ALL elements of the stream match the given condition function
- `NoneMatch(f func(interface{}) bool)`: Indicates whether NONE of elements of the stream match the given condition function.
- `Contains(value interface{})`: Indicates whether the provided value matches any of the values in the stream.

###### Void returns
- `ForEach(f func(interface{}))`: Iterates over all elements in the stream calling the provided function.
- `ParallelForEach(f func(interface{}), threads int, skipWait ...bool)`: Iterates over all elements in the stream calling the provided function. Creates multiple go channels to parallelize the operation. ParallelForeach does not use any thread values previously provided in any filtering method nor enables parallel filtering if any filtering is done prior to the `ParallelForEach` phase. Only use `ParallelForEach` if the order in which the elements are processed does not matter, otherwise see `ForEach`.

###### Int returns
- `Count()`: Counts the elements contained by the resulting stream.

## Parallelism

Parallel operations are allowed with the streams, once enabled, the when the intermediate operations are processed, they will be done in parallel. Note that parallelism for `ForEach` operations are handled separately, since `ForEach` is a terminal operation and executing it in parallel cannot guarantee the final order of the intermediate operations.

###### Enabling parallelism for intermediate operations.

There are multiple ways of enabling parallelism in a stream.

- As a variadic argument at the time of creation. (Supported by `From`, `FromArray` and `FromIterator`)

```go
stream := streams.
    FromArray(array, 8).          // Creates the stream from the given array, the second (variadic) argument indicates the amount of threads to be used when executing intermediate operations
    Filter(filterHandler).
    Filter(anotherFilterHandle).
    OrderBy(sorterHandler)
```

- As a variadic argument when adding a filter or except

```go
stream := streams.
    FromArray(array).
    Filter(filterHandler, 8).       // Adds a filter to the stream, the second (variadic) argument indicates the amount of threads to be used when executing intermediate operations
    Except(exceptHandler).
    OrderBy(sorterHandler)

stream := streams.
    FromArray(array).
    Filter(filterHandler).
    Except(exceptHandler, 8).       // Adds an except to the stream, the second (variadic) argument indicates the amount of threads to be used when executing intermediate operations
    OrderBy(sorterHandler)
```

- Explicitly to the stream at any point before a terminal operation

```go
stream := streams.
    FromArray(array).
    SetThreads(8).           // Sets the amount of threads to be used when executing intermediate operations
    Filter(filterHandler).
    Except(exceptHandler).
    OrderBy(sorterHandler)

stream := streams.
    FromArray(array).
    Filter(filterHandler).
    Except(exceptHandler).
    OrderBy(sorterHandler).
    SetThreads(8)            // Can be added at any point before invoking a terminal operation
```

When setting the amount of threads at any of the options above, keep in mind, the amount of threads will be capped by the maximum amount of CPUs available in the host machine.

Any number equal or lower than `0` can be provided if the maximum amount of threads is desired based on the CPUs cap.

###### Executing a For Each in parallel

The `ForEach` function by default does not run in parallel, regardless if threads were previously assigned for intermediate operations.
The reason for this design is, intermediate operations produce a result which may have been sorted by a provided sorting algorithm.
When running a `ForEach` in parallel, the order in which the foreach handlers will be ran for each element in the resulting stream cannot be guaranteed.

Use `ParallelForEach` only when the order of execution does not matter when processing the elements of the stream

The function `ParallelForeach(f function(interface{}), threads int, skipWait ...bool)` receives 3 args.

- `f`: The handler which will process each single element of the stream.
- `threads`: Indicates the amount of parallel go channels `ParallelForEach` will use. This will be capped by the amount of available CPUs in the host machine. Any number equal or lower than `0` can be provided if the maximum amount of threads is desired based on the CPUs cap.
- `skipWait`: This is an optional variadic argument that indicates if `ParallelForEach` should wait until all of the go channels are done executing. If set to `true`, the function `ParallelForEach` will return immediately after the intermediate operations are done and the `ParallelForEach` process will run in the background.

Example:

```go
stream := streams.
    FromArray(array, 8).       // Setting 8 cores for intermediate operations, such as filtering.
    Filter(filterHandler).
    Except(exceptHandler).
    OrderBy(sorterHandler).
    ParallelForEach(function (x interface{}) {

        // foreach logic here.

    }, 0)                     // Setting the maximum amount of cores available for the ParallelForEach process.
```