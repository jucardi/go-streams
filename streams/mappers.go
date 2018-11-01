package streams

import "strconv"

/** This file provides a few predefined mappers that can be used with the steams.Map **/

// MapIntToString returns a ConvertFunc which maps an int to a string.
//
//  Eg:  strArray := streams.From(intArray).Map(MapIntToString()).ToArray().([]string)
//
func MapIntToString() ConvertFunc {
	return func(i interface{}) interface{} {
		return strconv.Itoa(i.(int))
	}
}

// MapStringToInt returns a ConvertFunc which maps a string to an int.
//
// - errorHandler: Optional variadic arg, if provided, it will be invoked if the string to int
//                 conversion fails.
//
//  Eg:  errHandler := func(str string, err error) {
//           log.Errorf("unable to convert %s to int, %s", str, err.Error())
//       }
//       intArray := streams.From(strArray).Map(MapStringToInt(errHandler)).ToArray().([]int)
//
func MapStringToInt(errorHandler ...func(string, error)) ConvertFunc {
	return func(x interface{}) interface{} {
		str := x.(string)
		i, err := strconv.Atoi(str)
		if err != nil && len(errorHandler) > 0 && errorHandler[0] != nil {
			errorHandler[0](str, err)
		}
		return i
	}
}
