# repotool - a tool to aggregate source code repositories metadata

[![Build Status](https://travis-ci.org/DevMine/repotool.png?branch=master)](https://travis-ci.org/DevMine/repotool)
[![GoDoc](http://godoc.org/github.com/DevMine/repotool?status.svg)](http://godoc.org/github.com/DevMine/repotool)
[![GoWalker](http://img.shields.io/badge/doc-gowalker-blue.svg?style=flat)](https://gowalker.org/github.com/DevMine/repotool)

`repotool` is a command line tool that aggregates source code repositories
metadata (such as VCS type, commits and so on) and produces JSON objects out of
it.  It is also able to store repository information into a database.

A repository contains a list of commits that may contain, if you enable this
option, a list of deltas. A delta contains information about a file
touched by a commit. It may also contain patches if specified via an option.

Currently, only [git](http://git-scm.com/) is supported.

Below is an example of the data produced, without commit deltas and patches:

```
{
  "name": "repotool",
  "vcs": "git",
  "clone_url": "https://github.com/DevMine/repotool.git",
  "clone_path": "/home/robin/Hacking/repotool",
  "default_branch": "master",
  "commits": [
    {
      "vcs_id": "df55def5e6185447c6bd360ec1144a847d73b986",
      "message": "repotool: Add possibility to insert commit diff deltas into the db.\n\nFor this purpose, create a new 'commit_diff_deltas' table.\n",
      "author": {
        "name": "Robin Hahling",
        "email": "robin.hahling@gw-computing.net"
      },
      "committer": {
        "name": "Robin Hahling",
        "email": "robin.hahling@gw-computing.net"
      },
      "author_date": "2015-01-14T18:12:47+01:00",
      "commit_date": "2015-01-14T18:12:47+01:00",
      "file_changed_count": 2,
      "insertions_count": 89,
      "deletions_count": 2
    },
    ...
  ]
}
```

And with deltas enabled (without patches):

```
{
  "name": "repotool",
  "vcs": "git",
  "clone_url": "https://github.com/DevMine/repotool.git",
  "clone_path": "/home/robin/Hacking/repotool",
  "default_branch": "master",
  "commits": [
    {
      "vcs_id": "863f9ed113f06829359d0fd4040ae4a6b5c1cf5e",
      "message": "tools/batch: Use a channel to create a pool of tasks for goroutines.\n\nUse a channel on which each tasks (ie call to repotool) is added.\nThis allows to have goroutines picking up tasks from the channel as soon\nas they are done. This way, there is no waiting time as long as there\nare tasks in the pool.\n",
      "author": {
        "name": "Robin Hahling",
        "email": "robin.hahling@gw-computing.net"
      },
      "committer": {
        "name": "Robin Hahling",
        "email": "robin.hahling@gw-computing.net"
      },
      "author_date": "2015-01-13T15:24:40+01:00",
      "commit_date": "2015-01-13T15:24:40+01:00",
      "diff_delta": [
        {
          "status": "modified",
          "binary": false,
          "old_file_path": "tools/batch.go",
          "new_file_path": "tools/batch.go"
        }
      ],
      "file_changed_count": 1,
      "insertions_count": 25,
      "deletions_count": 23
    },
    ...
  ]
}
```

And you can even include patches:

```
{
  "name": "repotool",
  "vcs": "git",
  "clone_url": "https://github.com/DevMine/repotool.git",
  "clone_path": "/home/robin/Hacking/repotool",
  "default_branch": "master",
  "commits": [
    {
      "vcs_id": "fe8aaac0c7650d8ce9c8f4ddeaa63105b3dd0e9e",
      "message": "repotool: Print repository name before processing db insertions.\n",
      "author": {
        "name": "Robin Hahling",
        "email": "robin.hahling@gw-computing.net"
      },
      "committer": {
        "name": "Robin Hahling",
        "email": "robin.hahling@gw-computing.net"
      },
      "author_date": "2015-01-14T18:14:18+01:00",
      "commit_date": "2015-01-14T18:14:18+01:00",
      "diff_delta": [
        {
          "patch": "diff --git a/repotool.go b/repotool.go\nindex ba1eed0..d1ce7a3 100644\n--- a/repotool.go\n+++ b/repotool.go\n@@ -97,8 +97,9 @@ func main() {\n \t\t}\n \t\tdefer db.Close()\n \n-\t\tfmt.Fprintf(os.Stderr, \"inserting %d commits into the database...\\n\",\n-\t\t\tlen(repository.GetCommits()))\n+\t\tfmt.Fprintf(os.Stderr,\n+\t\t\t\"inserting %d commits from %s repository into the database...\\n\",\n+\t\t\tlen(repository.GetCommits()), repository.GetName())\n \t\ttic := time.Now()\n \t\tinsertRepoData(db, repository)\n \t\ttoc := time.Now()\n",
          "status": "modified",
          "binary": false,
          "old_file_path": "repotool.go",
          "new_file_path": "repotool.go"
        }
      ],
      "file_changed_count": 1,
      "insertions_count": 3,
      "deletions_count": 2
    },
    ...
  ]
}
```

## Installation

`repotool` depends on [git2go](https://github.com/libgit2/git2go), which is a
[Go](http://golang.org/) binding to [libgit2](https://libgit2.github.com/), a C
library that implements `git` core methods. Hence, you need `libgit2` installed
on your system unless you statically compile `libgit2` into `git2go`.

If the requirements are met, installing `repotool` is as simple as running this
command in a terminal (assuming [Go](http://golang.org/) is installed):

    go get github.com/DevMine/repotool/cmd/...

Or you can download a binary for your platform from the DevMine project's
[downloads page](http://devmine.ch/downloads).

## Usage

`repotool` produces JSON, provided that you feed it with a path to a source code
repository managed by a VCS which can be either in the form of a directory or a
tar archive. By default, informative messages are outputted to `stderr` whereas
JSON is outputted to `stdout`. To see the list of available options, use the
`-h` flag. Example usage:

    repotool ~/Code/myawesomeproject > myawesomeproject.json

`repotool-db` can be used to insert data into the PostgreSQL database.
You need to provide a configuration file in argument. Simply copy
`repotool.conf.sample` to `repotool.conf` and adjust database connection
information at the very least. See this
[README.md](https://github.com/DevMine/repotool/blob/master/db/README.md) for
more information about the database schema.
`repotool-db` can be used to process multiple repositories in parallel. This is
why, as opposed to `repotool`, it does not simply take a repository as argument
but it takes a directory where it expects to find source code repositories. The
depth at which repositories are expected to be found can be specified with the
`depth` flag. To see the list of available options, use the `-h` flag. Example
usage:

    repotool-db -c repotool.conf ~/Code

With the configuration file, you can also tell `repotool-db` to insert commit
deltas and commits patches (the latter works only if you enable commit deltas,
quite logically). Simply set the `commit_deltas` and to `true`. Note that the
`commit_patches` option is ignored for now. However, you should know that
inserting `commit_patches` slow things down a lot. `repotool-db` can process
repositories concurrently by recursively traversing directories, spawning
goroutines in the process. When using it, bear in mind that `repotool-db` is IO
and CPU intensive, hence do not spawn too many goroutines or you might reach the
number of open files limit. The number of goroutines can be adjusted with the
`-g` parameter. Using about the same number of goroutines as the number of cpu
cores should be a reasonable choice.

As `libgit2` does not support reading information directly from a tar archive,
when given a git repository as a tar archive, `repotool`, or `repotool-db` will
extract part of the archive into a temporary location. You can specify where
using `tmp_dir` in the configuration file for `repotool-db` or by given the
information as argument to `repotool`. We advise specifying a path to a ramdisk
for increased performance and reduced main storage I/Os. When using a ramdisk
with limited capacity, you shall specify the largest size for a tar archive to
be extracted in `tmp_dir` using the `tmp_dir_file_size_limit` option from the
configuration file for `repotool-db` or by using the appropriate flag for
`repotool`. Every tar archive larger than this size will be extracted in its
storage location instead.
