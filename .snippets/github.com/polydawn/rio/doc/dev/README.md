developer readme
================

(This document is for someone hacking on the `rio` code itself.  Go up one dir for user-facing documentation!)



big picture
-----------

### Filesets and packed Wares

`rio` handles data in two catagorically different ways:

- Sets of files, unpacked on a real filesystem on your host
- Packed "wares", which represent sets of files and all their metadata, but are organized in (furious handwaving and abstractions) some other way.

Packed "wares" are content-addressable: `WareID`s are hashes.

Filesets are often cached by the content-addressable WareID they were unpacked from, but aren't exactly themselves considered content-addressable, because...

### universal IDs

... don't exist.  It's not a thing.  We give up.

Git hashes are content-addressable, and IPFS hashes are content-addressable, and they're not the same, and *that's okay*.

`rio.WareID` is a two-part object: it expresses the pack format, and the hash, *which is only defined within the namespace of that pack format*.

Example: `"git:asdf"` and `"ipfs:qwer"` are both valid WareIDs, and may even refer to the same files and metadata when unpacked.

*There is also no single hash which defines the WareID of an unpacked filesystem*.  This is just the other side of the same coin.
You can only compute a hash of a fileset by picking which packing format you want to use.

### storing and transfering packed Wares

In addition to this Fileset/Ware dichotomy, we have "Warehouses", which define a place to store Wares, and also implicitly some of their protocol.

#### mixing warehouses and pack formats

Warehouses tend to come in two types: things that grok key/value storage, and things that are... more complicated than that.

Key/value-style warehouses are often reusable.
Local filesystem dirs, S3 buckets, GCS buckets, and more can be used as k/v warehouses.
Transmats like 'tar' can use any of these k/v warehouses interchangably: just call the transmat with the appropriate URL style and hand over the relevant auth tokens.
You can mirror the exact same 'tar'-packed Ware between a local filesystem warehouse and an S3 warehouse, and it will keep the same WareID.

Other warehouses are more tightly tied to a specific pack format.
For example, git: git repositories only store git objects.  You have to be using git trees and commits as a packing format in order to use a git repo as a warehouse for your data.

Sometimes things can be hybridized even further: for example, you can use a 'tar' transmat to pack tarball files, then
turn around and place those files in a git commit!
Whether or not this is a good idea is a whole different question of course.



code layout
-----------

- `go.polydawn.net/rio` -- main package.  Interface definitions.  Other projects using `rio` as a library should import this package -- and few others.
- `go.polydawn.net/rio/fs` -- types for paths, and an abstract filesystem.  All `rio` code uses this when describing files and filesystem metadata.
  - `go.polydawn.net/rio/fs/osfs` -- a concrete implementation of the `fs` interfaces, implemented with a regular filesystem.
- `go.polydawn.net/rio/fsOp` -- operations on filesystems.  Distinct from `fs` because `fsOp` is more intention-oriented; `fs` is a fairly direct proxy to syscalls, and much less friendly.
- `go.polydawn.net/rio/warehouse/*` -- implementations of storage warehouses.  Local filesystem, S3, GCS, IPFS... each get their own package under here.
- `go.polydawn.net/transmat/*` -- implementations of filesystem packing formats.  E.g. `tar`.
  - REVIEW: so is this name a bug and the whole package should be `s/transmat/packing/`??  Probably
- `go.polydawn.net/transmat/mixins/fshash` -- helper functions for accumulating a hash for a fileset.  Used in some of the transmat implementations.
- `go.polydawn.net/lib/*` -- grabbag library functions; these are things that *probably* make sense even more broadly than rio, but are vendored here for simplicity's sake.

Overall, seen from the outside (as a consumer of `rio`-as-a-library):

- You want to handle Filesets.
  - You can handle them using the `fs` types for specifying paths
  - and you *may* use the `fs` types if necessary for inspecting or manipulating metadata...
    - but usually you want to use a transmat to handle filesets for you!
- Plug together a `packing` and a `warehouse` sytem to get a transmat.
  - for example `tar`+`s3`, or `tar`+`fs`: both construct valid transmats.
- And from the perspective of library consumers: that's it, no more details should need to leak.

### the fs vs fsOp split

Where 'fs' ends and 'fsOp' begins can be difficult to define, so here are some examples.

1. `Chown()` belongs in `fs`, because it's basically proxying a syscall.
2. `PlaceFile()` belongs in `fsOp` because it composes multiple syscalls...
3. More importantly, `PlaceFile()` belongs in `fsOp` because it implements some application-level sandboxing logic!  It has *opinions*.



when filesystems are insane^W fun
---------------------------------

The one place we're *not* offering infinte pluggable flexibility is around the basic file set and filesystem attributes themselves.
Here are some opinions we have:

- 'atime' and 'ctime' don't exist.
  - 'ctime' is impossible to set without a kernel hook or being a filesystem driver yourself; so it's out of bounds by virtue of being ridiculously inconvenient to work with.  Fortunately, it's also rarely used.
  - 'atime' properties are ridiculous to maintain, because nearly every operation in the universe may alter them anyway.  The behavior of 'atime' updates also varies *wildly* across modern filesystems, making it a very silly property indeed.  Fortunately, like with ctimes, nearly no applications in the wild look use atimes anyway (for the same reasons!).
- 'mtime' properties do exist.
  - Some programs do refer to them and make behavioral choices on mtime properties, so `rio` functions do honor mtime.
  - But do note `rio` transmats will default to flattening mtime properties them unless otherwise instructed, because they're usually useless noise.
    - So theologically speaking: `rio` transmats will honor Wares which are packed with noise included; but will try to reduce the amount of noise in the universe by not producing new Wares with noise included, unless you explicit ask for that.
- The order of entries in the filesystem doesn't matter.
  - In reality, on some filesystems, listing directory contents results in sorted order results; in others, it's mostly stable between read-only calls; in others yet, things may be even more random.  This is a bummer, and whenever relevant, `rio` will sort things to normalize behaviors.
- *Symlinks have mtimes.*
  - This is a controvertial one.  Mac filesystems are notable for disagreeing.  `rio` **errors** when given a symlink on Mac systems rather than having silently divergent behavior on Mac vs Linux systems.  `rio` prioritizes consistent precision and correctness over portability in this case.

These opinions can all be seen in the `fsOp` package, and the behaviors of transmats.
These opinions are *not* expressed in the `fs/*` packages, so if you need to reach in deeper, you can.
