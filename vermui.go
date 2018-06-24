package vermui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/maruel/panicparse/stack"
	"github.com/verdverm/tview"

	"github.com/verdverm/vermui/events"
)

var app *tview.Application
var appLock sync.RWMutex
var rootView tview.Primitive

// Init initializes vermui library. This function should be called before any others.
// After initialization, the library must be finalized by 'Close' function.
func Init() error {
	app = tview.NewApplication()

	err := events.Init(app)
	if err != nil {
		return err
	}

	events.AddGlobalHandler("/", events.DefaultHandler)

	events.AddGlobalHandler("/sys/redraw", func(e events.Event) {
		Draw()
	})

	return nil
}

// blocking call
func Start() error {

	// catch panics, clean up, format error
	defer func() {
		e := recover()
		if e != nil {
			Stop()
			// Print a formatted panic output
			fmt.Fprintf(os.Stderr, "Captured a panic(value=%v) lib.Start()... Exit vermui and clean terminal...\nPrint stack trace:\n\n", e)
			//debug.PrintStack()
			gs, err := stack.ParseDump(bytes.NewReader(debug.Stack()), os.Stderr)
			if err != nil {
				debug.PrintStack()
				os.Exit(1)
			}
			p := &stack.Palette{}
			buckets := stack.SortBuckets(stack.Bucketize(gs, stack.AnyValue))
			srcLen, pkgLen := stack.CalcLengths(buckets, false)
			for _, bucket := range buckets {
				io.WriteString(os.Stdout, p.BucketHeader(&bucket, false, len(buckets) > 1))
				io.WriteString(os.Stdout, p.StackLines(&bucket.Signature, srcLen, pkgLen, false))
			}
			panic(e)
		}
	}()

	// start the event engine
	go events.Start()

	err := rootView.Mount(nil)
	if err != nil {
		panic(err)
	}

	// blocking
	app.SetRoot(rootView, true)
	return app.Run()
}

// Close finalizes vermui library,
// should be called after successful initialization when vermui's functionality isn't required anymore.
func Stop() error {
	//appLock.Lock()
	//defer appLock.Unlock()

	app.Stop()
	err := events.Stop()
	if err != nil {
		return err
	}
	return nil
}

func Application() *tview.Application {
	return app
}

func Draw() {
	appLock.RLock()
	defer appLock.RUnlock()

	go app.Draw()
}

func Clear() {
	//appLock.Lock()
	//defer appLock.Unlock()

	if app == nil {
		// really shouldn't get here, but the event stream is still running
		return
	}
	screen := app.Screen()
	if screen != nil {
		screen.Clear()
		screen.Sync()
	}
}

func GetRootView() tview.Primitive {
	return rootView
}

func SetRootView(v tview.Primitive) {
	rootView = v
}

func GetFocus() (p tview.Primitive) {
	//appLock.RLock()
	//defer appLock.RLock()

	return app.GetFocus()
}

func SetFocus(p tview.Primitive) {
	//appLock.Lock()
	//defer appLock.Unlock()

	if app == nil {
		// really shouldn't get here, but the event stream is still running
		return
	}

	// go app.Screen().HideCursor()
	app.SetFocus(p)
	Draw()
}
func Unfocus() {
	//appLock.Lock()
	//defer appLock.Unlock()

	if app == nil {
		// really shouldn't get here, but the event stream is still running
		return
	}

	// go app.Screen().HideCursor()
	app.SetFocus(rootView)
	Draw()
}

func AddGlobalHandler(path string, handler func(events.Event)) {
	events.AddGlobalHandler(path, handler)
}

func RemoveGlobalHandler(path string) {
	events.RemoveGlobalHandler(path)
}

func ClearGlobalHandlers() {
	events.ClearGlobalHandlers()
}

func AddWidgetHandler(widget tview.Primitive, path string, handler func(events.Event)) {
	events.AddWidgetHandler(widget, path, handler)
}

func RemoveWidgetHandler(widget tview.Primitive, path string) {
	events.RemoveWidgetHandler(widget, path)
}

func ClearWidgetHandlers(widget tview.Primitive) {
	events.ClearWidgetHandlers(widget)
}
