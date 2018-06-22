/*
	Helper functions for emitting structured logs to the rio.Monitor.

	These functions encompass most common lifecycle events in a transmat,
	and using them A) saves typing and B) keeps the common stuff formatted
	in a common way between transmats.
	Transmats can of course also write their own log events raw; it is freetext.
*/
package log

import (
	"fmt"
	"time"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
)

func CacheHasIt(mon rio.Monitor, ware api.WareID) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogInfo,
			Msg:   fmt.Sprintf("cache already has ware %q", ware),
			Detail: [][2]string{
				{"wareID", ware.String()},
			},
		},
	}
}

// Log path for a 'rio.ErrWarehouseUnavailable'; mode is "read" or "write".
func WarehouseUnavailable(mon rio.Monitor, err error, wh api.WarehouseAddr, ware api.WareID, mode string) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogWarn,
			Msg:   fmt.Sprintf("%s while dialing warehouse %q for %s: %s", rio.ErrWarehouseUnavailable, wh, mode, err),
			Detail: [][2]string{
				{"warehouse", string(wh)},
				{"wareID", ware.String()},
				{"error", err.Error()},
			},
		},
	}
}

// Log path for a 'rio.ErrWareNotFound'.
func WareNotFound(mon rio.Monitor, err error, wh api.WarehouseAddr, ware api.WareID) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogInfo,
			Msg:   fmt.Sprintf("%s from warehouse %q for ware %q", rio.ErrWareNotFound, wh, ware),
			Detail: [][2]string{
				{"warehouse", string(wh)},
				{"wareID", ware.String()},
				{"error", err.Error()},
			},
		},
	}
}

func WareReaderOpened(mon rio.Monitor, wh api.WarehouseAddr, ware api.WareID) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogInfo,
			Msg:   fmt.Sprintf("read for ware %q opened from warehouse %q", ware, wh),
			Detail: [][2]string{
				{"warehouse", string(wh)},
				{"wareID", ware.String()},
			},
		},
	}
}

// This logs a cache hit where the "object store" (as git calls it, for example)
// has the object we need -- as opposed to our fileset cache, which presumably
// has already missed, or we would've returned that already.
// It means we *aren't* doing network ops, but an unpacking still needs to run.
func WareObjCacheHit(mon rio.Monitor, ware api.WareID) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogInfo,
			Msg:   fmt.Sprintf("raw objects for ware %q found cached, unpacking them to fileset", ware),
			Detail: [][2]string{
				{"wareID", ware.String()},
			},
		},
	}
}

func MirrorNoop(mon rio.Monitor, wh api.WarehouseAddr, ware api.WareID) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogInfo,
			Msg:   fmt.Sprintf("mirror skip: warehouse at %q already has ware %q", wh, ware),
			Detail: [][2]string{
				{"warehouse", string(wh)},
				{"wareID", ware.String()},
			},
		},
	}
}

// Emit debug log entry for implicit parent dir creation.
// This is mostly a tar thing and probably shouldn't be in the general mixins;
// the fact that it's here is a hint that we need some serious refactor on logs.
//
// ALSO a FIXME: we would like to comment on the wareID here, but in calling contexts,
// that's *not actually topically relevant*... we need logger helpers that handle this.
func DirectoryInferred(mon rio.Monitor, inferred, path fs.RelPath) {
	if mon.Chan == nil {
		return
	}
	mon.Chan <- rio.Event{
		Log: &rio.Event_Log{
			Time:  time.Now(),
			Level: rio.LogDebug,
			Msg:   fmt.Sprintf("unpacking: inferring dir %q for parent of path %q", inferred, path),
			Detail: [][2]string{
				{"inferred", inferred.String()},
				{"path", path.String()},
			},
		},
	}
}
