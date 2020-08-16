package util

import "strconv"

func String2Uint(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return uint(i)
	}
}
