package statusbar

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/verdverm/tview"
	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

const emptyMsg = "press 'Ctrl-<space>' to enter a command or '/path/to/something' to navigate"

type StatusBar struct {
	*tview.TextView

	curr    string   // current input (potentially partial)
	hIdx    int      // where we are in history
	history []string // command history
}

func New() *StatusBar {
	S := &StatusBar{
		history: []string{},
	}

	status := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	status.SetTitle(" Status ").SetTitleAlign(tview.AlignRight).SetBorder(true).SetBorderPadding(0, 0, 1, 0)

	S.TextView = status

	return S
}

func (S *StatusBar) Mount(context map[string]interface{}) error {
	vermui.AddWidgetHandler(S, "/sys/key/C-s", func(e events.Event) {
		S.SetBorderColor(tcell.ColorFuchsia)
		vermui.SetFocus(S.TextView)
	})
	S.SetDoneFunc(func(key tcell.Key) {
		S.SetBorderColor(tcell.ColorWhite)
		vermui.Unfocus()
	})

	vermui.AddWidgetHandler(S, "/user/error", func(evt events.Event) {
		str := fmt.Sprintf("[red]%v[white]", evt.Data.(*events.EventCustom).Data())

		S.Clear()
		fmt.Fprint(S, str)
		vermui.Draw()

		go func() {
			time.Sleep(time.Second * 6)
			text := S.GetText()
			if text == str {
				S.Clear()
				fmt.Fprint(S, "[lime]ok[white]")
				vermui.Draw()
			}
		}()
	})

	vermui.AddWidgetHandler(S, "/status/message", func(evt events.Event) {
		str := evt.Data.(*events.EventCustom).Data().(string)
		S.history = append(S.history, str)

		S.Clear()
		fmt.Fprint(S, str)
		vermui.Draw()

		go func() {
			time.Sleep(time.Second * 6)
			text := S.GetText()
			if text == str {
				S.Clear()
				fmt.Fprint(S, "[lime]ok[white]")
				vermui.Draw()
			}
		}()
	})

	return nil
}
func (S *StatusBar) Unmount() error {

	vermui.RemoveWidgetHandler(S, "/sys/key/C-s")
	vermui.RemoveWidgetHandler(S, "/user/error")
	vermui.RemoveWidgetHandler(S, "/status/message")

	return nil
}

// InputHandler returns the handler for this primitive.
func (S *StatusBar) InputHandler() func(tcell.Event, func(tview.Primitive)) {
	return S.WrapInputHandler(func(event tcell.Event, setFocus func(p tview.Primitive)) {
		switch evt := event.(type) {
		case *tcell.EventKey:

			dist := 1

			// Process key evt.
			switch key := evt.Key(); key {

			// Upwards, back in history
			case tcell.KeyHome:
				dist = len(S.history)
				fallthrough
			case tcell.KeyPgUp:
				dist += 4
				fallthrough
			case tcell.KeyUp: // Regular character.
				S.hIdx -= dist
				if S.hIdx < 0 {
					S.hIdx = 0
				}

			// Downwards, more recent in history
			case tcell.KeyEnd:
				dist = len(S.history)
				fallthrough
			case tcell.KeyPgDn:
				dist += 4
				fallthrough
			case tcell.KeyDown:
				S.hIdx += dist
				if S.hIdx >= len(S.history) {
					S.hIdx = len(S.history) - 1
				}

			}

		}
		str := ""
		if len(S.history) > 0 {
			str = S.history[S.hIdx]
		}
		S.Clear()
		fmt.Fprint(S, str)
		vermui.Draw()
	})
}
