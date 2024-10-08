package utils

import "time"

func PointerFromFloat64(a float64) *float64 {
	return &a
}

func TimeMustParse(layout, timestamp string) time.Time {
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		panic(err)
	}
	return t
}
