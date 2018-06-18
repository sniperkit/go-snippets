package interrupts

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/uoregon-libraries/gopkg/logger"
)

var isDone int32

// Logger is set to the default logger, but can be overridden for custom behavior
var Logger = logger.DefaultLogger

// TrapIntTerm catches interrupt and termination signals to let processes exit
// more cleanly.  A second signal of either type will immediately end the
// process, however.
func TrapIntTerm(quit func()) {
	var sigInt = make(chan os.Signal, 1)
	signal.Notify(sigInt, syscall.SIGINT)
	signal.Notify(sigInt, syscall.SIGTERM)
	go func() {
		for range sigInt {
			if done() {
				Logger.Warnf("Force-interrupt detected; shutting down.")
				os.Exit(1)
			}

			Logger.Infof("Interrupt detected; attempting to clean up.  Another signal will immediately end the process.")
			atomic.StoreInt32(&isDone, 1)
			quit()
		}
	}()
}

func done() bool {
	return atomic.LoadInt32(&isDone) == 1
}
