package main

import (
	"strconv"
	"strings"
)

func normalizeNumber(str string) int64 {
	//Remove any leading spaces or $ signs
	str = strings.TrimLeft(str, " ")
	str = strings.TrimLeft(str, "$")
	str = strings.TrimLeft(str, " ")
	str = strings.TrimRight(str, " ")
	s := strings.Split(str, ",")
	if len(s) > 0 {
		var s1 string
		for _, data := range s {
			s1 += data
		}
		num, err := strconv.Atoi(s1)
		if err == nil {
			return int64(num)
		}
	}
	return 0
}
