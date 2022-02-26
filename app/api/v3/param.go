package v3

import "strconv"

func mustStringToUint(v string) uint {
	num, _ := strconv.Atoi(v)
	return uint(num)
}
