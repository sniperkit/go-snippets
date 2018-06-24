# Overview
[![Build Status](https://travis-ci.org/git-hook/captain-hook.svg?branch=master)](https://travis-ci.org/git-hook/captain-hook)

This package provides a framework for git hooks following the [*.d configuration approach].

In other words, captain-hook allows us to express each git hook as a
*pipeline* instead of a single script.

[*.d configuration approach]: http://blog.siphos.be/2013/05/the-linux-d-approach/

# Installation

Add captain-hook to your repository as a git submodule:

    mkdir -p hooks
    git submodule add https://github.com/ericcrosson/captain-hook hooks/captain-hook
    hooks/captain-hook/setup

Invoke `setup` to create the following directory structure:

    ▲ tree -I captain-hook
    .
    └── hooks
        ├── commit-msg -> captain-hook/run-hook
        ├── commit-msg.d
        ├── post-checkout -> captain-hook/run-hook
        ├── post-checkout.d
        ├── post-clone -> captain-hook/run-hook
        ├── post-clone.d
        ├── post-commit -> captain-hook/run-hook
        ├── post-commit.d
        ├── pre-commit -> captain-hook/run-hook
        ├── pre-commit.d
        ├── pre-rebase -> captain-hook/run-hook
        ├── pre-rebase.d
        ├── prepare-commit-msg -> captain-hook/run-hook
        └── prepare-commit-msg.d

# Use

Every client-side git hook is installed, and will be called by git
naturally.  Each git hook is intercepted by captain-hook, a wrapper to
invoke the scripts in its corresponding `<git-hook-name>.d` directory.
Your job is to define the git-hook pipelines by populating these
directories with carefully-named scripts.

## Naming scripts

Captain-hook invokes scripts in the order specified by the names of
the scripts in the relevant hook-directory.  The scripts are invoked
in numerical order, from 000-999. Note that scripts beginning with
anything other than a three-digit number between one and nine-hundred
ninety-nine, inclusive, will be ignored.  It is advisable to follow
this outlined hierarchy for naming scripts:

000-099 System level initialization scripts <br>
100-899 User level git hook actions <br>
900-999 System level teardown scripts <br>

In order to fully leverage the benefits of this numeric system,
consider allowing distance between filenames of discrete parts of a
single git-hook pipeline.  For example, the pre-push pipeline may
build and test, but rather than naming your scripts 100-build.sh and
101-test.sh, leave room for future expansion: 110-build.sh and
120-test.sh.  This ensures minimal changes to (proven) working parts
of the system as well as isolation in git diffing tools to make code
review easier.

## Executing scripts

If one of the scripts exits with a non-zero code, the remaining
scripts are not invoked in order to limit unexpected behavior.
Sometimes this is undesirable, as we would like to separate the
cleanup code from the actionable code that created the mess.  For this
reason captain-hook provides the concept of an *ensured script* that
is run every time git invokes a hook, no matter how successful the
hook execution has been.  Typically, ensured scripts are the last
scripts to run, performing cleanup or unstashes to undo modifications
to the repository made by earlier scripts in the git hook.

Sometimes it is desirable to "short-circuit" a git-hook similar to the
above, but without reporting an error.  For example, perhaps the first
script in a pipeline is to check if any action needs to be taken, and
it concludes the answer is "no."  In this case, return a negative
exit-code to skip the remaining non-ensured scripts in the hook
directory.  Remember, ensured scripts are run *every* time the hook is
invoked by git.

# Supported platforms

This project has been tested on:

- [X] GNU/Linux

- [X] macOS

# Extras

Astute readers will question the concept of a post-clone script.  Note
that this technology
is [now possible](https://github.com/git-hook/post-clone).

# License

This project is protected by
the [Apache license](https://www.apache.org/licenses/LICENSE-2.0).
