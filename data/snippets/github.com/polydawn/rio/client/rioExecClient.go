package rioexecclient

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
)

var (
	_ rio.UnpackFunc = UnpackFunc
	_ rio.PackFunc   = PackFunc
	_ rio.MirrorFunc = MirrorFunc
)

func UnpackFunc(
	ctx context.Context,
	wareID api.WareID,
	path string,
	filters api.FilesetFilters,
	placementMode rio.PlacementMode,
	warehouses []api.WarehouseAddr,
	monitor rio.Monitor,
) (gotWareID api.WareID, err error) {
	// Marshal args.
	args, err := UnpackArgs(wareID, path, filters, placementMode, warehouses, monitor)
	if err != nil {
		return api.WareID{}, err
	}
	// Bulk of invoking and handling process messages is shared code.
	return packOrUnpack(ctx, args, monitor)
}

func PackFunc(
	ctx context.Context,
	packType api.PackType,
	path string,
	filters api.FilesetFilters,
	warehouse api.WarehouseAddr,
	monitor rio.Monitor,
) (api.WareID, error) {
	// Marshal args.
	args, err := PackArgs(packType, path, filters, warehouse, monitor)
	if err != nil {
		return api.WareID{}, err
	}
	// Bulk of invoking and handling process messages is shared code.
	return packOrUnpack(ctx, args, monitor)
}

// internal implementation of message parsing for both pack and unpack.
// (they "conincidentally" have the same API.)
func packOrUnpack(
	ctx context.Context,
	args []string,
	monitor rio.Monitor,
) (api.WareID, error) {
	if monitor.Chan != nil {
		defer close(monitor.Chan)
	}

	// Spawn process.
	cmd := exec.Command("rio", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: failed to start: %s", err)
	}
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	if err = cmd.Start(); err != nil {
		return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: failed to start: %s", err)
	}

	// Set up reaction to ctx.done: send a sig to the child proc.
	//  (No, you couldn't set this up without a goroutine -- you can't select with the IO we're about to do;
	//  and No, you couldn't do it until after cmd.Start -- the Process handle doesn't exist until then.)
	go func() {
		<-ctx.Done() // FIXME goroutine leak occurs when the process ends gracefully
		cmd.Process.Signal(os.Interrupt)
		time.Sleep(100 * time.Millisecond)
		cmd.Process.Signal(os.Kill)
	}()

	// Consume stdout, converting it to Monitor.Chan sends.
	//  When exiting because the child sent its 'result' message correctly, the
	//  msgSlot will hold the final data (or error); we'll return it at the end,
	//  but we also check the exit code for a match.
	//  (We're relying on the child proc getting signal'd to close the stdout pipe
	//  and in turn release us here in case of ctx.done.)
	unmarshaller := refmt.NewUnmarshallerAtlased(json.DecodeOptions{}, stdout, rio.Atlas)
	var msgSlot rio.Event
	for {
		// Peel off a message.
		if err := unmarshaller.Unmarshal(&msgSlot); err != nil {
			if err == io.EOF {
				// In case of unexpected EOF, there must have been a panic on the other side;
				//  it'll be more informative to break here and return the error from Wait,
				//  which will include the stderr capture.
				break
			}
			return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: API parse error: %s", err)
		}

		// If it's the final "result" message, prepare to return.
		if msgSlot.Result != nil {
			// Bail.  We'll review this last message frame in a second.
			break
		}
		// For all other messages: forward to the monitor channel (if it exists!)
		if monitor.Chan != nil {
			select {
			case <-ctx.Done():
				break
			case monitor.Chan <- msgSlot:
				// continue
			}
		}
	}

	// Wait for process complete.
	//  The exit code SHOULD be redundant with the error we SHOULD have already
	//  deserialized... but we check that it all matches up.
	code, err := waitFor(cmd)
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: wait error: %s (stderr: %q)", err, stderrBuf.String())
	}
	if code == 0 {
		// If the exit code was success, we'd sure better have gotten the rightly formatted result message.
		if msgSlot.Result == nil {
			return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: exited zero, but no clear result?! (stderr: %q)", err, stderrBuf.String())
		}
		if msgSlot.Result.Error != nil {
			return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: exited zero, but result had error, category=%s: %s", msgSlot.Result.Error.Category, msgSlot.Result.Error)
		}
		return msgSlot.Result.WareID, nil // This is the happy path return!
	}
	// For non-zero exits: Check match for sanity.
	exitCategory := rio.CategoryForExitCode(code)
	if msgSlot.Result == nil {
		return api.WareID{}, Errorf(exitCategory, "no message available (stderr: %q)", stderrBuf.String())
	}
	if msgSlot.Result.Error.Category() != exitCategory {
		return api.WareID{}, Errorf(exitCategory, "no message available (stderr: %q)", stderrBuf.String())
	}
	return api.WareID{}, msgSlot.Result.Error // This is the clean error path!
}

func MirrorFunc(
	ctx context.Context,
	wareID api.WareID,
	target api.WarehouseAddr,
	sources []api.WarehouseAddr,
	monitor rio.Monitor,
) (api.WareID, error) {
	// TODO all of it
	return api.WareID{}, nil
}
