package debugging

import (
	"bytes"
	"fmt"
	"runtime"
)

//GetNormalStack gets the normal stack trace
func GetNormalStack(getAllGoRoutines bool) string {
	stackBuf := make([]byte, 1<<16)
	runtime.Stack(stackBuf, getAllGoRoutines)
	stackBuf = bytes.Trim(stackBuf, "\x00")
	return string(stackBuf)
}

//GetPrettyStack gets the pretty stack trace
func GetPrettyStack() string {
	var buf bytes.Buffer
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s:%d", file, line)))
	}
	return buf.String()
}
