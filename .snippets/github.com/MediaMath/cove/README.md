# Cove - Libraries & CLI's that wrap the go command.

This is a small set of wrapper libraries and cli's that wrap around the golang go command.

## Runtime dependencies

Unlike most go projects this one does have a dependency on the go development environment.  It will not work with just the runtime.

## gosh - Get Over SsH - simple script for getting go packages at supplied uri's.

Note that gosh has issues being built with go1.8 and above. 

### Explicit ssh urls

Gosh works by taking a package/git ssh url pair, seperated by a ','.

```bash
gosh github.com/MediaMath/cove,git@github.com:MediaMath/cove.git github.com/MediaMath/foo/bar,git@github.com:MediaMath/foo.git
```

### Implied github ssh urls.

If you do not provide a second argument with a ',' to gosh it will attempt to imply a github url.

```bash
#will use git@github.com:MediaMath/cove.git and git@github.com:MediaMath/foo.git
gosh github.com/MediaMath/cove github.com/MediaMath/foo/bar
```

### To Install gosh

```bash
go install github.com/MediaMath/cove/gosh
```


## cvr - wrapper around go code coverage.

The go code coverage facility is one of the few that does not respond well to path traversals.  This tool provides an ability to:

- Open go coverage reports in a web browser from the command line in 1 command.
- Generate go coverage reports for multiple projects that match query paths in 1 go. 

Run in a go project directory, this will generate and open the go coverage html report in  the configured web browser:

```bash
cvr
```

Run in any directory, this will create a list of html coverage reports in the specified output directory.

```bash
cvr -o=path/to/some/dir github.com/MediaMath/... 
```

### To Install cvr

```bash
go install github.com/MediaMath/cove/cvr
```
