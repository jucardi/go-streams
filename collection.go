package main

import (
	"fmt"
	"github.com/jucardi/collections/streams"
	"reflect"
)

// Returns the first index of the target string `t`, or
// -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Returns `true` if the target string t is in the
// slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

// Returns `true` if one of the strings in the slice
// satisfies the predicate `f`.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// Returns `true` if all of the strings in the slice
// satisfy the predicate `f`.
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Returns a new slice containing all strings in the
// slice that satisfy the predicate `f`.
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Returns a new slice containing the results of applying
// the function `f` to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func main() {
	var strs = []string{"peach", "apple", "pear", "plum"}

	var appleFunc = func(x reflect.Value) bool {
		return "apple" == x.String()
	}

	var trueFunc = func(x reflect.Value) bool {
		return true
	}

	result := streams.From(strs).
		Filter(appleFunc).
		//AllMatch(trueFunc)
		First().String()

	fmt.Println(result)

	match := streams.From(strs).AllMatch(trueFunc)

	fmt.Println(match)
	fmt.Println(streams.From(strs).Contains("apple"))
	//f := func(x reflect.Value) bool {
	//	return "apple" == x.String()
	//}
	//
	//fmt.Print(f)
	//Filter(func(x interface{}) bool {
	//	return "peach" == string(x)
	//}).
	//First())
	//	// Here we try out our various collection functions.
	//	var strs = []string{"peach", "apple", "pear", "plum"}
	//
	//	fmt.Println(Index(strs, "pear"))
	//
	//	fmt.Println(Include(strs, "grape"))
	//
	//	fmt.Println(Any(strs, func(v string) bool {
	//		return strings.HasPrefix(v, "p")
	//	}))
	//
	//	fmt.Println(All(strs, func(v string) bool {
	//		return strings.HasPrefix(v, "p")
	//	}))
	//
	//	fmt.Println(Filter(strs, func(v string) bool {
	//		return strings.Contains(v, "e")
	//	}))
	//
	//	// The above examples all used anonymous functions,
	//	// but you can also use named functions of the correct
	//	// type.
	//	fmt.Println(Map(strs, strings.ToUpper))
	//
}
