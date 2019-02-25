package edgar

import (
	"errors"
	"math"
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

func filingScale(strs []string, t filingDocType) map[scaleEntity]scaleFactor {
	ret := make(map[scaleEntity]scaleFactor)
	if t == filingDocEN {
		ret[scaleEntityShares] = scaleNone
	} else {
		ret[scaleEntityShares] = scaleMillion
	}
	ret[scaleEntityMoney] = scaleMillion
	ret[scaleEntityPerShare] = scaleNone
	for _, str := range strs {
		s := strings.ToLower(str)
		parts := strings.Split(s, ",")
		for _, part := range parts {
			if strings.Contains(part, "share") {
				// Share scale
				if strings.Contains(part, "thousand") {
					ret[scaleEntityShares] = scaleThousand
				} else if strings.Contains(part, "million") {
					ret[scaleEntityShares] = scaleMillion
				}
			} else if strings.Contains(part, "$") || strings.Contains(part, "usd") {
				//Money scale
				if strings.Contains(part, "thousand") {
					ret[scaleEntityMoney] = scaleThousand
				} else if strings.Contains(part, "billion") {
					ret[scaleEntityMoney] = scaleBillion
				}
			}
		}
	}
	return ret
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

func round(val float64) float64 {
	return math.Floor(val*100) / 100
}

// For now we will say that two numbers are in the same scale
// if it is within 50% of each other
// This is a fallback cross check if the scale used for metrics is accurate
func isSameScale(one float64, two float64) bool {
	val := (one - two) / two
	if one < two {
		val = (two - one) / one
	}
	if val <= 1 {
		return true
	}
	return false
}
