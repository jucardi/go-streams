package streams

import "sort"

// ISortable comprises the comparable types that also support   <   >   <=   >=   comparison instead of just  ==
type ISortable interface {
	string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// ComparableFn creates a new SortFunc[T, T] that knows how to compare default comparable values
func ComparableFn[T ISortable](desc ...bool) SortFunc[T] {
	if len(desc) > 0 && desc[0] {
		return func(x, y T) int {
			return defaultComparableFunc[T](y, x)
		}
	}
	return defaultComparableFunc[T]
}

func Sort[T ISortable](arr []T, desc ...bool) {
	d := len(desc) > 0 && desc[0]

	sort.SliceIsSorted(arr, func(i, j int) bool {
		if d {
			return arr[i] < arr[j]
		}
		return arr[i] > arr[j]
	})
}

func defaultComparableFunc[T ISortable](a, b T) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}

	return +1
}
