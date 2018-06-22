package testutil

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/rio/caps"
)

type ConveyRequirement struct {
	Name      string
	Predicate func() bool
}

/*
	Require that the tests are not running with the "short" flag enabled.
*/
var RequiresLongRun = ConveyRequirement{"run long tests", func() bool { return !testing.Short() }}

/*
	Require that the test process is running with enough capabilities to be able to manage file ownership.
*/
var RequiresCanManageOwnership = ConveyRequirement{"have caps for managing file ownership", caps.Scan().CanManageOwnership}

/*
	Require that the test process is running with enough capabilities to be able to make bind mounts.
*/
var RequiresCanMountBind = ConveyRequirement{"have caps for mounting binds", caps.Scan().CanMountBind}

/*
	Require that the test process is running with enough capabilities to be able to make any/all mounts.
*/
var RequiresCanMountAny = ConveyRequirement{"have caps for any mounting", caps.Scan().CanMountAny}

/*
	Require than an env var *not* be set.

	We use this for things like `RequiresEnvBlank(RIO_TEST_SKIP_AUFS)`.
*/
func RequiresEnvBlank(key string) ConveyRequirement {
	return ConveyRequirement{
		fmt.Sprintf("env %q must not be set", key),
		func() bool { return os.Getenv(key) == "" },
	}
}

/*
	Decorates a GoConvey test to check a set of `ConveyRequirement`s,
	returning a dummy test func that skips (with an explanation!) if any
	of the requirements are unsatisfied; if all is well, it yields
	the real test function unchanged.  Provide the `...ConveyRequirement`s
	first, followed by the `func()` (like the argument order in `Convey`).
*/
func Requires(items ...interface{}) func(c convey.C) {
	// parse args
	// not the most robust parsing.  just panics if there's weird stuff
	var requirements []ConveyRequirement
	for _, it := range items {
		if req, ok := it.(ConveyRequirement); ok {
			requirements = append(requirements, req)
		} else {
			break
		}
	}
	action := items[len(items)-1]
	// examine requirements
	var widest int
	for _, req := range requirements {
		if len(req.Name) > widest {
			widest = len(req.Name)
		}
	}
	// check requirements
	var requirementsListing bytes.Buffer
	var names []string
	allSat := true
	for _, req := range requirements {
		sat := req.Predicate()
		allSat = allSat && sat
		names = append(names, req.Name)
		fmt.Fprintf(&requirementsListing, "requirement %*q: %v\n", widest+2, req.Name, sat)
	}
	// act
	if allSat {
		return func(c convey.C) {
			// attempted: inserting another convey that makes a single 'true=true' assertion so we see the prereqs and a green check mark.
			// doesn't work: doing so causes a leaf node, in which everything is run :/ even if skipped, the remaining `So` that aren't
			// in another block get attached to it, which makes verrry odd reading, and causes an extra repetition of anything
			// that isn't in another convey block.
			//	convey.SkipConvey(title, func() { convey.So(true, convey.ShouldBeTrue) })
			switch action := action.(type) {
			case func():
				action()
			case func(c convey.C):
				action(c)
			}
		}
	} else {
		title := "Prereqs: " + strings.Join(names, ", ")
		return func(c convey.C) {
			convey.Convey(title, nil)
			c.Println()
			c.Print(requirementsListing.String())
		}
	}
}
