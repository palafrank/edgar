package edgar

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Timestamp time.Time

func (t Timestamp) String() string {
	return time.Time(t).Format("2006-01-02")
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*t = getDate(s)
	return nil
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

func getDate(dateStr string) Timestamp {
	year := getYear(dateStr)
	month := getMonth(dateStr)
	day := getDay(dateStr)
	ts := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return Timestamp(ts)
}

func getDateString(ts time.Time) string {
	return ts.Format("2006-01-02")
}
