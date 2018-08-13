package main

import (
	"strconv"
	"strings"
)

func normalizeNumber(str string) uint64 {
	s := strings.Split(str, ",")
	if len(s) > 0 {
		var s1 string
		for _, data := range s {
			s1 += data
		}
		num, err := strconv.Atoi(s1)
		if err == nil {
			return uint64(num)
		}
	}
	return 0
}
