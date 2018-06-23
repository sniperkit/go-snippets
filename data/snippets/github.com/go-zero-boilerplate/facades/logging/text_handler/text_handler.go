package text_handler

//TODO: Code copied and tweaked from https://github.com/apex/log/blob/master/handlers/text/text.go

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	apex "github.com/francoishill/log"

	"github.com/fatih/color"
)

var (
	//DefaultTimeStampFormat is the default timestamp format for the Handler
	DefaultTimeStampFormat = "2006-01-02 15:04:05"
	//DefaultMessageWidth is the default message width for the Handler
	DefaultMessageWidth = 50
)

//New text handler. Empty `timeStampFormat` will default to `DefaultTimeStampFormat`. Zero `messageWidth` will default to `DefaultMessageWidth`
func New(infoWriter, errorWriter io.Writer, timeStampFormat string, messageWidth int) *Handler {
	if timeStampFormat == "" {
		timeStampFormat = DefaultTimeStampFormat
	}
	if messageWidth <= 0 {
		messageWidth = DefaultMessageWidth
	}
	return &Handler{
		timeStampFormat: timeStampFormat,
		messageWidth:    messageWidth,
		InfoWriter:      infoWriter,
		ErrorWriter:     errorWriter,
	}
}

// Colors mapping.
var Colors = [...]*color.Color{
	apex.DebugLevel:     color.New(color.FgWhite),
	apex.InfoLevel:      color.New(color.FgBlue),
	apex.NoticeLevel:    color.New(color.FgYellow),
	apex.WarnLevel:      color.New(color.FgYellow),
	apex.ErrorLevel:     color.New(color.FgRed),
	apex.CriticalLevel:  color.New(color.FgHiRed),
	apex.AlertLevel:     color.New(color.FgHiRed),
	apex.EmergencyLevel: color.New(color.FgHiRed),
}

// Strings mapping.
var Strings = [...]string{
	apex.DebugLevel:     "DEBUG",
	apex.InfoLevel:      "INFO",
	apex.NoticeLevel:    "NOTICE",
	apex.WarnLevel:      "WARN",
	apex.ErrorLevel:     "ERROR",
	apex.CriticalLevel:  "CRIT",
	apex.AlertLevel:     "ALERT",
	apex.EmergencyLevel: "EMER",
}

// field used for sorting.
type field struct {
	Name  string
	Value interface{}
}

// by sorts projects by call count.
type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Handler implementation.
type Handler struct {
	mu sync.RWMutex

	timeStampFormat string
	messageWidth    int

	InfoWriter  io.Writer
	ErrorWriter io.Writer
}

func (h *Handler) pickWriterFromLevel(level apex.Level) io.Writer {
	if level <= apex.InfoLevel {
		return h.InfoWriter
	}
	return h.ErrorWriter
}

// HandleLog implements apex.Handler.
func (h *Handler) HandleLog(e *apex.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	color := Colors[e.Level]
	colorFunc := color.SprintfFunc()
	level := Strings[e.Level]

	var fields []field

	for k, v := range e.Fields {
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))

	writer := h.pickWriterFromLevel(e.Level)
	formattedTime := time.Now().Format(h.timeStampFormat)
	messageWidthStr := fmt.Sprintf("%d", h.messageWidth)
	printFmt := "[%s] %s %-" + messageWidthStr + "s"
	fmt.Fprintf(writer, "%s", fmt.Sprintf(printFmt, formattedTime, colorFunc("%-6s", level), e.Message))

	if len(fields) > 0 {
		fieldKeyVals := []string{}
		for _, f := range fields {
			valueExpanded := strings.Replace(fmt.Sprintf("%+v", f.Value), "\n", "\\n", -1)
			fieldKeyVals = append(fieldKeyVals, fmt.Sprintf("%s=%s", colorFunc("%s", f.Name), valueExpanded))
		}
		fmt.Fprintf(writer, " %s", strings.Join(fieldKeyVals, " "))
	}

	fmt.Fprintln(writer)

	return nil
}
