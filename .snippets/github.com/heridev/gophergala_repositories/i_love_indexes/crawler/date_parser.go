package crawler

import (
	"regexp"
	"strconv"
	"time"

	"gopkg.in/errgo.v1"
)

var (
	lighttpdDateRegexp = regexp.MustCompile(`^(\d{4})-([A-Za-z]{3})-(\d{2})\s(\d{2}):(\d{2}):(\d{2})$`)
	apacheDateRegexp   = regexp.MustCompile(`^(\d{2})-([A-Za-z]{3})-(\d{4})\s(\d{2}):(\d{2})$`)
	apacheDateRegexp2  = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})\s(\d{2}):(\d{2})$`)
	monthMap           = map[string]time.Month{
		"Jan": time.January,
		"Feb": time.February,
		"Mar": time.March,
		"Apr": time.April,
		"May": time.May,
		"Jun": time.June,
		"Jul": time.July,
		"Aug": time.August,
		"Sep": time.September,
		"Oct": time.October,
		"Nov": time.November,
		"Dec": time.December,
	}
)

func LighttpdParseDate(date string) (time.Time, error) {
	values := lighttpdDateRegexp.FindAllStringSubmatch(date, -1)
	if len(values) != 1 || len(values[0]) != 7 {
		return time.Time{}, errgo.Newf("invalid format: %v", date)
	}
	return time.Date(
		mustInt(values[0][1]), monthMap[values[0][2]],
		mustInt(values[0][3]), mustInt(values[0][4]),
		mustInt(values[0][5]), mustInt(values[0][6]),
		0, time.Local), nil
}

func ApacheParseDate(date string) (time.Time, error) {
	if apacheDateRegexp.MatchString(date) {
		values := apacheDateRegexp.FindAllStringSubmatch(date, -1)
		if len(values) != 1 || len(values[0]) != 6 {
			return time.Time{}, errgo.Newf("invalid format: %v", date)
		}
		return time.Date(
			mustInt(values[0][3]), monthMap[values[0][2]],
			mustInt(values[0][1]), mustInt(values[0][4]),
			mustInt(values[0][5]), 0, 0, time.Local), nil
	} else if apacheDateRegexp2.MatchString(date) {
		values := apacheDateRegexp2.FindAllStringSubmatch(date, -1)
		if len(values) != 1 || len(values[0]) != 6 {
			return time.Time{}, errgo.Newf("invalid format: %v", date)
		}
		return time.Date(
			mustInt(values[0][1]), time.Month(mustInt(values[0][2])),
			mustInt(values[0][3]), mustInt(values[0][4]),
			mustInt(values[0][5]), 0, 0, time.Local), nil
	} else {
		return time.Time{}, errgo.Newf("unkown format: %v", date)
	}
}

func mustInt(s string) int {
	return int(mustInt64(s))
}

func mustInt64(s string) int64 {
	if s == "-" {
		return -1
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
