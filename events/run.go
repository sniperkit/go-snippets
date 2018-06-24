package events

import (
	"sync"

	"github.com/verdverm/tview"
)

var defaultEventStream *EventStream
var defaultWgtMgr WgtMgr
var wgtMgrMuxtx sync.Mutex

var DefaultHandler = func(e Event) {}

var sysEvtChs []chan Event
var customEventCh = make(chan Event, 256)

func Init(app *tview.Application) error {
	sysEvtChs = make([]chan Event, 0)
	go hookEventsFromApp(app)

	defaultEventStream = NewEventStream()
	defaultEventStream.Init()
	defaultEventStream.Merge("tcell", NewSysEvtCh())
	defaultEventStream.Merge("custom", customEventCh)

	defaultWgtMgr = NewWgtMgr()
	defaultEventStream.Hook(defaultWgtMgr.WgtHandlersHook())

	return nil
}

func Start() error {
	defaultEventStream.Loop()
	return nil
}

func Stop() error {
	defaultEventStream.StopLoop()
	return nil
}

func Merge(name string, ec chan Event) {
	defaultEventStream.Merge(name, ec)
}

func AddGlobalHandler(path string, handler func(Event)) {
	defaultEventStream.Handle(path, handler)
}

func RemoveGlobalHandler(path string) {
	defaultEventStream.RemoveHandle(path)
}

func ClearGlobalHandlers() {
	defaultEventStream.ResetHandlers()
}

func AddWidgetHandler(wgt tview.Primitive, path string, handler func(Event)) {
	if _, ok := defaultWgtMgr[wgt.Id()]; !ok {
		defaultWgtMgr.AddWgt(wgt)
	}

	defaultWgtMgr.AddWgtHandler(wgt.Id(), path, handler)
}

func RemoveWidgetHandler(wgt tview.Primitive, path string) {
	_, ok := defaultWgtMgr[wgt.Id()]
	if !ok {
		return
	}

	defaultWgtMgr.RmWgtHandler(wgt.Id(), path)
}

func ClearWidgetHandlers(wgt tview.Primitive) {
	_, ok := defaultWgtMgr[wgt.Id()]
	if !ok {
		return
	}

	defaultWgtMgr.ClearWgtHandlers(wgt.Id())
}
