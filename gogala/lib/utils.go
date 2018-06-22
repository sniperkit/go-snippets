package lib

import (
	"bytes"
	"time"
)

func MakeArgs(a ...interface{}) []interface{} {
	args := make([]interface{}, 0)
	for _, v := range a {
		args = append(args, v)
	}
	return args
}

func AppendString(a ...string) string {
	var b bytes.Buffer

	for i := 0; i < len(a); i++ {
		b.WriteString(a[i])
	}
	return b.String()
}

func PrintTimeStamp() string {
	return time.Now().Format(time.Kitchen)
}
