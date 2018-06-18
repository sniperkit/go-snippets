// Package dtime 工具包 时间辅助操作
package dtime

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ShortTimeLayout  = "2006-01-02"
	MiddleTimeLayout = "2006-01-02 15:04"
	LongTimeLayout   = "2006-01-02 15:04:05"
)

// JSONShortTime 只展示年月日 e.g. 2017--11-01
type JSONShortTime struct{ JSONTime }

func (t JSONShortTime) String() string {
	return t.Format(ShortTimeLayout)
}

// MarshalJSON outputs JSON presentation
func (t JSONShortTime) MarshalJSON() ([]byte, error) {
	return t.doMarshalJSON(ShortTimeLayout)
}

func (t JSONShortTime) ToDB() interface{} {
	return int64(t.JSONTime)
}

func (t *JSONShortTime) FromDB(b []byte) error {
	st, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = JSONTime(st).AsShortTime()
	return nil
}

// UnmarshalJSON unmarshal string to JSONTime
func (t *JSONShortTime) UnmarshalJSON(b []byte) error {
	jt, err := t.doUnmarshalJSON(ShortTimeLayout, b)
	if err != nil {
		return err
	}
	t.JSONTime = jt
	return nil
}

// JSONMiddleTime  展示年月日时分,e.g. 2017-11-01 14:23
type JSONMiddleTime struct{ JSONTime }

func (t JSONMiddleTime) String() string {
	return t.Format(MiddleTimeLayout)
}

func (t JSONMiddleTime) ToDB() interface{} {
	return int64(t.JSONTime)
}

func (t *JSONMiddleTime) FromDB(b []byte) error {
	st, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*t = JSONTime(st).AsMiddleTime()
	return nil
}

// MarshalJSON outputs JSON presentation
func (t JSONMiddleTime) MarshalJSON() ([]byte, error) {
	return t.doMarshalJSON(MiddleTimeLayout)
}

// UnmarshalJSON unmarshal string to JSONTime
func (t *JSONMiddleTime) UnmarshalJSON(b []byte) error {
	jt, err := t.doUnmarshalJSON(MiddleTimeLayout, b)
	if err != nil {
		return err
	}
	t.JSONTime = jt
	return nil
}

// JSONTime JSON 时间， 时间戳
type JSONTime int64

func (t JSONTime) doMarshalJSON(layout string) ([]byte, error) {
	if int64(t) <= 0 {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format(layout))), nil
}

// MarshalJSON outputs JSON presentation
func (t JSONTime) MarshalJSON() ([]byte, error) {
	return t.doMarshalJSON(LongTimeLayout)
}

func (t JSONTime) doUnmarshalJSON(layout string, b []byte) (JSONTime, error) {
	str := string(b[1 : len(b)-1])
	if str == "" {
		return t, nil
	}
	tm, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		return t, err
	}
	return JSONTime(tm.Unix()), nil
}

// UnmarshalJSON unmarshal string to JSONTime
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	jt, err := t.doUnmarshalJSON(LongTimeLayout, b)
	if err != nil {
		return nil
	}
	*t = jt
	return nil
}

// AsShortTime 转换为短日期格式
func (t JSONTime) AsShortTime() JSONShortTime { return JSONShortTime{t} }

// AsMiddleTime 转换为中长日期格式
func (t JSONTime) AsMiddleTime() JSONMiddleTime { return JSONMiddleTime{t} }

func (t JSONTime) String() string {
	return t.Format(LongTimeLayout)
}

func (t JSONTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t JSONTime) Add(d time.Duration) JSONTime {
	sec := d / time.Second
	return JSONTime(int64(t) + int64(sec))
}

func (t JSONTime) Format(layout string) string {
	return time.Unix(int64(t), 0).Format(layout)
}

// IsZero 是否为零值
func (t JSONTime) IsZero() bool { return t == 0 }

func Now() JSONTime { return JSONTime(time.Now().Unix()) }

func Today() JSONTime {
	y, m, d := time.Now().Date()
	return JSONTime(time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix())
}

// ParseJSONTime 解析时间戳
func ParseJSONTime(layout, v string) JSONTime {
	t, err := time.ParseInLocation(layout, v, time.Local)
	if err != nil {
		return 0
	}
	return JSONTime(t.Unix())
}

// ParseTime 解析时间
func ParseTime(layout, v string) time.Time {
	t, _ := time.ParseInLocation(layout, v, time.Local)
	return t
}
