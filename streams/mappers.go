package streams

import "strconv"

/** This file provides a few predefined mappers that can be used with the steams.Map **/

// IMappers defines a list of built-in ConvertFunc mappers for different values types to be using in
// mapping functions such as `Map[F, T]` and `MapNonComparable[F, T]`
type IMappers interface {
	// IntToString returns a ConvertFunc which maps an int to a string.
	IntToString() ConvertFunc[int, string]
	// StringToInt returns a ConvertFunc which maps a string to an int.
	//
	//   - errorHandler: Optional variadic arg, if provided, it will be invoked if the string to int
	//     conversion fails.
	//
	//     Eg:  errHandler := func(str string, err error) {
	//     log.Errorf("unable to convert %s to int, %s", str, err.Error())
	//     }
	//     intArray := streams.From(strArray).Map(MapStringToInt(errHandler)).ToArray().([]int)
	StringToInt(errorHandler ...func(string, error)) ConvertFunc[string, int]
}

var defaultMappers mappers

type mappers struct{}

func (mappers) IntToString() ConvertFunc[int, string] {
	return strconv.Itoa
}

func (mappers) StringToInt(errorHandler ...func(string, error)) ConvertFunc[string, int] {
	return func(x string) int {
		i, err := strconv.Atoi(x)
		if err != nil && len(errorHandler) > 0 && errorHandler[0] != nil {
			errorHandler[0](x, err)
		}
		return i
	}
}

func mapStream[From, To comparable](from IStream[From], f ConvertFunc[From, To]) IList[To] {
	return NewList[To](mapIterable[From, To](from.ToCollection(), f))
}

func mapIterable[From, To any](from IIterable[From], f ConvertFunc[From, To]) []To {
	return mapIterator[From, To](from.Iterator(), f)
}

func mapIterator[From, To any](from IIterator[From], f ConvertFunc[From, To]) (ret []To) {
	for old := from.Current(); from.HasNext(); old = from.Next() {
		n := f(old)
		ret = append(ret, n)
	}

	return
}
