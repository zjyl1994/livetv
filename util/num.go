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

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func String2Int64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i64
	}
}
