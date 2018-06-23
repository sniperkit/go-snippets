/*
	The git transmat can unpack filesystems from the Git version control system.

	The features of this are intentionally limited for rio's purposes:
	the git transmat is read-only (unpack; no pack);
	`rio unpack git` must specify a hash (this should come as no surprise, since
	it's the rule for all Rio pack types, but it is different than git-checkout.
	Neither git branches nor tags are acceptable, being indirect and mutable);
	and the unpacked filesystem will *not* include the `.git` dir (because
	paradoxically, the internal layout of the `.git` objects is can be
	quite unpredictable, and is definitely not a function of the commit hash!).

	Packing into git is not supported because the semantics don't align:
	commit hashes are not a pure function of the packed file contents (due
	to additional info like commit timestamps and the parent commits),
	nor is it valid to store a single commit in git without a branch or tag name,
	and therefore the `api.PackFunc` signiture is almost totally incongruent.
	Git is designed for version control; not object storage.  This is okay.
*/
package git

import (
	"go.polydawn.net/go-timeless-api"
)

const PackType = api.PackType("git")
