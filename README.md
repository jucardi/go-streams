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
go get github.com/jucardi/go-streams@latest
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

Now let's do a simple forach operation

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
pineaple
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

Currently Streams can be obtained in two ways:
- From an array via `streams.FromArray[T comparable](array []T)`, where the `array` parameter can be any array or slice
implementation.
- From an implementation of `streams.ICollection[T comparable]` via `streams.FromCollection[T comparable](col ICollection[T])`

**NOTE:** There is a global function `From[T comparable]` that accepts any value and properly creates the stream depending
on whether it is an array or a collection. However it panics if the input is neither.

## Stream operations and pipelines

Stream operations are divided into intermediate and terminal operations, and are combined to form stream pipelines. A 
stream pipeline consists of a source (such as an iterable, a collection, an array, a generator function, or an I/O 
channel); followed by zero or more intermediate operations such as `Filter()`, `Exclude()` or `Sort()` and a terminal
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

- `Filter(function (x T) bool)`: Filters any element that does not meet the condition provided by the function.
- `Except(function (x T) bool)`: Filters all elements that meet the condition provided by the function.
- `Sort(function (x, y T) int, desc ...bool)`: Sorts the elements ascending in the stream using the provided comparable
  function.
- `SetThreads(threads int)`: Sets the number of threads to be used when processing the stream
- `Distinct()`: Ensures that the resulting stream operation only includes unique values


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
- `ToArray()`: Returns an array of elements from the resulting stream
- `ToCollection()`: Returns a `ICollection[T]` of elements from the resulting stream
- `ToIterable()`: Returns a `IIterable[T]` of elements from the resulting stream
- `ToList()`: Returns a `IList[T]` of elements from the resulting stream
- `ToDistinct()`: Returns a `ISet[T]` of elements from the resulting stream with only unique values

###### Boolean returns
- `IsEmpty()`: Indicates whether the resulting stream contains no elements
- `Contains(value T)`: Indicates whether the provided value matches any of the values in the stream.
- `AnyMatch(f func(x T) bool)`: Indicates whether ANY elements of the stream match the given condition function.
- `AllMatch(f func(x T) bool)`: Indicates whether ALL elements of the stream match the given condition function
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
    Sort(sorterHandler)
    IfEmpty().
    Then(func (resultStream IStream[T] { 
        // Do something if empty
    }).
    Else(func (resultStream IStream[T] { 
        // Do something if not empty
    })
```

- `IfEmpty()`: Returns a `IThen[T]` handler if empty
- `IfAnyMatch(f ConditionalFunc[T])`: Returns a `IThen[T]` handler if any elements match the provided condition
- `IfAllMatch(f ConditionalFunc[T])`: Returns a `IThen[T]` handler if all elements match the provided condition
- `IfNoneMatch(f ConditionalFunc[T])`: Returns a `IThen[T]` handler if none elements match the provided condition

###### Void returns
- `ForEach(f func(x T))`: Iterates over all elements in the stream calling the provided function.
- `ParallelForEach(c T), threads int, skipWait ...bool)`: Iterates over all elements in the stream calling the provided
function. Creates multiple go channels to parallelize the operation. ParallelForeach does not use any thread values 
previously provided in any filtering method nor enables parallel filtering if any filtering is done prior to the 
`ParallelForEach` phase. Only use `ParallelForEach` if the order in which the elements are processed does not matter, 
otherwise see `ForEach`.

###### Int returns
- `Count()`: Counts the elements contained by the resulting stream.

### Mapping functions

In previous versions of this library, there was a `Map()` function that would help to convert the elements of the initial
type of the stream into something else (E.g: array of strings to integers)

This library was updated in favor of Golang Generics to avoid having to work with `interface{}` and `reflect`. However,
Golang Generics do not support adding additional generic types to functions that are defined in interfaces or structs,
so the `Map` function cannot be invoked from the stream since it cannot accept the target type of the conversion.

For this reason, the `Map` function was moved as a static function in this package, and new mapping functions were added
to simplify some mapping invocations

The current mapping functions built-in this package are:

- `Map[From, To comparable](source any, f ConvertFunc[From, To]) IList[To]`: Map maps the elements of the source to a new
element, using the mapping function provided. Outputs a new IList with the converted elements. This function accepts the
following sources where `From` and `To` are `comparable` types:
  - `[]From`
  - `IIterable[From]`
  - `ICollection[From]`
  - `IList[From]`
  - `IIterator[From]`
  - `IStream[From]`



- `MapNonComparable[From, To any](source any, f ConvertFunc[From, To]) []To`: This function is similar to `Map`, maps
the elements of the source to a new element, using the mapping function provided. Outputs an array containing the mapped
elements instead of a collection, and unlike `Map` the source accepts non-comparable types. This function accepts the
following sources where `From` and `To` accept any type including non-comparable:
  - `[]From`
  - `IIterable[From]`
  - `ICollection[From]`
  - `IList[From]`
  - `IIterator[From]`


- `MapToPtr[T any](source any) []*T`: Converts a collection of `T` to a collection of `*T`. Many structs are not `comparable`
which makes `T` unsupported by `IStream`, `ICollection`, `IList` and `ISet`. Pointers however are considered comparable,
so this function outputs an array of `*T` which can be used in the types mentioned.

###### Examples:
*1. Mapping from a stream*
```go
// Given the following source array of number strings
sourceArray := []string{"1", "5", "8", "100", "23", "6", "abc"}

mappedList := streams.Map[string, int](   
    streams.                            
        From[string](sourceArray).      // The source stream for the mapping with any stream functions needed prior to
        Filter(func(x string) bool {    // be processed for mapping. In this example, all items that cannot be parsed to
            _, err := strconv.Atoi(x)   // an int will be filtered
            return err == nil
        }),
    func(item string) int {             // The mapping function
        ret, _ := strconv.Atoi(x)
        return ret
    }
})
```
*2. Using the mapping result as a stream to continue processing*
```go
// Given the following source array of number strings
sourceArray := []string{"1", "5", "8", "100", "23", "6", "abc"}

mappedList := streams.Map[string, int](   
    sourceArray,                       // Using the array as the source instead of a stream
    func(item string) int {            // The mapping function
        ret, _ := strconv.Atoi(x)
        return ret                     // for any elements that cannot be converted to int, returns default int (0)
    }
}).
    ToStream().                        // Obtain a stream from the resulting list
    Filter(func (x int) bool {         // Apply any stream functions to the new stream
        return x > 0
    }).
    ToList()
```
*3. Mapping using the `MapNoComparable` function*
```go
// Given the following struct
type SomeStruct struct {
    Name       string
    Score      int
    StringFn   func() string
}

// And given the following array
arr := []SomeStruct{
    {
        Name: "abcd",
        StringFn: func() string {
            return "something"
        },
    },
    {
        Score: 123,
        StringFn: func() string {
            return "something else"
        },
    },
}

// Because the structure has a `func` field, the structure `SomeStruct` is not a `comparable` type,
// so that's where the func `MapNoComparable` becomes handy. This example transforms the source array
// into an array functions []func() string using the function in the `StringFn` field

funcs := streams.
    MapNonComparable[SomeStruct, func() string](
        arr,
        func(x SomeStruct) func() string {
            return x.StringFn
        })
```
*4. Using the `MapToPtr` function*
```go
// Using the same struct and array as the example above

// Given the following struct
type SomeStruct struct {
    Name       string
    Score      int
    StringFn   func() string
}

// And given the following array
arr := []SomeStruct{
    {
        Name: "abcd",
        StringFn: func() string {
            return "something"
        },
    },
    {
        Score: 123,
        StringFn: func() string {
            return "something else"
        },
    },
}


// Say in this case, we would like to obtain an array of functions []func() string contained in the
// source structure, but only for those structures that have a `Score` of 50 or higher.
//
// streams.From[SomeStruct] will not work because `SomeStruct` is not `comparable`, so we can do the
// following:

// Convert `arr` from []SomeStruct to []*SomeStruct
newArr := MapToPtr[SomeStruct](arr)

// Now, `newArr` can be used in streams because *SomeStruct is `comparable`
funcs := streams.
    Map[*SomeStruct, func() string](
        streams.                               // The source stream for mapping
            From[*SomeStruct](newArr).
            Filter(func (x *SomeStruct){       // The filtering in the source stream
                return x.Score >= 50
            }),
        func(x *SomeStruct) func() string {    // The mapping function
            return x.StringFn
        },
    )
 
```

### Parallelism

Parallel operations are allowed with the streams, once enabled, when the intermediate operations are processed, they
will be done in parallel. Note that parallelism for `ForEach` operations are handled separately, since `ForEach` is a
terminal operation and executing it in parallel cannot guarantee the final order of the intermediate operations.

###### Enabling parallelism for intermediate operations.

There are multiple ways of enabling parallelism in a stream.

- As a variadic argument at the time of creation. (Supported by `From`, `FromArray` and `FromIterator`)

```go
stream := streams.
    FromArray[T](array, 8).      // Creates the stream from the given array, the second (variadic) argument indicates 
                                 // the amount of threads to be used when executing intermediate operations
    Filter(filterHandler).
    Filter(anotherFilterHandle).
    Sort(sorterHandler)
```

- As a variadic argument when adding a filter or except

```go
stream := streams.
    FromArray[T](array).
    Filter(filterHandler, 8).    // Adds a filter to the stream, the second (variadic) argument indicates the amount
                                 // of threads to be used when executing intermediate operations
    Except(exceptHandler).
    Sort(sorterHandler)

stream := streams.
    FromArray[T](array).
    Filter(filterHandler).
    Except(exceptHandler, 8).    // Adds an except to the stream, the second (variadic) argument indicates the amount of threads to be used when executing intermediate operations
    Sort(sorterHandler)
```

- Explicitly to the stream at any point before a terminal operation

```go
stream := streams.
    FromArray(array).
    SetThreads(8).           // Sets the amount of threads to be used when executing intermediate operations
    Filter(filterHandler).
    Except(exceptHandler).
    Sort(sorterHandler)

stream := streams.
    FromArray(array).
    Filter(filterHandler).
    Except(exceptHandler).
    Sort(sorterHandler).
    SetThreads(8)            // Can be added at any point before invoking a terminal operation
```

When setting the amount of threads at any of the options above, keep in mind, the amount of threads will be capped by the maximum amount of CPUs available in the host machine.

Any number equal or lower than `0` can be provided if the maximum amount of threads is desired based on the CPUs cap.

###### Executing a For Each in parallel

The `ForEach` function by default does not run in parallel, regardless if threads were previously assigned for intermediate operations.
The reason for this design is, intermediate operations produce a result which may have been sorted by a provided sorting algorithm.
When running a `ForEach` in parallel, the order in which the foreach handlers will be ran for each element in the resulting stream cannot be guaranteed.

Use `ParallelForEach` only when the order of execution does not matter when processing the elements of the stream

The function `ParallelForeach(f function(x T), threads int, skipWait ...bool)` receives 3 args.

- `f`: The handler which will process each single element of the stream.
- `threads`: Indicates the amount of parallel go channels `ParallelForEach` will use. This will be capped by the amount of available CPUs in the host machine. Any number equal or lower than `0` can be provided if the maximum amount of threads is desired based on the CPUs cap.
- `skipWait`: This is an optional variadic argument that indicates if `ParallelForEach` should wait until all of the go channels are done executing. If set to `true`, the function `ParallelForEach` will return immediately after the intermediate operations are done and the `ParallelForEach` process will run in the background.

Example:

```go
stream := streams.
    FromArray[T](array, 8).   // Setting 8 cores for intermediate operations, such as filtering.
    Filter(filterHandler).
    Except(exceptHandler).
    Sort(sorterHandler).
    ParallelForEach(function (x T) {

        // foreach logic here.

    }, 0)                     // Setting the maximum amount of cores available for the ParallelForEach process.
```