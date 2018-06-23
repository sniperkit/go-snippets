package fs

import (
	"fmt"
	"io"
	"os"
	"syscall"

	. "github.com/warpfork/go-errcat"
)

type ErrorCategory string

const (
	/*
		The catch-all for errors we haven't particularly mapped.
	*/
	ErrMisc ErrorCategory = "fs-misc"

	ErrUnexpectedEOF ErrorCategory = "fs-unexpected-eof"
	ErrNotExists     ErrorCategory = "fs-not-exists"
	ErrAlreadyExists ErrorCategory = "fs-already-exists"
	ErrNotDir        ErrorCategory = "fs-not-dir"   // contextually, may be a form of either ErrNotExists or ErrAlready exists: shows up when e.g. lstat doesnotexist/deeper/path or mkdir aregularfile/deeper/path.
	ErrRecursion     ErrorCategory = "fs-recursion" // returned when cycles detected in symlinks.
	ErrShortWrite    ErrorCategory = "fs-shortwrite"
	ErrPermission    ErrorCategory = "fs-permission"

	/*
		Error returned when operating in a confined filesystem slice and an
		operation performed would result in effects outside the area, e.g.
		calling `PlaceFile` with "./reasonable/path" but "./reasonable" happens
		to be a symlink to "../../.." -- the symlink is valid, but to traverse
		it would violate confinement.

		Not all functions which do symlink checks will verify if the symlink target
		is within the operational area; they may return ErrBreakout upon encountering
		any symlink, even if following it would still be within bounds.
		Check the function's documentation for more info on how it behaves.

		Note that any function returning ErrBreakout is, by nature, doing so in a
		best-effort sense: if there are concurrent modifcations to the operational
		area of the filesystem by any other processes, it is *impossible* to
		avoid a TOCTOU violation.
	*/
	ErrBreakout ErrorCategory = "fs-breakout"
)

/*
   Categorize any errors into errcat-style errors with a category of `fs.ErrorCategory` type.
   Well-recognized errors will be normalized to a specific type,
   and all other errors will have a category of `ErrMisc`.
*/
func NormalizeIOError(ioe error) error {
	// Value switches.  Relatively fast -- and thus checked first.
	switch ioe {
	case nil:
		return nil
	case io.EOF, io.ErrUnexpectedEOF: // we don't believe in returning expected EOFs as errors.
		return Recategorize(ErrUnexpectedEOF, ioe)
	case io.ErrShortWrite:
		return Recategorize(ErrShortWrite, ioe)
	}
	// Complicated things there are no stdlib predicates for.
	switch e2 := ioe.(type) {
	case *os.PathError:
		switch e2.Err {
		case syscall.ENOTDIR:
			return ErrorDetailed(ErrNotDir, e2.Error(), map[string]string{"path": e2.Path})
		}
	}
	// Predicates.  God knows what they'll match;
	//  literally turing complete exhaustive checking is the only option.
	switch {
	case os.IsNotExist(ioe):
		switch e2 := ioe.(type) {
		case *os.PathError: // Rejoice in the one kind of error that provides a clear path.
			return ErrorDetailed(ErrNotExists, e2.Error(), map[string]string{"path": e2.Path})
		case *os.LinkError:
			return ErrorDetailed(ErrNotExists, e2.Error(), map[string]string{"pathOld": e2.Old, "pathNew": e2.New})
		case *os.SyscallError:
			return Recategorize(ErrNotExists, ioe) // has no path info :(
		default: // 'os.ErrExist' is stringly typed :(
			return Recategorize(ErrNotExists, ioe) // has no path info :(
		}
	case os.IsExist(ioe):
		return Recategorize(ErrAlreadyExists, ioe)
	case os.IsPermission(ioe):
		return Recategorize(ErrPermission, ioe)
	}
	// No matches.  Categorize to a placeholder.  At least it'll be serializable.
	return Recategorize(ErrMisc, ioe)
}

func NewBreakoutError(OpArea AbsolutePath, OpPath RelPath, LinkPath RelPath, LinkTarget string) error {
	return ErrorDetailed(
		ErrBreakout,
		fmt.Sprintf(
			"breakout error: refusing to traverse symlink at %q->%q while placing %q in %q",
			LinkPath, LinkTarget, OpPath, OpArea),
		map[string]string{
			"opArea":     OpArea.String(),
			"opPath":     OpPath.String(),
			"linkPath":   LinkPath.String(),
			"linkTarget": LinkTarget,
		},
	)
}
