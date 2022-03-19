package v3

import "strconv"

func mustStringToInt(v string) int {
	num, _ := strconv.Atoi(v)
	return num
}

func mustStringToUint(v string) uint {
	num, _ := strconv.Atoi(v)
	return uint(num)
}
