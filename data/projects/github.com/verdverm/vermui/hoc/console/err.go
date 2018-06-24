package console

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/verdverm/tview"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

type ErrConsoleWidget struct {
	*tview.TextView
}

func NewErrConsoleWidget() *ErrConsoleWidget {
	textView := tview.NewTextView()
	textView.
		SetTextColor(tcell.ColorMaroon).
		SetScrollable(true).
		SetChangedFunc(func() {
			vermui.Draw()
			textView.ScrollToEnd()
		})

	textView.SetTitle(" errors ").
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	C := &ErrConsoleWidget{
		TextView: textView,
	}

	return C
}

func (C *ErrConsoleWidget) Mount(context map[string]interface{}) error {

	vermui.AddGlobalHandler("/user/error", func(evt events.Event) {
		str := evt.Data.(*events.EventCustom).Data()
		text := fmt.Sprintf("[%s] %v\n", evt.When().Format("2006-01-02 15:04:05"), str)
		fmt.Fprintf(C, "%s", text)
	})

	vermui.AddGlobalHandler("/sys/err", func(ev events.Event) {
		err := ev.Data.(*events.EventError)
		line := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), err)
		fmt.Fprintf(C, "[red]SYSERR %v[white]\n", line)
	})

	return nil
}
func (C *ErrConsoleWidget) Unmount() error {
	vermui.RemoveWidgetHandler(C, "/user/error")
	vermui.RemoveWidgetHandler(C, "/sys/err")
	return nil
}
