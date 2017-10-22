## Stream Collections for Go, inspired in Java 8 Streams and .NET Linq

Provides structs to support functional-style operations on streams of elements, such as map-reduce transformations on collections, filtering, sorting and parallel foreach operations. For example:

```Go
// Given the following array
var fruitArray = []string{"peach", "apple", "pear", "plum", "pineapple", "banana", "kiwi", "orange"}

fruitsThatStartWithP := streams.
    FromArray(testArray).                               // Creates a stream from the given array
        Filter(func(v interface{}) bool {               // Applies the given function filter
            return strings.HasPrefix(v.(string), "p")
        }).
        ToArray().([]string)                            // Converts back to an array

```
The resulting array will be `{"peach", "pear", "plum", "pineapple"}`

Here we use an array of string as the source of the stream, perform a filter operation provided by a function that receives
a single element of the collection and determines whether the element should remain in the stream by returning a boolean.

The characteristics of a Stream are inspired in the stream features provided by Java 8. The following characteristics apply
in the go-streams.

- No storage. A stream is not a data structure that stores elements; instead, it conveys elements from a source such as a data structure, an array, a generator function, or an I/O channel, through a pipeline of computational operations.
- Functional in nature. An operation on a stream produces a result, but does not modify its source. For example, filtering a Stream obtained from a collection produces a new Stream without the filtered elements, rather than removing elements from the source collection.
- Laziness-seeking. Many stream operations, such as filtering and sorting, can be implemented lazily, exposing opportunities for optimization. For example, "find the first String with three consecutive vowels" need not examine all the input strings. Stream operations are divided into intermediate (Stream-producing) operations and terminal (value- or side-effect-producing) operations. Intermediate operations are always lazy.
- Consumable. The elements of a stream are only visited once during the life of a stream. Like an Iterator, a new stream must be generated to revisit the same elements of the source.

Currently Streams can be obtain in two ways:
- From an array via `streams.FromArray(array interface{})`, where the `array` parameter can be any array or slice implementation.
- From an implementation of `streams.ICollection` via `streams.FromCollection(collection ICollection)`

**Streams specifically for basic types, such as `StringStream` and `IntStream` will be coming soon**

## Stream operations and pipelines

Stream operations are divided into intermediate and terminal operations, and are combined to form stream pipelines. A stream pipeline consists of a source (such as a Collection, an array, a generator function, or an I/O channel); followed by zero or more intermediate operations such as `Filter()`, `Exclude()` or `Sort()` and a terminal operation such as `ForEach()` or `First()`

Intermediate operations return a new stream. They are always lazy; executing an intermediate operation such as `Filter()` or `OrderBy()` does not actually perform any action, but instead register the action into the stream that, when traversed, will execute all filtering and sorting criteria at once. Traversal of the pipeline source does not begin until the terminal operation of the pipeline is executed, such as `First()`, `Last()`, `ToArray()`

Terminal operations, such as `ForEach()`, `First()`, `Last()`, `ParallelForEach()` or `ToArray()`, may traverse the stream to produce a result or a side-effect. After the terminal operation is performed, the stream pipeline is considered consumed, and can no longer be used; if you need to traverse the same data source again, you must return to the data source to get a new stream. In almost all cases, terminal operations are eager, completing their traversal of the data source and processing of the pipeline before returning