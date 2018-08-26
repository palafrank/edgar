package edgar

import (
	"strconv"
	"strings"
)

func normalizeNumber(str string) int64 {
	negative := int64(1)
	//Remove any leading spaces or $ signs
	if strings.Contains(str, "(") && strings.Contains(str, ")") {
		negative = -1
	}
	str = strings.TrimLeft(str, " ")
	str = strings.TrimLeft(str, "$")
	str = strings.TrimLeft(str, " ")
	str = strings.TrimRight(str, " ")
	str = strings.TrimLeft(str, "(")
	str = strings.TrimRight(str, ")")

	//TODO: Ignoring decimals for now
	s := strings.Split(str, ".")
	s = strings.Split(s[0], ",")

	if len(s) > 0 {
		var s1 string
		for _, data := range s {
			s1 += data
		}
		num, err := strconv.Atoi(s1)
		if err == nil {
			return int64(num) * negative
		}
	}
	return 0
}

func getYear(date string) int {
	strs := strings.Split(date, "-")
	if len(strs) != 3 {
		return 0
	}
	year, _ := strconv.Atoi(strs[0])
	return year
}

func getMonth(date string) int {
	strs := strings.Split(date, "-")
	if len(strs) != 3 {
		return 0
	}
	year, _ := strconv.Atoi(strs[1])
	return year
}

func getDay(date string) int {
	strs := strings.Split(date, "-")
	if len(strs) != 3 {
		return 0
	}
	year, _ := strconv.Atoi(strs[2])
	return year
}

func getDate(dateStr string) Date {
	var d date
	d.year = getYear(dateStr)
	d.month = getMonth(dateStr)
	d.day = getDay(dateStr)
	return d
}
