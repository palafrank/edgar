package edgar

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func normalizeNumber(str string) (float64, error) {
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

	dec := ""
	s := strings.Split(str, ".")
	if len(s) > 1 {
		dec = s[1]
	}
	s = strings.Split(s[0], ",")

	if len(s) > 0 {
		var s1 string
		for _, data := range s {
			s1 += data
		}
		if dec != "" {
			s1 += "."
			s1 += dec
		}
		num, err := strconv.ParseFloat(s1, 64)
		if err == nil {
			return num * negative, nil
		}
	}
	return 0, errors.New("Error normalizing number")
}

func filingScale(strs []string) map[scaleEntity]scaleFactor {
	ret := make(map[scaleEntity]scaleFactor)
	ret[scaleEntityShares] = scaleNone
	ret[scaleEntityMoney] = scaleMillion
	ret[scaleEntityPerShare] = scaleNone
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

func getFinDataXBRLTag(onclick string) (string, error) {
	if strings.Contains(onclick, "showAR") {
		d := strings.Split(onclick, `'`)
		if len(d) == 3 {
			if strings.Contains(d[1], "defref") {
				return d[1], nil
			}
		}
	}
	return "", errors.New("Not a financial tag")
}

func setCollectedData(data interface{}, fieldNum int) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	bit, bitOk := t.Field(fieldNum).Tag.Lookup("bit")
	if bitOk {
		bitLoc, err := strconv.Atoi(bit)
		if err == nil {
			field := v.FieldByName("CollectedData")
			if field.CanSet() {
				var mask uint64 = 0x01
				obj := field.Uint()
				obj |= mask << uint8(bitLoc)
				field.SetUint(obj)
			}
		}
	}
}

func clearCollectedData(data interface{}, fieldNum int) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	bit, bitOk := t.Field(fieldNum).Tag.Lookup("bit")
	if bitOk {
		bitLoc, err := strconv.Atoi(bit)
		if err == nil {
			field := v.FieldByName("CollectedData")
			if field.CanSet() {
				var mask uint64 = 0x01
				obj := field.Uint()
				obj &= ^(mask << uint8(bitLoc))
				field.SetUint(obj)
			}
		}
	}
}

func isCollectedDataSet(data interface{}, fieldName string) bool {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	f, ok := t.FieldByName(fieldName)
	if !ok {
		return false
	}
	bit, bitOk := f.Tag.Lookup("bit")
	if bitOk {
		bitLoc, err := strconv.Atoi(bit)
		if err == nil {
			field := v.FieldByName("CollectedData")
			if field.CanSet() {
				var mask uint64 = 0x01
				obj := field.Uint()
				if obj&(mask<<uint8(bitLoc)) != 0 {
					return true
				}
			}
		}
	}
	return false
}
