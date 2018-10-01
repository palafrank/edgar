package edgar

import (
	"strconv"
	"strings"
)

func normalizeNumber(str string) float64 {
	negative := float64(1)
	//Remove any leading spaces or $ signs
	if strings.Contains(str, "(") && strings.Contains(str, ")") {
		negative *= -1
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
		num, err := strconv.ParseFloat(s1, 64)
		if err == nil {
			return num * negative
		}
	}
	return 0
}

func filingScale(strs []string) map[scaleEntity]scaleFactor {
	ret := make(map[scaleEntity]scaleFactor)
	ret[scaleEntityShares] = scaleNone
	ret[scaleEntityMoney] = scaleMillion
	for _, str := range strs {
		for key, val := range filingScales {
			if strings.Contains(strings.ToLower(str), strings.ToLower(key)) {
				//Some scale available in this line
				ret[val.entity] = val.scale
			}
		}
	}
	return ret
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
