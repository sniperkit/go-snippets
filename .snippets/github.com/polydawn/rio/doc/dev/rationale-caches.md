Rationale: Caches
=================

Fileset cache implementations in rio are... *interesting*.
Long story short, they're pretty much an IPC convention all their own.
This sounds terrifying, but it's not without reason.

Here are those reasons.

### Why are is the filesystem being used to pass so much data?!

Well, that kind of answers itself, doesn't it?

We want one process to own the entire "get the files" workflow.

We want other processes to own the "stitch up a tree of filesets" workflow
(and then another one yet to handle launching containers, if that's what you're up to; etcetera).

Using paths on the filesystem to contain all the individual filesets and passing
around those paths as handles to the filesets is the only way to wield information
that may be in the gigabytes scale or larger.

### Why are cache behaviors baked in at such a low level?  Can't we extract them?

Components of the Timeless Stack are typically meant to be "batteries-included"
when working together, but split apart cleanly if used individually.
So, it's reasonable to ask: why can't the caching layer be cleanly split from
the rest of `rio` behaviors?

The answer is two-... er, *three*-fold:

1. It simply doesn't turn out as well as you'd think when you try to implement it that way.
   (See the "why aren't caches something I can layer *around* an UnpackFunc" section.)
2. Caches are the one implementation detail we have an absolute desire *not*
   to admit to the user, *ever*.  (See the "why are caches so hidden" section.)
3. *It's not what you want.*
   All of the user-facing (or even library-user facing) functionality that's reasonable
   to use is of the "put it here" form, not the "yield me a path" form.

### Why are caches so hidden?

*(Also can be stated as: Why can't I ask a rio command to show me the path to a cached fileset?)*

An overwhelming amount of systemic correctness depends on the correctness of caches.
It's imperative that users, no matter how well-meaning, do not muss with them.

The amount of the system that depends completely on the caches being unmolested
is pretty much 100%.  Thus, we don't want to mention the detail of their paths
to a user or caller *ever*, lest they feel encouraged to go mucking around in there.

### Why aren't caches something I can layer *around* an UnpackFunc?

Because it doesn't come out very elegantly.

You'd end up with some mess like this:

```
type FilesetCache interface {
	Yield(
		ctx context.Context,
		wareID api.WareID,
		filters api.FilesetFilters,
		warehouses []api.WarehouseAddr,
		monitor rio.Monitor,
	) (
		api.WareID,
		fs.FS,
		error,
	)
}
```

It's *almost* the signiture of an `UnpackFunc` all over again, but *slightly* different.
Just different enough that it doesn't compose at all.

Another possible approach is to punt on caching until the unpack-tree level,
and thus try to fold this ugly twist of API surface area under the covers of the
other API shifts there -- but this also doesn't shape up well:
Single unpack calls should be able to use (and optionally populate) a cache, no?

Finally, there's a small and almost comical yet still real issue with cache pools:
caches can only be shared by transmats of the same pack&hash type (!).
So, for each packType, a cache needs to mkdir "{packType}/committed".
What's the problem?  Well, only the unpack funcs themselves know their own pack type whitelist
(n.b., this also checks out from other angles of the bigger picture: someone on the
far side of the RPC gap doesn't know what ware types the rio command will support
unless it asks it!)... so, it's a shame if our cache layer does a mkdir for any ol'
string it gets in the type part of a WareID, just to find out it's not a real supported
type one function call later!

### You don't want to be handling caches at the higher level.

You just don't.

All the relevant user stories are better off with caching, but a PITA path full of
quicksand to implement if you think you can do the caching and composing at a higher level.
You want to stick to the ways that rio can be invoked already, none of which have
any reason to admit to the internal cache path.

Ways Rio can be invoked:

1. `rio unpack` -- a user wants a single thing in a single place on their host filesystem.
   This could either ignore cache, or use cache and copy-placer.
   Or, in poweruser mode one might even use a `--placer=bind` option.
2. `rio unpack-tree` -- a user wants to unpack a whole set of stuff onto their host filesystem.
   Mostly, this is like a for-loop over `rio unpack`.  **Mostly**.
   There is an **asterisk** here; wait for it.
3. `rio unpack-tree --faux-chroot` -- repeatr wants to unpack a whole set of stuff.
   Why is this different than when a user asks for it?  Here's that **asterisk** again.

Tree-unpacks are actually different than just a rack of unpacks in a loop:
stiching them together should verify that symlink breakouts aren't an issue.
This is the source of the above *asterisks*.

When a user issues a single unpack command, it's free to follow symlinks in the path
the user gave to the unpack target -- anywhere; presumably the user meant whatever they said,
and knows their host filesystem and its symlinks reasonably well.

When a *tree* of filesets is being unpacked and stitched together, it's a different story:
it's expected that all the filesystem changes will be under the root path given to the
unpack-tree command (!), and nowhere else (!!!).
This means the unpack-tree process needs to keep checking for each unpack that its
place in the tree won't be shifted unreasonably by symlinks created by earlier unpacks:
what if one of the earlier unpacks creates a symlink pointing somewhere outside the
unpack-tree root path?

(As another fun layer of detail, unpack-tree operations can come in two modes:
an unpack-tree that's expected to leave artifacts within a subdirectory of a bigger
filesystem should reject traversal of breakout'ing symlinks entirely;
whereas for an unpack-tree setting up a container filesystem for repeatr can
reasonably massage any absolutized or overly-up-up symlinks and manually "traverse"
them in the same way they would be interpreted if already chroot'd, because
the result will make sense when seen inside the container.)

Main thing to learn here: you *never* want a higher-level caller doing stitch
themselves.  It's actually quite error-prone.  You want to let rio own it for you.

### So unite all these justifications -- you did *what* in reaction to this?!

Caches as a shared filesystem IPC interface behavioral specification.

(God have mercy.  But it's just an Admission of the True.)

We implement several rio features to quietly share the same understanding of how
caches will be laid out on the filesystem.

Even when one of these components is invoking the others (e.g., an unpack-tree
invoking individual unpacks) using the typical external-process RPC mechanism,
they're *also* communicating via the filesystem,
almost as if the internal side-effecting filesystem behaviors are part of the API's contract.
(`s/almost//` -- it is.)

Specifically: when `rio unpack-tree` invokes `rio unpack`, it does so with
a path of `"-"` and a flag `--placer="none"`.
The unpack-tree command is expecting the unpack command to populate the cache,
and then due to the "none" placer instruction, simply exit.
The unpack-tree command will then do all of the placer behaviors after all
individual unpacks are done.
The unpack commands do not send any message to their parent unpack-tree command
about the cache paths, nor does the unpack-tree command deign to explicitly command
them about what cache paths to use; they simply share a tacit understanding that
they will all be using the same essential parts of path naming conventions.

**From the "caches should be deeply hidden" reasoning:**
It actually seems *safer* to design a quietly shared understanding on the
filesystem between the unpack and unpack-tree tools -- both of which are already
required to be consenting adults in this picture -- and expose *nothing* in the
CLI or RPC API... than it does to make an "advanced" mode RPC API which does
show such details as the cache paths.
People will be much more scared of using internal details like a cache
filesystem behavior spec in a far-away internal package than they will be of
a string returned from a published RPC API func -- *as they should be*.

**From the "it doesn't layer around unpack elegantly" reasoning:**
Filling the cache is handled in the logic for each unpack func.
(They have some utility code to mix in which reduces the redundancy of this, of course.)
This works well.

**From the "it's what you want" reasoning:**
All the CLI commands and RPC APIs we publish are still primarily oriented around the
user saying "put it {here}" -- which is as it should be.
Some esoteric combinations of flags which an end-user would be unlikely to have any
direct use for together are used by the internals of unpack-tree -- but none of them
actually cause drastic or inconsistent changes in behavior, which is parsimonious.

### There you have it.

That's why caches work the way they do.
