package numconvert

import "strconv"

func BytesToInt64(in string) int64 {
	valuesInt64, _ := strconv.ParseInt(in, 10, 64)
	return valuesInt64
}

//string to Float64
func StringToFloat64(values string) float64 {
	valuesFloat, _ := strconv.ParseFloat(values, 64)
	return valuesFloat
}
